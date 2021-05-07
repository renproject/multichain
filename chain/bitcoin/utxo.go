package bitcoin

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/renproject/multichain/api/utxo"
	"github.com/renproject/pack"
)

// Version of Bitcoin transactions supported by the multichain.
const Version int32 = 2

// The TxBuilder is an implementation of a UTXO-compatible transaction builder
// for Bitcoin.
type TxBuilder struct {
	params *chaincfg.Params
}

// NewTxBuilder returns a transaction builder that builds UTXO-compatible
// Bitcoin transactions for the given chain configuration (this means that it
// can be used for regnet, testnet, and mainnet, but also for networks that are
// minimally modified forks of the Bitcoin network).
func NewTxBuilder(params *chaincfg.Params) TxBuilder {
	return TxBuilder{params: params}
}

// BuildTx returns a Bitcoin transaction that consumes funds from the given
// inputs, and sends them to the given recipients. The difference in the sum
// value of the inputs and the sum value of the recipients is paid as a fee to
// the Bitcoin network. This fee must be calculated independently of this
// function. Outputs produced for recipients will use P2PKH, P2SH, P2WPKH, or
// P2WSH scripts as the pubkey script, based on the format of the recipient
// address.
func (txBuilder TxBuilder) BuildTx(inputs []utxo.Input, recipients []utxo.Recipient) (utxo.Tx, error) {
	msgTx := wire.NewMsgTx(Version)

	// Inputs
	for _, input := range inputs {
		hash := chainhash.Hash{}
		copy(hash[:], input.Hash)
		index := input.Index.Uint32()
		msgTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(&hash, index), nil, nil))
	}

	// Outputs
	for _, recipient := range recipients {
		addr, err := btcutil.DecodeAddress(string(recipient.To), txBuilder.params)
		if err != nil {
			return nil, err
		}
		script, err := txscript.PayToAddrScript(addr)
		if err != nil {
			return nil, err
		}
		value := recipient.Value.Int().Int64()
		if value < 0 {
			return nil, fmt.Errorf("expected value >= 0, got value %v", value)
		}
		msgTx.AddTxOut(wire.NewTxOut(value, script))
	}

	return &Tx{inputs: inputs, recipients: recipients, msgTx: msgTx, signed: false}, nil
}

// Tx represents a simple Bitcoin transaction that implements the Bitcoin Compat
// API.
type Tx struct {
	inputs     []utxo.Input
	recipients []utxo.Recipient

	msgTx *wire.MsgTx

	signed bool
}

// Hash returns the transaction hash of the given underlying transaction.
func (tx *Tx) Hash() (pack.Bytes, error) {
	txhash := tx.msgTx.TxHash()
	return pack.NewBytes(txhash[:]), nil
}

// Inputs returns the UTXO inputs in the underlying transaction.
func (tx *Tx) Inputs() ([]utxo.Input, error) {
	return tx.inputs, nil
}

// Outputs returns the UTXO outputs in the underlying transaction.
func (tx *Tx) Outputs() ([]utxo.Output, error) {
	hash, err := tx.Hash()
	if err != nil {
		return nil, fmt.Errorf("bad hash: %v", err)
	}
	outputs := make([]utxo.Output, len(tx.msgTx.TxOut))
	for i := range outputs {
		outputs[i].Outpoint = utxo.Outpoint{
			Hash:  hash,
			Index: pack.NewU32(uint32(i)),
		}
		outputs[i].PubKeyScript = pack.Bytes(tx.msgTx.TxOut[i].PkScript)
		if tx.msgTx.TxOut[i].Value < 0 {
			return nil, fmt.Errorf("bad output %v: value is less than zero", i)
		}
		outputs[i].Value = pack.NewU256FromU64(pack.NewU64(uint64(tx.msgTx.TxOut[i].Value)))
	}
	return outputs, nil
}

// Sighashes returns the digests that must be signed before the transaction
// can be submitted by the client.
func (tx *Tx) Sighashes() ([]pack.Bytes32, error) {
	sighashes := make([]pack.Bytes32, len(tx.inputs))

	for i, txin := range tx.inputs {
		pubKeyScript := txin.PubKeyScript
		sigScript := txin.SigScript
		value := txin.Value.Int().Int64()
		if value < 0 {
			return []pack.Bytes32{}, fmt.Errorf("expected value >= 0, got value %v", value)
		}

		var hash []byte
		var err error
		if sigScript == nil {
			if txscript.IsPayToWitnessPubKeyHash(pubKeyScript) {
				hash, err = txscript.CalcWitnessSigHash(pubKeyScript, txscript.NewTxSigHashes(tx.msgTx), txscript.SigHashAll, tx.msgTx, i, value)
			} else {
				hash, err = txscript.CalcSignatureHash(pubKeyScript, txscript.SigHashAll, tx.msgTx, i)
			}
		} else {
			if txscript.IsPayToWitnessScriptHash(pubKeyScript) {
				hash, err = txscript.CalcWitnessSigHash(sigScript, txscript.NewTxSigHashes(tx.msgTx), txscript.SigHashAll, tx.msgTx, i, value)
			} else {
				hash, err = txscript.CalcSignatureHash(sigScript, txscript.SigHashAll, tx.msgTx, i)
			}
		}
		if err != nil {
			return []pack.Bytes32{}, err
		}

		sighash := [32]byte{}
		copy(sighash[:], hash)
		sighashes[i] = pack.NewBytes32(sighash)
	}

	return sighashes, nil
}

// Sign consumes a list of signatures, and adds them to the list of UTXOs in
// the underlying transactions.
func (tx *Tx) Sign(signatures []pack.Bytes65, pubKey pack.Bytes) error {
	if tx.signed {
		return fmt.Errorf("already signed")
	}
	if len(signatures) != len(tx.msgTx.TxIn) {
		return fmt.Errorf("expected %v signatures, got %v signatures", len(tx.msgTx.TxIn), len(signatures))
	}

	for i, rsv := range signatures {
		var err error

		// Decode the signature and the pubkey script.
		r := new(big.Int).SetBytes(rsv[:32])
		s := new(big.Int).SetBytes(rsv[32:64])
		signature := btcec.Signature{
			R: r,
			S: s,
		}
		pubKeyScript := tx.inputs[i].Output.PubKeyScript
		sigScript := tx.inputs[i].SigScript

		// Support segwit.
		if sigScript == nil {
			if txscript.IsPayToWitnessPubKeyHash(pubKeyScript) || txscript.IsPayToWitnessScriptHash(pubKeyScript) {
				tx.msgTx.TxIn[i].Witness = wire.TxWitness([][]byte{append(signature.Serialize(), byte(txscript.SigHashAll)), pubKey})
				continue
			}
		} else {
			if txscript.IsPayToWitnessScriptHash(sigScript) || txscript.IsPayToWitnessScriptHash(sigScript) {
				tx.msgTx.TxIn[i].Witness = wire.TxWitness([][]byte{append(signature.Serialize(), byte(txscript.SigHashAll)), pubKey, sigScript})
				continue
			}
		}

		// Support non-segwit
		builder := txscript.NewScriptBuilder()
		builder.AddData(append(signature.Serialize(), byte(txscript.SigHashAll)))
		builder.AddData(pubKey)
		if sigScript != nil {
			builder.AddData(sigScript)
		}
		tx.msgTx.TxIn[i].SignatureScript, err = builder.Script()
		if err != nil {
			return err
		}
	}

	tx.signed = true
	return nil
}

// Serialize serializes the UTXO transaction to bytes
func (tx *Tx) Serialize() (pack.Bytes, error) {
	buf := new(bytes.Buffer)
	if err := tx.msgTx.Serialize(buf); err != nil {
		return pack.Bytes{}, err
	}
	return pack.NewBytes(buf.Bytes()), nil
}

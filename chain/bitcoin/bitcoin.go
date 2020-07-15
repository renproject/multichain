package bitcoin

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/renproject/multichain/compat/bitcoincompat"
	"github.com/renproject/pack"
)

// Version of Bitcoin transactions supported by the multichain.
const Version int32 = 2

type txBuilder struct{}

// NewTxBuilder returns an implementation of the transaction builder interface
// from the Bitcoin Compat API, and exposes the functionality to build simple
// Bitcoin transactions.
func NewTxBuilder() bitcoincompat.TxBuilder {
	return txBuilder{}
}

// BuildTx returns a simple Bitcoin transaction that consumes the funds from the
// given outputs, and sends the to the given recipients. The difference in the
// sum value of the inputs and the sum value of the recipients is paid as a fee
// to the Bitcoin network.
//
// It is assumed that the required signature scripts require the SIGHASH_ALL
// signatures and the serialized public key:
//
//  builder := txscript.NewScriptBuilder()
//  builder.AddData(append(signature.Serialize(), byte(txscript.SigHashAll)))
//  builder.AddData(serializedPubKey)
//
// Outputs produced for recipients will use P2PKH, P2SH, P2WPKH, or P2WSH
// scripts as the pubkey script, based on the format of the recipient address.
func (txBuilder) BuildTx(inputs []bitcoincompat.Output, recipients []bitcoincompat.Recipient) (bitcoincompat.Tx, error) {
	msgTx := wire.NewMsgTx(Version)

	// Inputs
	for _, input := range inputs {
		hash := chainhash.Hash(input.Outpoint.Hash)
		index := input.Outpoint.Index.Uint32()
		msgTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(&hash, index), nil, nil))
	}

	// Outputs
	for _, recipient := range recipients {
		script, err := txscript.PayToAddrScript(recipient.Address)
		if err != nil {
			return &Tx{}, err
		}
		value := int64(recipient.Value.Uint64())
		if value < 0 {
			return &Tx{}, fmt.Errorf("expected value >= 0, got value = %v", value)
		}
		msgTx.AddTxOut(wire.NewTxOut(value, script))
	}

	return &Tx{inputs: inputs, recipients: recipients, msgTx: msgTx, signed: false}, nil
}

// Tx represents a simple Bitcoin transaction that implements the Bitcoin Compat
// API.
type Tx struct {
	inputs     []bitcoincompat.Output
	recipients []bitcoincompat.Recipient

	msgTx *wire.MsgTx

	signed bool
}

func (tx *Tx) Hash() pack.Bytes32 {
	return pack.NewBytes32(tx.msgTx.TxHash())
}

func (tx *Tx) Sighashes() ([]pack.Bytes32, error) {
	sighashes := make([]pack.Bytes32, len(tx.inputs))

	for i := range tx.inputs {
		pubKeyScript := tx.inputs[i].PubKeyScript
		value := int64(tx.inputs[i].Value.Uint64())
		if value < 0 {
			return []pack.Bytes32{}, fmt.Errorf("expected value >= 0, got value = %v", value)
		}

		var hash []byte
		var err error
		if txscript.IsPayToWitnessPubKeyHash(pubKeyScript) {
			hash, err = txscript.CalcWitnessSigHash(pubKeyScript, txscript.NewTxSigHashes(tx.msgTx), txscript.SigHashAll, tx.msgTx, i, value)
		} else {
			hash, err = txscript.CalcSignatureHash(pubKeyScript, txscript.SigHashAll, tx.msgTx, i)
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

func (tx *Tx) Sign(signatures []pack.Bytes65, pubKey pack.Bytes) error {
	if tx.signed {
		return fmt.Errorf("signed")
	}
	if len(signatures) != len(tx.msgTx.TxIn) {
		return fmt.Errorf("expected %v signatures, got %v signatures", len(tx.msgTx.TxIn), len(signatures))
	}

	for i, rsv := range signatures {
		// Decode the signature and the pubkey script.
		r := new(big.Int).SetBytes(rsv[:32])
		s := new(big.Int).SetBytes(rsv[32:64])
		signature := btcec.Signature{
			R: r,
			S: s,
		}
		pubKeyScript := tx.inputs[i].PubKeyScript

		// Support the consumption of SegWit outputs.
		if txscript.IsPayToWitnessPubKeyHash(pubKeyScript) || txscript.IsPayToWitnessScriptHash(pubKeyScript) {
			tx.msgTx.TxIn[i].Witness = wire.TxWitness([][]byte{append(signature.Serialize(), byte(txscript.SigHashAll)), pubKey})
			continue
		}

		// Support the consumption of non-SegWite outputs.
		builder := txscript.NewScriptBuilder()
		builder.AddData(append(signature.Serialize(), byte(txscript.SigHashAll)))
		builder.AddData(pubKey)
		signatureScript, err := builder.Script()
		if err != nil {
			return err
		}
		tx.msgTx.TxIn[i].SignatureScript = signatureScript
	}

	tx.signed = true
	return nil
}

func (tx *Tx) Serialize() (pack.Bytes, error) {
	buf := new(bytes.Buffer)
	if err := tx.msgTx.Serialize(buf); err != nil {
		return pack.Bytes{}, err
	}
	return pack.NewBytes(buf.Bytes()), nil
}

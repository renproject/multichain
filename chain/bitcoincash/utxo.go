package bitcoincash

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/renproject/multichain/api/utxo"
	"github.com/renproject/multichain/chain/bitcoin"
	"github.com/renproject/pack"
)

// SighashForkID used to distinguish between Bitcoin Cash and Bitcoin
// transactions by masking hash types.
const SighashForkID = txscript.SigHashType(0x40)

// SighashMask used to mask hash types.
const SighashMask = txscript.SigHashType(0x1F)

// Version of Bitcoin Cash transactions supported by the multichain.
const Version int32 = 1

// ClientOptions are used to parameterise the behaviour of the Client.
type ClientOptions = bitcoin.ClientOptions

// DefaultClientOptions returns ClientOptions with the default settings. These
// settings are valid for use with the default local deployment of the
// multichain. In production, the host, user, and password should be changed.
func DefaultClientOptions() ClientOptions {
	return bitcoin.DefaultClientOptions().WithHost("http://127.0.0.1:19443")
}

// A Client interacts with an instance of the Bitcoin network using the RPC
// interface exposed by a Bitcoin node.
type Client = bitcoin.Client

// NewClient returns a new Client.
var NewClient = bitcoin.NewClient

// The TxBuilder is an implementation of a UTXO-compatible transaction builder
// for Bitcoin.
type TxBuilder struct {
	params *chaincfg.Params
}

// NewTxBuilder returns an implementation of the transaction builder interface
// from the Bitcoin Compat API, and exposes the functionality to build simple
// Bitcoin Cash transactions.
func NewTxBuilder(params *chaincfg.Params) utxo.TxBuilder {
	return TxBuilder{params: params}
}

// BuildTx returns a simple Bitcoin Cash transaction that consumes the funds
// from the given outputs, and sends the to the given recipients. The difference
// in the sum value of the inputs and the sum value of the recipients is paid as
// a fee to the Bitcoin Cash network.
//
// It is assumed that the required signature scripts require the SIGHASH_ALL
// signatures and the serialized public key:
//
//  builder := txscript.NewScriptBuilder()
//  builder.AddData(append(signature.Serialize(), byte(txscript.SigHashAll|SighashForkID)))
//  builder.AddData(serializedPubKey)
//
// Outputs produced for recipients will use P2PKH, or P2SH scripts as the pubkey
// script, based on the format of the recipient address.
func (txBuilder TxBuilder) BuildTx(inputs []utxo.Input, recipients []utxo.Recipient) (utxo.Tx, error) {
	msgTx := wire.NewMsgTx(Version)

	// Address encoder-decoder
	addrEncodeDecoder := NewAddressEncodeDecoder(txBuilder.params)

	// Inputs
	for _, input := range inputs {
		hash := chainhash.Hash{}
		copy(hash[:], input.Hash)
		index := input.Index.Uint32()
		msgTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(&hash, index), nil, nil))
	}

	// Outputs
	for _, recipient := range recipients {
		addrBytes, err := addrEncodeDecoder.DecodeAddress(recipient.To)
		if err != nil {
			return &Tx{}, err
		}
		addr, err := addressFromRawBytes(addrBytes, txBuilder.params)
		if err != nil {
			return &Tx{}, err
		}
		script, err := txscript.PayToAddrScript(addr.BitcoinAddress())
		if err != nil {
			return &Tx{}, err
		}
		value := recipient.Value.Int().Int64()
		if value < 0 {
			return nil, fmt.Errorf("expected value >= 0, got value = %v", value)
		}
		msgTx.AddTxOut(wire.NewTxOut(value, script))
	}

	return &Tx{inputs: inputs, recipients: recipients, msgTx: msgTx, signed: false}, nil
}

// Tx represents a simple Bitcoin Cash transaction that implements the Bitcoin
// Compat API.
type Tx struct {
	inputs     []utxo.Input
	recipients []utxo.Recipient

	msgTx *wire.MsgTx

	signed bool
}

// Hash returns the transaction hash of the given underlying transaction. It
// implements the multichain.UTXOTx interface
func (tx *Tx) Hash() (pack.Bytes, error) {
	txhash := tx.msgTx.TxHash()
	return pack.NewBytes(txhash[:]), nil
}

// Inputs returns the UTXO inputs in the underlying transaction. It implements
// the multichain.UTXOTx interface
func (tx *Tx) Inputs() ([]utxo.Input, error) {
	return tx.inputs, nil
}

// Outputs returns the UTXO outputs in the underlying transaction. It implements
// the multichain.UTXOTx interface
func (tx *Tx) Outputs() ([]utxo.Output, error) {
	hash, err := tx.Hash()
	if err != nil {
		return nil, err
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
		pubKeyScript := txin.Output.PubKeyScript
		sigScript := txin.SigScript
		value := txin.Output.Value.Int().Int64()
		if value < 0 {
			return []pack.Bytes32{}, fmt.Errorf("expected value >= 0, got value = %v", value)
		}

		var hash []byte
		if sigScript == nil {
			hash = CalculateBip143Sighash(pubKeyScript, txscript.NewTxSigHashes(tx.msgTx), txscript.SigHashAll, tx.msgTx, i, value)
		} else {
			hash = CalculateBip143Sighash(sigScript, txscript.NewTxSigHashes(tx.msgTx), txscript.SigHashAll, tx.msgTx, i, value)
		}

		sighash := [32]byte{}
		copy(sighash[:], hash)
		sighashes[i] = pack.NewBytes32(sighash)
	}
	return sighashes, nil
}

// Sign consumes a list of signatures, and adds them to the list of UTXOs in
// the underlying transactions. It implements the multichain.UTXOTx interface
func (tx *Tx) Sign(signatures []pack.Bytes65, pubKey pack.Bytes) error {
	if tx.signed {
		return fmt.Errorf("already signed")
	}
	if len(signatures) != len(tx.msgTx.TxIn) {
		return fmt.Errorf("expected %v signatures, got %v signatures", len(tx.msgTx.TxIn), len(signatures))
	}

	for i, rsv := range signatures {
		r := new(big.Int).SetBytes(rsv[:32])
		s := new(big.Int).SetBytes(rsv[32:64])
		signature := btcec.Signature{
			R: r,
			S: s,
		}

		builder := txscript.NewScriptBuilder()
		builder.AddData(append(signature.Serialize(), byte(txscript.SigHashAll|SighashForkID)))
		builder.AddData(pubKey)
		if tx.inputs[i].SigScript != nil {
			builder.AddData(tx.inputs[i].SigScript)
		}
		signatureScript, err := builder.Script()
		if err != nil {
			return err
		}
		tx.msgTx.TxIn[i].SignatureScript = signatureScript
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

// CalculateBip143Sighash computes the sighash digest of a transaction's input
// using the new, optimized digest calculation algorithm defined in BIP0143.
// This function makes use of pre-calculated sighash fragments stored within the
// passed HashCache to eliminate duplicate hashing computations when calculating
// the final digest, reducing the complexity from O(N^2) to O(N). Additionally,
// signatures now cover the input value of the referenced unspent output. This
// allows offline, or hardware wallets to compute the exact amount being spent,
// in addition to the final transaction fee. In the case the wallet if fed an
// invalid input amount, the real sighash will differ causing the produced
// signature to be invalid.
//
// https://github.com/bitcoin/bips/blob/master/bip-0143.mediawiki
func CalculateBip143Sighash(subScript []byte, sigHashes *txscript.TxSigHashes, hashType txscript.SigHashType, tx *wire.MsgTx, idx int, amt int64) []byte {

	// As a sanity check, ensure the passed input index for the transaction
	// is valid.
	if idx > len(tx.TxIn)-1 {
		fmt.Printf("CalculateBip143Sighash: i %d with %d inputs", idx, len(tx.TxIn))
		return nil
	}

	// We'll utilize this buffer throughout to incrementally calculate
	// the signature hash for this transaction.
	var sigHash bytes.Buffer

	// First write out, then encode the transaction's version number.
	var bVersion [4]byte
	binary.LittleEndian.PutUint32(bVersion[:], uint32(tx.Version))
	sigHash.Write(bVersion[:])

	// Next write out the possibly pre-calculated hashes for the sequence
	// numbers of all inputs, and the hashes of the previous outs for all
	// outputs.
	var zeroHash chainhash.Hash

	// If anyone can pay isn't active, then we can use the cached
	// hashPrevOuts, otherwise we just write zeroes for the prev outs.
	if hashType&txscript.SigHashAnyOneCanPay == 0 {
		sigHash.Write(sigHashes.HashPrevOuts[:])
	} else {
		sigHash.Write(zeroHash[:])
	}

	// If the sighash isn't anyone can pay, single, or none, the use the
	// cached hash sequences, otherwise write all zeroes for the
	// hashSequence.
	if hashType&txscript.SigHashAnyOneCanPay == 0 &&
		hashType&SighashMask != txscript.SigHashSingle &&
		hashType&SighashMask != txscript.SigHashNone {
		sigHash.Write(sigHashes.HashSequence[:])
	} else {
		sigHash.Write(zeroHash[:])
	}

	// Next, write the outpoint being spent.
	sigHash.Write(tx.TxIn[idx].PreviousOutPoint.Hash[:])
	var bIndex [4]byte
	binary.LittleEndian.PutUint32(bIndex[:], tx.TxIn[idx].PreviousOutPoint.Index)
	sigHash.Write(bIndex[:])

	// For p2wsh outputs, and future outputs, the script code is the
	// original script, with all code separators removed, serialized
	// with a var int length prefix.
	wire.WriteVarBytes(&sigHash, 0, subScript)

	// Next, add the input amount, and sequence number of the input being
	// signed.
	var bAmount [8]byte
	binary.LittleEndian.PutUint64(bAmount[:], uint64(amt))
	sigHash.Write(bAmount[:])
	var bSequence [4]byte
	binary.LittleEndian.PutUint32(bSequence[:], tx.TxIn[idx].Sequence)
	sigHash.Write(bSequence[:])

	// If the current signature mode isn't single, or none, then we can
	// re-use the pre-generated hashoutputs sighash fragment. Otherwise,
	// we'll serialize and add only the target output index to the signature
	// pre-image.
	if hashType&SighashMask != txscript.SigHashSingle &&
		hashType&SighashMask != txscript.SigHashNone {
		sigHash.Write(sigHashes.HashOutputs[:])
	} else if hashType&SighashMask == txscript.SigHashSingle && idx < len(tx.TxOut) {
		var b bytes.Buffer
		wire.WriteTxOut(&b, 0, 0, tx.TxOut[idx])
		sigHash.Write(chainhash.DoubleHashB(b.Bytes()))
	} else {
		sigHash.Write(zeroHash[:])
	}

	// Finally, write out the transaction's locktime, and the sig hash
	// type.
	var bLockTime [4]byte
	binary.LittleEndian.PutUint32(bLockTime[:], tx.LockTime)
	sigHash.Write(bLockTime[:])
	var bHashType [4]byte
	binary.LittleEndian.PutUint32(bHashType[:], uint32(hashType|SighashForkID))
	sigHash.Write(bHashType[:])

	return chainhash.DoubleHashB(sigHash.Bytes())
}

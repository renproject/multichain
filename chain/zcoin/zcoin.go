package zcoin

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/renproject/multichain/compat/bitcoincompat"
	"github.com/renproject/pack"
)

// Version of Zcoin transactions supported by the multichain.
const Version int32 = 1

type txBuilder struct {
	params *chaincfg.Params
}

// NewTxBuilder returns an implementation the transaction builder interface from
// the Bitcoin Compat API, and exposes the functionality to build simple Zcoin
// transactions.
func NewTxBuilder(params *chaincfg.Params) bitcoincompat.TxBuilder {
	return txBuilder{params: params}
}

// BuildTx returns a simple Zcoin transaction that consumes the funds from the
// given outputs, and sends the to the given recipients. The difference in the
// sum value of the inputs and the sum value of the recipients is paid as a fee
// to the Zcoin network.
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
func (txBuilder txBuilder) BuildTx(inputs []bitcoincompat.Output, recipients []bitcoincompat.Recipient) (bitcoincompat.Tx, error) {
	msgTx := wire.NewMsgTx(Version)

	// Inputs
	for _, input := range inputs {
		hash := chainhash.Hash(input.Outpoint.Hash)
		index := input.Outpoint.Index.Uint32()
		msgTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(&hash, index), nil, nil))
	}

	// Outputs
	for _, recipient := range recipients {
		var script []byte
		var err error
		switch addr := recipient.Address.(type) {
		case AddressPubKeyHash:
			script, err = txscript.PayToAddrScript(addr.AddressPubKeyHash)
		default:
			script, err = txscript.PayToAddrScript(recipient.Address)
		}
		if err != nil {
			return nil, err
		}
		value := int64(recipient.Value.Uint64())
		if value < 0 {
			return nil, fmt.Errorf("expected value >= 0, got value = %v", value)
		}
		msgTx.AddTxOut(wire.NewTxOut(value, script))
	}

	return &Tx{inputs: inputs, recipients: recipients, msgTx: msgTx, signed: false}, nil
}

// Tx represents a simple Zcoin transaction that implements the Bitcoin Compat
// API.
type Tx struct {
	inputs     []bitcoincompat.Output
	recipients []bitcoincompat.Recipient

	msgTx        *wire.MsgTx

	signed bool
}

func (tx *Tx) Hash() pack.Bytes32 {
	serial, err := tx.Serialize()
	if err != nil {
		return pack.Bytes32{}
	}
	return pack.NewBytes32(chainhash.DoubleHashH(serial))
}

func (tx *Tx) Sighashes() ([]pack.Bytes32, error) {
	sighashes := make([]pack.Bytes32, len(tx.inputs))
	for i, txin := range tx.inputs {
		pubKeyScript := txin.PubKeyScript

		var hash []byte
		var err error
               hash, err = txscript.CalcSignatureHash(pubKeyScript, txscript.SigHashAll, tx.msgTx, i)
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

// AddressPubKeyHash represents an address for P2PKH transactions for Zcoin that
// is compatible with the Bitcoin Compat API.
type AddressPubKeyHash struct {
	*btcutil.AddressPubKeyHash
	params *chaincfg.Params
}

// NewAddressPubKeyHash returns a new AddressPubKeyHash that is compatible with
// the Bitcoin Compat API.
func NewAddressPubKeyHash(pkh []byte, params *chaincfg.Params) (AddressPubKeyHash, error) {
	addr, err := btcutil.NewAddressPubKeyHash(pkh, params)
	return AddressPubKeyHash{AddressPubKeyHash: addr, params: params}, err
}

// String returns the string encoding of the transaction output destination.
//
// Please note that String differs subtly from EncodeAddress: String will return
// the value as a string without any conversion, while EncodeAddress may convert
// destination types (for example, converting pubkeys to P2PKH addresses) before
// encoding as a payment address string.

// EncodeAddress returns the string encoding of the payment address associated
// with the Address value. See the comment on String for how this method differs
// from String.
func (addr AddressPubKeyHash) EncodeAddress() string {
	hash := *addr.AddressPubKeyHash.Hash160()
	var prefix []byte
	switch addr.params {
	case &chaincfg.RegressionNetParams:
		prefix = regnet.p2pkhPrefix
	case &chaincfg.TestNet3Params:
		prefix = testnet.p2pkhPrefix
	case &chaincfg.MainNetParams:
		prefix = mainnet.p2pkhPrefix
	}
	return encodeAddress(hash[:], prefix)
}

func encodeAddress(hash, prefix []byte) string {
	var (
		body  = append(prefix, hash...)
		chk   = checksum(body)
		cksum [4]byte
	)
	copy(cksum[:], chk[:4])
	return base58.Encode(append(body, cksum[:]...))
}

func checksum(input []byte) (cksum [4]byte) {
	var (
		h  = sha256.Sum256(input)
		h2 = sha256.Sum256(h[:])
	)
	copy(cksum[:], h2[:4])
	return
}

type netParams struct {
	name   string
	params *chaincfg.Params

	p2shPrefix    []byte
	p2pkhPrefix   []byte
}

var (
	mainnet = netParams{
		name:   "mainnet",
		params: &chaincfg.MainNetParams,

		p2pkhPrefix: []byte{0x52},
		p2shPrefix:  []byte{0x07},
	}
	testnet = netParams{
		name:   "testnet",
		params: &chaincfg.TestNet3Params,

		p2pkhPrefix: []byte{0x41},
		p2shPrefix:  []byte{0xB2},
	}
	regnet = netParams{
		name:   "regtest",
		params: &chaincfg.RegressionNetParams,

		p2pkhPrefix: []byte{0x41},
		p2shPrefix:  []byte{0xB2},
	}
)

package zcash

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/codahale/blake2"
	"github.com/renproject/multichain/compat/bitcoincompat"
	"github.com/renproject/pack"
	"golang.org/x/crypto/ripemd160"
)

// Version of Zcash transactions supported by the multichain.
const Version int32 = 4

type txBuilder struct {
	params *chaincfg.Params
}

// NewTxBuilder returns an implementation the transaction builder interface from
// the Bitcoin Compat API, and exposes the functionality to build simple Zcash
// transactions.
func NewTxBuilder(params *chaincfg.Params) bitcoincompat.TxBuilder {
	return txBuilder{params: params}
}

// BuildTx returns a simple Zcash transaction that consumes the funds from the
// given outputs, and sends the to the given recipients. The difference in the
// sum value of the inputs and the sum value of the recipients is paid as a fee
// to the Zcash network.
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
func (txBuilder txBuilder) BuildTx(inputs []bitcoincompat.Input, recipients []bitcoincompat.Recipient) (bitcoincompat.Tx, error) {
	msgTx := wire.NewMsgTx(Version)

	// Inputs
	for _, input := range inputs {
		hash := chainhash.Hash(input.Output.Outpoint.Hash)
		index := input.Output.Outpoint.Index.Uint32()
		msgTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(&hash, index), nil, nil))
	}

	// Outputs
	for _, recipient := range recipients {
		addr, err := DecodeAddress(string(recipient.Address))
		if err != nil {
			return &Tx{}, err
		}
		script, err := txscript.PayToAddrScript(addr.BitcoinCompatAddress())
		if err != nil {
			return &Tx{}, err
		}
		value := int64(recipient.Value.Uint64())
		if value < 0 {
			return nil, fmt.Errorf("expected value >= 0, got value = %v", value)
		}
		msgTx.AddTxOut(wire.NewTxOut(value, script))
	}

	var expiryHeight uint32
	switch txBuilder.params {
	case &chaincfg.RegressionNetParams:
		expiryHeight = regnet.expiryHeight
	case &chaincfg.TestNet3Params:
		expiryHeight = testnet.expiryHeight
	case &chaincfg.MainNetParams:
		expiryHeight = mainnet.expiryHeight
	}

	return &Tx{inputs: inputs, recipients: recipients, msgTx: msgTx, expiryHeight: expiryHeight, signed: false}, nil
}

// Tx represents a simple Zcash transaction that implements the Bitcoin Compat
// API.
type Tx struct {
	inputs     []bitcoincompat.Input
	recipients []bitcoincompat.Recipient

	msgTx        *wire.MsgTx
	expiryHeight uint32

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
		sigScript := txin.SigScript
		value := int64(txin.Output.Value.Uint64())
		if value < 0 {
			return []pack.Bytes32{}, fmt.Errorf("expected value >= 0, got value = %v", value)
		}

		hash, err := calculateSighash(regnet, sigScript, txscript.SigHashAll, tx.msgTx, i, value, tx.expiryHeight)
		if err != nil {
			return []pack.Bytes32{}, err
		}

		sighash := [32]byte{}
		copy(sighash[:], hash)
		sighashes[i] = pack.NewBytes32(sighash)
	}
	return sighashes, nil
}

func (tx *Tx) Outputs() ([]bitcoincompat.Output, error) {
	hash := tx.Hash()
	outputs := make([]bitcoincompat.Output, len(tx.msgTx.TxOut))
	for i := range outputs {
		outputs[i].Outpoint = bitcoincompat.Outpoint{
			Hash:  hash,
			Index: pack.NewU32(uint32(i)),
		}
		outputs[i].PubKeyScript = pack.Bytes(tx.msgTx.TxOut[i].PkScript)
		if tx.msgTx.TxOut[i].Value < 0 {
			return nil, fmt.Errorf("bad output %v: value is less than zero", i)
		}
		outputs[i].Value = pack.NewU64(uint64(tx.msgTx.TxOut[i].Value))
	}
	return outputs, nil
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
	w := new(bytes.Buffer)
	pver := uint32(0)
	enc := wire.BaseEncoding

	if err := binary.Write(w, binary.LittleEndian, uint32(tx.msgTx.Version)|(1<<31)); err != nil {
		return pack.Bytes{}, err
	}

	var versionGroupID = versionOverwinterGroupID
	if tx.msgTx.Version == versionSapling {
		versionGroupID = versionSaplingGroupID
	}

	if err := binary.Write(w, binary.LittleEndian, versionGroupID); err != nil {
		return pack.Bytes{}, err
	}

	// If the encoding nVersion is set to WitnessEncoding, and the Flags
	// field for the MsgTx aren't 0x00, then this indicates the transaction
	// is to be encoded using the new witness inclusionary structure
	// defined in BIP0144.
	doWitness := enc == wire.WitnessEncoding && tx.msgTx.HasWitness()
	if doWitness {
		// After the txn's Version field, we include two additional
		// bytes specific to the witness encoding. The first byte is an
		// always 0x00 marker byte, which allows decoders to
		// distinguish a serialized transaction with witnesses from a
		// regular (legacy) one. The second byte is the Flag field,
		// which at the moment is always 0x01, but may be extended in
		// the future to accommodate auxiliary non-committed fields.
		if _, err := w.Write(witnessMarkerBytes); err != nil {
			return pack.Bytes{}, err
		}
	}

	count := uint64(len(tx.msgTx.TxIn))
	if err := writeVarInt(w, pver, count); err != nil {
		return pack.Bytes{}, err
	}

	for _, ti := range tx.msgTx.TxIn {
		if err := writeTxIn(w, pver, tx.msgTx.Version, ti); err != nil {
			return pack.Bytes{}, err
		}
	}

	count = uint64(len(tx.msgTx.TxOut))
	if err := writeVarInt(w, pver, count); err != nil {
		return pack.Bytes{}, err
	}

	for _, to := range tx.msgTx.TxOut {
		if err := writeTxOut(w, pver, tx.msgTx.Version, to); err != nil {
			return pack.Bytes{}, err
		}
	}

	// If this transaction is a witness transaction, and the witness
	// encoded is desired, then encode the witness for each of the inputs
	// within the transaction.
	if doWitness {
		for _, ti := range tx.msgTx.TxIn {
			if err := writeTxWitness(w, pver, tx.msgTx.Version, ti.Witness); err != nil {
				return pack.Bytes{}, err
			}
		}
	}

	if err := binary.Write(w, binary.LittleEndian, tx.msgTx.LockTime); err != nil {
		return pack.Bytes{}, err
	}

	if err := binary.Write(w, binary.LittleEndian, tx.expiryHeight); err != nil {
		return pack.Bytes{}, err
	}

	if tx.msgTx.Version == versionSapling {
		// valueBalance
		if err := binary.Write(w, binary.LittleEndian, uint64(0)); err != nil {
			return pack.Bytes{}, err
		}

		// nShieldedSpend
		if err := writeVarInt(w, pver, 0); err != nil {
			return pack.Bytes{}, err
		}

		// nShieldedOutput
		if err := writeVarInt(w, pver, 0); err != nil {
			return pack.Bytes{}, err
		}
	}

	if err := writeVarInt(w, pver, 0); err != nil {
		return pack.Bytes{}, err
	}

	return pack.NewBytes(w.Bytes()), nil
}

// An Address represents a Zcash address.
type Address interface {
	btcutil.Address
	BitcoinCompatAddress() btcutil.Address
}

// AddressPubKeyHash represents an address for P2PKH transactions for Zcash that
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
func (addr AddressPubKeyHash) String() string {
	return addr.EncodeAddress()
}

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

// ScriptAddress returns the raw bytes of the address to be used when inserting
// the address into a txout's script.
func (addr AddressPubKeyHash) ScriptAddress() []byte {
	return addr.AddressPubKeyHash.ScriptAddress()
}

// IsForNet returns whether or not the address is associated with the passed
// bitcoin network.
func (addr AddressPubKeyHash) IsForNet(params *chaincfg.Params) bool {
	return addr.AddressPubKeyHash.IsForNet(params)
}

// BitcoinCompatAddress returns the address as if it was a Bitcoin address.
func (addr AddressPubKeyHash) BitcoinCompatAddress() btcutil.Address {
	return addr.AddressPubKeyHash
}

// AddressScriptHash represents an address for P2SH transactions for Zcash that
// is compatible with the Bitcoin Compat API.
type AddressScriptHash struct {
	*btcutil.AddressScriptHash
	params *chaincfg.Params
}

// NewAddressScriptHash returns a new AddressScriptHash that is compatible with
// the Bitcoin Compat API.
func NewAddressScriptHash(pkh []byte, params *chaincfg.Params) (AddressScriptHash, error) {
	addr, err := btcutil.NewAddressScriptHash(pkh, params)
	return AddressScriptHash{AddressScriptHash: addr, params: params}, err
}

// String returns the string encoding of the transaction output destination.
//
// Please note that String differs subtly from EncodeAddress: String will return
// the value as a string without any conversion, while EncodeAddress may convert
// destination types (for example, converting pubkeys to P2PKH addresses) before
// encoding as a payment address string.
func (addr AddressScriptHash) String() string {
	return addr.EncodeAddress()
}

// BitcoinCompatAddress returns the address as if it was a Bitcoin address.
func (addr AddressScriptHash) BitcoinCompatAddress() btcutil.Address {
	return addr.AddressScriptHash
}

// EncodeAddress returns the string encoding of the payment address associated
// with the Address value. See the comment on String for how this method differs
// from String.
func (addr AddressScriptHash) EncodeAddress() string {
	hash := *addr.AddressScriptHash.Hash160()
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

// ScriptAddress returns the raw bytes of the address to be used when inserting
// the address into a txout's script.
func (addr AddressScriptHash) ScriptAddress() []byte {
	return addr.AddressScriptHash.ScriptAddress()
}

// IsForNet returns whether or not the address is associated with the passed
// bitcoin network.
func (addr AddressScriptHash) IsForNet(params *chaincfg.Params) bool {
	return addr.AddressScriptHash.IsForNet(params)
}

func DecodeAddress(addr string) (Address, error) {
	var decoded = base58.Decode(addr)
	if len(decoded) != 26 && len(decoded) != 25 {
		return nil, base58.ErrInvalidFormat
	}

	var cksum [4]byte
	copy(cksum[:], decoded[len(decoded)-4:])
	if checksum(decoded[:len(decoded)-4]) != cksum {
		return nil, base58.ErrChecksum
	}

	if len(decoded)-6 != ripemd160.Size && len(decoded)-5 != ripemd160.Size {
		return nil, errors.New("incorrect payload len")
	}

	var addrType uint8
	var params *chaincfg.Params
	var err error
	var hash [20]byte
	if len(decoded) == 26 {
		addrType, params, err = parsePrefix(decoded[:2])
		copy(hash[:], decoded[2:22])
	} else {
		addrType, params, err = parsePrefix(decoded[:1])
		copy(hash[:], decoded[1:21])
	}
	if err != nil {
		return nil, err
	}

	switch addrType {
	case 0: // P2PKH
		addr, err := btcutil.NewAddressPubKeyHash(hash[:], params)
		if err != nil {
			return nil, err
		}
		return &AddressPubKeyHash{AddressPubKeyHash: addr, params: params}, nil
	case 1: // P2SH
		addr, err := btcutil.NewAddressScriptHash(hash[:], params)
		if err != nil {
			return nil, err
		}
		return &AddressScriptHash{AddressScriptHash: addr, params: params}, nil
	}

	return nil, errors.New("unknown address")
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

func parsePrefix(prefix []byte) (uint8, *chaincfg.Params, error) {
	if bytes.Equal(prefix, mainnet.p2pkhPrefix) {
		return 0, &chaincfg.MainNetParams, nil
	}
	if bytes.Equal(prefix, mainnet.p2shPrefix) {
		return 1, &chaincfg.MainNetParams, nil
	}
	if bytes.Equal(prefix, testnet.p2pkhPrefix) {
		return 0, &chaincfg.TestNet3Params, nil
	}
	if bytes.Equal(prefix, testnet.p2shPrefix) {
		return 1, &chaincfg.TestNet3Params, nil
	}
	if bytes.Equal(prefix, regnet.p2pkhPrefix) {
		return 0, &chaincfg.RegressionNetParams, nil
	}
	if bytes.Equal(prefix, regnet.p2shPrefix) {
		return 1, &chaincfg.RegressionNetParams, nil
	}
	return 0, nil, btcutil.ErrUnknownAddressType
}

type netParams struct {
	name   string
	params *chaincfg.Params

	p2shPrefix    []byte
	p2pkhPrefix   []byte
	upgradeParams []upgradeParams
	expiryHeight  uint32
}

type upgradeParams struct {
	ActivationHeight uint32
	BranchID         []byte
}

const (
	sighashMask                 = 0x1f
	blake2BSighash              = "ZcashSigHash"
	prevoutsHashPersonalization = "ZcashPrevoutHash"
	sequenceHashPersonalization = "ZcashSequencHash"
	outputsHashPersonalization  = "ZcashOutputsHash"

	versionOverwinter        int32  = 3
	versionOverwinterGroupID uint32 = 0x3C48270
	versionSapling                  = 4
	versionSaplingGroupID           = 0x892f2085
)

var (
	witnessMarkerBytes = []byte{0x00, 0x01}

	mainnet = netParams{
		name:   "mainnet",
		params: &chaincfg.MainNetParams,

		p2pkhPrefix: []byte{0x1C, 0xB8},
		p2shPrefix:  []byte{0x1C, 0xBD},
		upgradeParams: []upgradeParams{
			{0, []byte{0x00, 0x00, 0x00, 0x00}},
			{347500, []byte{0x19, 0x1B, 0xA8, 0x5B}},
			{419200, []byte{0xBB, 0x09, 0xB8, 0x76}},
			{653600, []byte{0x60, 0x0E, 0xB4, 0x2B}},
		},
		expiryHeight: 10000000,
	}
	testnet = netParams{
		name:   "testnet",
		params: &chaincfg.TestNet3Params,

		p2pkhPrefix: []byte{0x1D, 0x25},
		p2shPrefix:  []byte{0x1C, 0xBA},
		upgradeParams: []upgradeParams{
			{0, []byte{0x00, 0x00, 0x00, 0x00}},
			{207500, []byte{0x19, 0x1B, 0xA8, 0x5B}},
			{280000, []byte{0xBB, 0x09, 0xB8, 0x76}},
			{584000, []byte{0x60, 0x0E, 0xB4, 0x2B}},
		},
		expiryHeight: 10000000,
	}
	regnet = netParams{
		name:   "regtest",
		params: &chaincfg.RegressionNetParams,

		p2pkhPrefix: []byte{0x1D, 0x25},
		p2shPrefix:  []byte{0x1C, 0xBA},
		upgradeParams: []upgradeParams{
			{0, []byte{0x00, 0x00, 0x00, 0x00}},
			{60, []byte{0x19, 0x1B, 0xA8, 0x5B}},
			{80, []byte{0xBB, 0x09, 0xB8, 0x76}},
			{100, []byte{0x60, 0x0E, 0xB4, 0x2B}},
		},
		expiryHeight: 10000000,
	}
)

func calculateSighash(
	network netParams,
	subScript []byte,
	hashType txscript.SigHashType,
	tx *wire.MsgTx,
	idx int,
	amt int64,
	expiryHeight uint32,
) ([]byte, error) {
	sigHashes, err := txSighashes(tx)
	if err != nil {
		return nil, err
	}

	// As a sanity check, ensure the passed input index for the transaction
	// is valid.
	if idx > len(tx.TxIn)-1 {
		return nil, fmt.Errorf("blake2bSignatureHash error: idx %d but %d txins", idx, len(tx.TxIn))
	}

	// We'll utilize this buffer throughout to incrementally calculate
	// the signature hash for this transaction.
	var sigHash bytes.Buffer

	// << GetHeader
	// First write out, then encode the transaction's nVersion number. Zcash current nVersion = 3
	var bVersion [4]byte
	binary.LittleEndian.PutUint32(bVersion[:], uint32(tx.Version)|(1<<31))
	sigHash.Write(bVersion[:])

	var versionGroupID = versionOverwinterGroupID
	if tx.Version == versionSapling {
		versionGroupID = versionSaplingGroupID
	}

	// << nVersionGroupId
	// Version group ID
	var nVersion [4]byte
	binary.LittleEndian.PutUint32(nVersion[:], versionGroupID)
	sigHash.Write(nVersion[:])

	// Next write out the possibly pre-calculated hashes for the sequence
	// numbers of all inputs, and the hashes of the previous outs for all
	// outputs.
	var zeroHash chainhash.Hash

	// << hashPrevouts
	// If anyone can pay isn't active, then we can use the cached
	// hashPrevOuts, otherwise we just write zeroes for the prev outs.
	if hashType&txscript.SigHashAnyOneCanPay == 0 {
		sigHash.Write(sigHashes.HashPrevOuts[:])
	} else {
		sigHash.Write(zeroHash[:])
	}

	// << hashSequence
	// If the sighash isn't anyone can pay, single, or none, the use the
	// cached hash sequences, otherwise write all zeroes for the
	// hashSequence.
	if hashType&txscript.SigHashAnyOneCanPay == 0 &&
		hashType&sighashMask != txscript.SigHashSingle &&
		hashType&sighashMask != txscript.SigHashNone {
		sigHash.Write(sigHashes.HashSequence[:])
	} else {
		sigHash.Write(zeroHash[:])
	}

	// << hashOutputs
	// If the current signature mode isn't single, or none, then we can
	// re-use the pre-generated hashoutputs sighash fragment. Otherwise,
	// we'll serialize and add only the target output index to the signature
	// pre-image.
	if hashType&sighashMask != txscript.SigHashSingle && hashType&sighashMask != txscript.SigHashNone {
		sigHash.Write(sigHashes.HashOutputs[:])
	} else if hashType&sighashMask == txscript.SigHashSingle && idx < len(tx.TxOut) {
		var (
			b bytes.Buffer
			h chainhash.Hash
		)
		if err := wire.WriteTxOut(&b, 0, 0, tx.TxOut[idx]); err != nil {
			return nil, err
		}

		var err error
		if h, err = blake2b(b.Bytes(), []byte(outputsHashPersonalization)); err != nil {
			return nil, err
		}
		sigHash.Write(h.CloneBytes())
	} else {
		sigHash.Write(zeroHash[:])
	}

	// << hashJoinSplits
	sigHash.Write(zeroHash[:])

	// << hashShieldedSpends
	if tx.Version == versionSapling {
		sigHash.Write(zeroHash[:])
	}

	// << hashShieldedOutputs
	if tx.Version == versionSapling {
		sigHash.Write(zeroHash[:])
	}

	// << nLockTime
	var lockTime [4]byte
	binary.LittleEndian.PutUint32(lockTime[:], tx.LockTime)
	sigHash.Write(lockTime[:])

	// << nExpiryHeight
	var expiryTime [4]byte
	binary.LittleEndian.PutUint32(expiryTime[:], expiryHeight)
	sigHash.Write(expiryTime[:])

	// << valueBalance
	if tx.Version == versionSapling {
		var valueBalance [8]byte
		binary.LittleEndian.PutUint64(valueBalance[:], 0)
		sigHash.Write(valueBalance[:])
	}

	// << nHashType
	var bHashType [4]byte
	binary.LittleEndian.PutUint32(bHashType[:], uint32(hashType))
	sigHash.Write(bHashType[:])

	if idx != math.MaxUint32 {
		// << prevout
		// Next, write the outpoint being spent.
		sigHash.Write(tx.TxIn[idx].PreviousOutPoint.Hash[:])
		var bIndex [4]byte
		binary.LittleEndian.PutUint32(bIndex[:], tx.TxIn[idx].PreviousOutPoint.Index)
		sigHash.Write(bIndex[:])

		// << scriptCode
		// For p2wsh outputs, and future outputs, the script code is the
		// original script, with all code separators removed, serialized
		// with a var int length prefix.
		// wire.WriteVarBytes(&sigHash, 0, subScript)
		if err := wire.WriteVarBytes(&sigHash, 0, subScript); err != nil {
			return nil, err
		}

		// << amount
		// Next, add the input amount, and sequence number of the input being
		// signed.
		if err := binary.Write(&sigHash, binary.LittleEndian, amt); err != nil {
			return nil, err
		}

		// << nSequence
		var bSequence [4]byte
		binary.LittleEndian.PutUint32(bSequence[:], tx.TxIn[idx].Sequence)
		sigHash.Write(bSequence[:])
	}

	var h chainhash.Hash
	if h, err = blake2b(sigHash.Bytes(), sighashKey(expiryHeight, network)); err != nil {
		return nil, err
	}

	return h.CloneBytes(), nil
}

func blake2b(data, key []byte) (h chainhash.Hash, err error) {
	bHash := blake2.New(&blake2.Config{
		Size:     32,
		Personal: key,
	})

	if _, err = bHash.Write(data); err != nil {
		return h, err
	}

	err = (&h).SetBytes(bHash.Sum(nil))
	return h, err
}

func sighashKey(activationHeight uint32, network netParams) []byte {
	var i int
	upgradeParams := network.upgradeParams
	for i = len(upgradeParams) - 1; i >= 0; i-- {
		if activationHeight >= upgradeParams[i].ActivationHeight {
			break
		}
	}
	return append([]byte(blake2BSighash), upgradeParams[i].BranchID...)
}

// txSighashes computes, and returns the cached sighashes of the given
// transaction.
func txSighashes(tx *wire.MsgTx) (h *txscript.TxSigHashes, err error) {
	h = &txscript.TxSigHashes{}

	if h.HashPrevOuts, err = calculateHashPrevOuts(tx); err != nil {
		return
	}

	if h.HashSequence, err = calculateHashSequence(tx); err != nil {
		return
	}

	if h.HashOutputs, err = calculateHashOutputs(tx); err != nil {
		return
	}

	return
}

// calculateHashPrevOuts calculates a single hash of all the previous
// outputs (txid:index) referenced within the passed transaction. This
// calculated hash can be re-used when validating all inputs spending segwit
// outputs, with a signature hash type of SigHashAll. This allows validation to
// re-use previous hashing computation, reducing the complexity of validating
// SigHashAll inputs from  O(N^2) to O(N).
func calculateHashPrevOuts(tx *wire.MsgTx) (chainhash.Hash, error) {
	var b bytes.Buffer
	for _, in := range tx.TxIn {
		// First write out the 32-byte transaction ID one of whose outputs are
		// being referenced by this input.

		b.Write(in.PreviousOutPoint.Hash[:])

		// Next, we'll encode the index of the referenced output as a little
		// endian integer.
		var buf [4]byte
		binary.LittleEndian.PutUint32(buf[:], in.PreviousOutPoint.Index)
		b.Write(buf[:])
	}

	return blake2b(b.Bytes(), []byte(prevoutsHashPersonalization))
}

// calculateHashSequence computes an aggregated hash of each of the
// sequence numbers within the inputs of the passed transaction. This single
// hash can be re-used when validating all inputs spending segwit outputs, which
// include signatures using the SigHashAll sighash type. This allows validation
// to re-use previous hashing computation, reducing the complexity of validating
// SigHashAll inputs from O(N^2) to O(N).
func calculateHashSequence(tx *wire.MsgTx) (chainhash.Hash, error) {
	var b bytes.Buffer
	for _, in := range tx.TxIn {
		var buf [4]byte
		binary.LittleEndian.PutUint32(buf[:], in.Sequence)
		b.Write(buf[:])
	}

	return blake2b(b.Bytes(), []byte(sequenceHashPersonalization))
}

// calculateHashOutputs computes a hash digest of all outputs created by
// the transaction encoded using the wire format. This single hash can be
// re-used when validating all inputs spending witness programs, which include
// signatures using the SigHashAll sighash type. This allows computation to be
// cached, reducing the total hashing complexity from O(N^2) to O(N).
func calculateHashOutputs(tx *wire.MsgTx) (_ chainhash.Hash, err error) {
	var b bytes.Buffer
	for _, out := range tx.TxOut {
		if err = wire.WriteTxOut(&b, 0, 0, out); err != nil {
			return chainhash.Hash{}, err
		}
	}

	return blake2b(b.Bytes(), []byte(outputsHashPersonalization))
}

// writeTxOut encodes to into the bitcoin protocol encoding for a transaction
// output (TxOut) to w.
//
// NOTE: This function is exported in order to allow txscript to compute the
// new sighashes for witness transactions (BIP0143).
func writeTxOut(w io.Writer, pver uint32, version int32, to *wire.TxOut) error {
	if err := binary.Write(w, binary.LittleEndian, uint64(to.Value)); err != nil {
		return err
	}
	return writeVarBytes(w, pver, to.PkScript)
}

// writeTxIn encodes ti to the bitcoin protocol encoding for a transaction
// input (TxIn) to w.
func writeTxIn(w io.Writer, pver uint32, version int32, ti *wire.TxIn) error {
	err := writeOutPoint(w, pver, version, &ti.PreviousOutPoint)
	if err != nil {
		return err
	}

	err = writeVarBytes(w, pver, ti.SignatureScript)
	if err != nil {
		return err
	}

	return binary.Write(w, binary.LittleEndian, ti.Sequence)
}

// writeOutPoint encodes op to the bitcoin protocol encoding for an OutPoint
// to w.
func writeOutPoint(w io.Writer, pver uint32, version int32, op *wire.OutPoint) error {
	_, err := w.Write(op.Hash[:])
	if err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, op.Index)
}

// writeTxWitness encodes the bitcoin protocol encoding for a transaction
// input's witness into to w.
func writeTxWitness(w io.Writer, pver uint32, version int32, wit [][]byte) error {
	err := writeVarInt(w, pver, uint64(len(wit)))
	if err != nil {
		return err
	}
	for _, item := range wit {
		err = writeVarBytes(w, pver, item)
		if err != nil {
			return err
		}
	}
	return nil
}

// writeVarInt serializes val to w using a variable number of bytes depending
// on its value.
func writeVarInt(w io.Writer, pver uint32, val uint64) error {
	if val < 0xfd {
		return binary.Write(w, binary.LittleEndian, uint8(val))
	}

	if val <= math.MaxUint16 {
		err := binary.Write(w, binary.LittleEndian, 0xfd)
		if err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, uint16(val))
	}

	if val <= math.MaxUint32 {
		err := binary.Write(w, binary.LittleEndian, 0xfe)
		if err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, uint32(val))
	}

	if err := binary.Write(w, binary.LittleEndian, 0xff); err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, val)
}

// writeVarBytes serializes a variable length byte array to w as a varInt
// containing the number of bytes, followed by the bytes themselves.
func writeVarBytes(w io.Writer, pver uint32, bytes []byte) error {
	slen := uint64(len(bytes))
	err := writeVarInt(w, pver, slen)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}

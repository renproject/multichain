package bitcoincash

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/renproject/multichain/compat/bitcoincompat"
	"github.com/renproject/pack"
	"golang.org/x/crypto/ripemd160"
)

// Version of Bitcoin Cash transactions supported by the multichain.
const Version int32 = 1

type txBuilder struct {
	params *chaincfg.Params
}

// NewTxBuilder returns an implementation of the transaction builder interface
// from the Bitcoin Compat API, and exposes the functionality to build simple
// Bitcoin Cash transactions.
func NewTxBuilder(params *chaincfg.Params) bitcoincompat.TxBuilder {
	return txBuilder{params: params}
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
		addr, err := DecodeAddress(string(recipient.Address), txBuilder.params)
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

	return &Tx{inputs: inputs, recipients: recipients, msgTx: msgTx, signed: false}, nil
}

// Tx represents a simple Bitcoin Cash transaction that implements the Bitcoin
// Compat API.
type Tx struct {
	inputs     []bitcoincompat.Input
	recipients []bitcoincompat.Recipient

	msgTx *wire.MsgTx

	signed bool
}

func (tx *Tx) Hash() pack.Bytes32 {
	return pack.NewBytes32(tx.msgTx.TxHash())
}

func (tx *Tx) Sighashes() ([]pack.Bytes32, error) {
	sighashes := make([]pack.Bytes32, len(tx.inputs))
	for i, txin := range tx.inputs {
		pubKeyScript := txin.Output.PubKeyScript
		sigScript := txin.SigScript
		value := int64(txin.Output.Value.Uint64())
		if value < 0 {
			return []pack.Bytes32{}, fmt.Errorf("expected value >= 0, got value = %v", value)
		}

		var hash []byte
		if sigScript != nil {
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

func (tx *Tx) Serialize() (pack.Bytes, error) {
	buf := new(bytes.Buffer)
	if err := tx.msgTx.Serialize(buf); err != nil {
		return pack.Bytes{}, err
	}
	return pack.NewBytes(buf.Bytes()), nil
}

// An Address represents a Bitcoin Cash address.
type Address interface {
	btcutil.Address
	BitcoinCompatAddress() btcutil.Address
}

// AddressLegacy represents a legacy Bitcoin address.
type AddressLegacy struct {
	btcutil.Address
}

// BitcoinCompatAddress returns the address as if it was a Bitcoin address.
func (addr AddressLegacy) BitcoinCompatAddress() btcutil.Address {
	return addr.Address
}

// AddressPubKeyHash represents an address for P2PKH transactions for
// Bitcoin Cash that is compatible with the Bitcoin-compat API.
type AddressPubKeyHash struct {
	*btcutil.AddressPubKeyHash
	params *chaincfg.Params
}

// NewAddressPubKeyHash returns a new AddressPubKeyHash
// that is compatible with the Bitcoin-compat API.
func NewAddressPubKeyHash(pkh []byte, params *chaincfg.Params) (AddressPubKeyHash, error) {
	addr, err := btcutil.NewAddressPubKeyHash(pkh, params)
	return AddressPubKeyHash{AddressPubKeyHash: addr, params: params}, err
}

// NewAddressPubKey returns a new AddressPubKey
// that is compatible with the Bitcoin-compat API.
func NewAddressPubKey(pk []byte, params *chaincfg.Params) (AddressPubKeyHash, error) {
	return NewAddressPubKeyHash(btcutil.Hash160(pk), params)
}

// String returns the string encoding of the transaction output
// destination.
//
// Please note that String differs subtly from EncodeAddress: String
// will return the value as a string without any conversion, while
// EncodeAddress may convert destination types (for example,
// converting pubkeys to P2PKH addresses) before encoding as a
// payment address string.
func (addr AddressPubKeyHash) String() string {
	return addr.EncodeAddress()
}

// EncodeAddress returns the string encoding of the payment address
// associated with the Address value.  See the comment on String
// for how this method differs from String.
func (addr AddressPubKeyHash) EncodeAddress() string {
	hash := *addr.AddressPubKeyHash.Hash160()
	encoded, err := EncodeAddress(0x00, hash[:], addr.params)
	if err != nil {
		panic(fmt.Errorf("invalid address: %v", err))
	}
	return encoded
}

// ScriptAddress returns the raw bytes of the address to be used
// when inserting the address into a txout's script.
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

// AddressScriptHash represents an address for P2SH transactions for
// Bitcoin Cash that is compatible with the Bitcoin-compat API.
type AddressScriptHash struct {
	*btcutil.AddressScriptHash
	params *chaincfg.Params
}

// NewAddressScriptHash returns a new AddressScriptHash
// that is compatible with the Bitcoin-compat API.
func NewAddressScriptHash(pkh []byte, params *chaincfg.Params) (AddressScriptHash, error) {
	addr, err := btcutil.NewAddressScriptHash(pkh, params)
	return AddressScriptHash{AddressScriptHash: addr, params: params}, err
}

// String returns the string encoding of the transaction output
// destination.
//
// Please note that String differs subtly from EncodeAddress: String
// will return the value as a string without any conversion, while
// EncodeAddress may convert destination types (for example,
// converting pubkeys to P2PKH addresses) before encoding as a
// payment address string.
func (addr AddressScriptHash) String() string {
	return addr.EncodeAddress()
}

// EncodeAddress returns the string encoding of the payment address
// associated with the Address value.  See the comment on String
// for how this method differs from String.
func (addr AddressScriptHash) EncodeAddress() string {
	hash := *addr.AddressScriptHash.Hash160()
	encoded, err := EncodeAddress(0x00, hash[:], addr.params)
	if err != nil {
		panic(fmt.Errorf("invalid address: %v", err))
	}
	return encoded
}

// ScriptAddress returns the raw bytes of the address to be used
// when inserting the address into a txout's script.
func (addr AddressScriptHash) ScriptAddress() []byte {
	return addr.AddressScriptHash.ScriptAddress()
}

// IsForNet returns whether or not the address is associated with the passed
// bitcoin network.
func (addr AddressScriptHash) IsForNet(params *chaincfg.Params) bool {
	return addr.AddressScriptHash.IsForNet(params)
}

// BitcoinCompatAddress returns the address as if it was a Bitcoin address.
func (addr AddressScriptHash) BitcoinCompatAddress() btcutil.Address {
	return addr.AddressScriptHash
}

// Alphabet used by Bitcoin Cash to encode addresses.
var Alphabet = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

// AlphabetReverseLookup used by Bitcoin Cash to decode addresses.
var AlphabetReverseLookup = func() map[rune]byte {
	lookup := map[rune]byte{}
	for i, char := range Alphabet {
		lookup[char] = byte(i)
	}
	return lookup
}()

// SighashForkID used to distinguish between Bitcoin Cash and Bitcoin
// transactions by masking hash types.
var SighashForkID = txscript.SigHashType(0x40)

// SighashMask used to mask hash types.
var SighashMask = txscript.SigHashType(0x1F)

// EncodeAddress using Bitcoin Cash address encoding, assuming that the hash
// data has no prefix or checksum.
func EncodeAddress(version byte, hash []byte, params *chaincfg.Params) (string, error) {
	if (len(hash)-20)/4 != int(version)%8 {
		return "", fmt.Errorf("invalid version: %d", version)
	}
	data, err := bech32.ConvertBits(append([]byte{version}, hash...), 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("invalid bech32 encoding: %v", err)
	}
	return EncodeToString(AppendChecksum(AddressPrefix(params), data)), nil
}

func DecodeAddress(addr string, params *chaincfg.Params) (Address, error) {
	// Legacy address decoding
	if address, err := btcutil.DecodeAddress(addr, params); err == nil {
		switch address.(type) {
		case *btcutil.AddressPubKeyHash, *btcutil.AddressScriptHash, *btcutil.AddressPubKey:
			return AddressLegacy{Address: address}, nil
		case *btcutil.AddressWitnessPubKeyHash, *btcutil.AddressWitnessScriptHash:
			return nil, fmt.Errorf("unsuported segwit bitcoin address type %T", address)
		default:
			return nil, fmt.Errorf("unsuported legacy bitcoin address type %T", address)
		}
	}

	if addrParts := strings.Split(addr, ":"); len(addrParts) != 1 {
		addr = addrParts[1]
	}

	decoded := DecodeString(addr)
	if !VerifyChecksum(AddressPrefix(params), decoded) {
		return nil, btcutil.ErrChecksumMismatch
	}

	addrBytes, err := bech32.ConvertBits(decoded[:len(decoded)-8], 5, 8, false)
	if err != nil {
		return nil, err
	}

	switch len(addrBytes) - 1 {
	case ripemd160.Size: // P2PKH or P2SH
		switch addrBytes[0] {
		case 0: // P2PKH
			addr, err := btcutil.NewAddressPubKeyHash(addrBytes[1:21], params)
			if err != nil {
				return nil, err
			}
			return &AddressPubKeyHash{AddressPubKeyHash: addr, params: params}, nil
		case 8: // P2SH
			addr, err := btcutil.NewAddressScriptHash(addrBytes[1:21], params)
			if err != nil {
				return nil, err
			}
			return &AddressScriptHash{AddressScriptHash: addr, params: params}, nil
		default:
			return nil, btcutil.ErrUnknownAddressType
		}
	default:
		return nil, errors.New("decoded address is of unknown size")
	}
}

// EncodeToString using Bitcoin Cash address encoding, assuming that the data
// has a prefix and checksum.
func EncodeToString(data []byte) string {
	addr := strings.Builder{}
	for _, d := range data {
		addr.WriteByte(Alphabet[d])
	}
	return addr.String()
}

// DecodeString using Bitcoin Cash address encoding.
func DecodeString(address string) []byte {
	data := []byte{}
	for _, c := range address {
		data = append(data, AlphabetReverseLookup[c])
	}
	return data
}

// AppendChecksum to the data payload.
//
// https://github.com/bitcoincashorg/bitcoincash.org/blob/master/spec/cashaddr.md#checksum
func AppendChecksum(prefix string, payload []byte) []byte {
	prefixedPayload := append(EncodePrefix(prefix), payload...)

	// Append 8 zeroes.
	prefixedPayload = append(prefixedPayload, 0, 0, 0, 0, 0, 0, 0, 0)

	// Determine what to XOR into those 8 zeroes.
	mod := PolyMod(prefixedPayload)

	checksum := make([]byte, 8)
	for i := 0; i < 8; i++ {
		// Convert the 5-bit groups in mod to checksum values.
		checksum[i] = byte((mod >> uint(5*(7-i))) & 0x1f)
	}
	return append(payload, checksum...)
}

// VerifyChecksum verifies whether the given payload is well-formed.
//
// https://github.com/bitcoincashorg/bitcoincash.org/blob/master/spec/cashaddr.md#checksum
func VerifyChecksum(prefix string, payload []byte) bool {
	return PolyMod(append(EncodePrefix(prefix), payload...)) == 0
}

// EncodePrefix string into bytes.
//
// https://github.com/bitcoincashorg/bitcoincash.org/blob/master/spec/cashaddr.md#checksum
func EncodePrefix(prefixString string) []byte {
	prefixBytes := make([]byte, len(prefixString)+1)
	for i := 0; i < len(prefixString); i++ {
		prefixBytes[i] = byte(prefixString[i]) & 0x1f
	}
	prefixBytes[len(prefixString)] = 0
	return prefixBytes
}

// AddressPrefix returns the string representations of an address prefix based
// on the network parameters: "bitcoincash" (for mainnet), "bchtest" (for
// testnet), and "bchreg" (for regtest). This function panics if the network
// parameters are not recognised.
func AddressPrefix(params *chaincfg.Params) string {
	if params == nil {
		panic(fmt.Errorf("non-exhaustive pattern: params %v", params))
	}
	switch params {
	case &chaincfg.MainNetParams:
		return "bitcoincash"
	case &chaincfg.TestNet3Params:
		return "bchtest"
	case &chaincfg.RegressionNetParams:
		return "bchreg"
	default:
		panic(fmt.Errorf("non-exhaustive pattern: params %v", params.Name))
	}
}

// PolyMod is used to calculate the checksum for Bitcoin Cash
// addresses.
//
//  uint64_t PolyMod(const data &v) {
//      uint64_t c = 1;
//      for (uint8_t d : v) {
//          uint8_t c0 = c >> 35;
//          c = ((c & 0x07ffffffff) << 5) ^ d;
//          if (c0 & 0x01) c ^= 0x98f2bc8e61;
//          if (c0 & 0x02) c ^= 0x79b76d99e2;
//          if (c0 & 0x04) c ^= 0xf33e5fb3c4;
//          if (c0 & 0x08) c ^= 0xae2eabe2a8;
//          if (c0 & 0x10) c ^= 0x1e4f43e470;
//      }
//      return c ^ 1;
//  }
//
// https://github.com/bitcoincashorg/bitcoincash.org/blob/master/spec/cashaddr.md
func PolyMod(v []byte) uint64 {
	c := uint64(1)
	for _, d := range v {
		c0 := byte(c >> 35)
		c = ((c & 0x07ffffffff) << 5) ^ uint64(d)

		if c0&0x01 > 0 {
			c ^= 0x98f2bc8e61
		}
		if c0&0x02 > 0 {
			c ^= 0x79b76d99e2
		}
		if c0&0x04 > 0 {
			c ^= 0xf33e5fb3c4
		}
		if c0&0x08 > 0 {
			c ^= 0xae2eabe2a8
		}
		if c0&0x10 > 0 {
			c ^= 0x1e4f43e470
		}
	}
	return c ^ 1
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

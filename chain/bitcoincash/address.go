package bitcoincash

import (
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"
	"golang.org/x/crypto/ripemd160"
)

var (
	// Alphabet used by Bitcoin Cash to encode addresses.
	Alphabet = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
	// AlphabetReverseLookup used by Bitcoin Cash to decode addresses.
	AlphabetReverseLookup = func() map[rune]byte {
		lookup := map[rune]byte{}
		for i, char := range Alphabet {
			lookup[char] = byte(i)
		}
		return lookup
	}()
)

// AddressEncodeDecoder implements the address.EncodeDecoder interface
type AddressEncodeDecoder struct {
	AddressEncoder
	AddressDecoder
}

// NewAddressEncodeDecoder constructs a new AddressEncodeDecoder with the
// chain specific configurations
func NewAddressEncodeDecoder(params *chaincfg.Params) AddressEncodeDecoder {
	return AddressEncodeDecoder{
		AddressEncoder: NewAddressEncoder(params),
		AddressDecoder: NewAddressDecoder(params),
	}
}

// AddressEncoder encapsulates the chain specific configurations and implements
// the address.Encoder interface
type AddressEncoder struct {
	params *chaincfg.Params
}

// NewAddressEncoder constructs a new AddressEncoder with the chain specific
// configurations
func NewAddressEncoder(params *chaincfg.Params) AddressEncoder {
	return AddressEncoder{params: params}
}

// AddressDecoder encapsulates the chain specific configurations and implements
// the address.Decoder interface
type AddressDecoder struct {
	params *chaincfg.Params
}

// NewAddressDecoder constructs a new AddressDecoder with the chain specific
// configurations
func NewAddressDecoder(params *chaincfg.Params) AddressDecoder {
	return AddressDecoder{params: params}
}

// EncodeAddress implements the address.Encoder interface
func (encoder AddressEncoder) EncodeAddress(rawAddr address.RawAddress) (address.Address, error) {
	rawAddrBytes := []byte(rawAddr)
	var encodedAddr string
	var err error

	switch len(rawAddrBytes) - 1 {
	case ripemd160.Size: // P2PKH or P2SH
		switch rawAddrBytes[0] {
		case 0: // P2PKH
			encodedAddr, err = encodeAddress(0x00, rawAddrBytes[1:21], encoder.params)
		case 8: // P2SH
			encodedAddr, err = encodeAddress(8, rawAddrBytes[1:21], encoder.params)
		default:
			return address.Address(""), btcutil.ErrUnknownAddressType
		}
	default:
		return encodeLegacyAddress(rawAddr, encoder.params)
	}

	if err != nil {
		return address.Address(""), fmt.Errorf("encoding: %v", err)
	}

	return address.Address(encodedAddr), nil
}

// DecodeAddress implements the address.Decoder interface
func (decoder AddressDecoder) DecodeAddress(addr address.Address) (address.RawAddress, error) {
	// Legacy address decoding
	if legacyAddr, err := btcutil.DecodeAddress(string(addr), decoder.params); err == nil {
		switch legacyAddr.(type) {
		case *btcutil.AddressPubKeyHash, *btcutil.AddressScriptHash, *btcutil.AddressPubKey:
			return decodeLegacyAddress(addr, decoder.params)
		case *btcutil.AddressWitnessPubKeyHash, *btcutil.AddressWitnessScriptHash:
			return nil, fmt.Errorf("unsuported segwit bitcoin address type %T", legacyAddr)
		default:
			return nil, fmt.Errorf("unsuported legacy bitcoin address type %T", legacyAddr)
		}
	}

	if addrParts := strings.Split(string(addr), ":"); len(addrParts) != 1 {
		addr = address.Address(addrParts[1])
	}

	decoded := DecodeString(string(addr))
	if !VerifyChecksum(AddressPrefix(decoder.params), decoded) {
		return nil, btcutil.ErrChecksumMismatch
	}

	addrBytes, err := bech32.ConvertBits(decoded[:len(decoded)-8], 5, 8, false)
	if err != nil {
		return nil, err
	}

	switch len(addrBytes) - 1 {
	case ripemd160.Size: // P2PKH or P2SH
		switch addrBytes[0] {
		case 0, 8: // P2PKH or P2SH
			return address.RawAddress(addrBytes), nil
		default:
			return nil, btcutil.ErrUnknownAddressType
		}
	default:
		return nil, errors.New("decoded address is of unknown size")
	}
}

func encodeLegacyAddress(rawAddr address.RawAddress, params *chaincfg.Params) (address.Address, error) {
	// Validate that the base58 address is in fact in correct format.
	encodedAddr := base58.Encode([]byte(rawAddr))
	if _, err := btcutil.DecodeAddress(encodedAddr, &chaincfg.RegressionNetParams); err != nil {
		return address.Address(""), fmt.Errorf("address validation error: %v", err)
	}

	return address.Address(encodedAddr), nil
}

func decodeLegacyAddress(addr address.Address, params *chaincfg.Params) (address.RawAddress, error) {
	// Decode the checksummed base58 format address.
	decoded, ver, err := base58.CheckDecode(string(addr))
	if err != nil {
		return nil, fmt.Errorf("checking: %v", err)
	}
	if len(decoded) != 20 {
		return nil, fmt.Errorf("expected len 20, got len %v", len(decoded))
	}

	// Validate the address format.
	switch ver {
	case params.PubKeyHashAddrID, params.ScriptHashAddrID:
		return address.RawAddress(pack.NewBytes(base58.Decode(string(addr)))), nil
	default:
		return nil, fmt.Errorf("unexpected address prefix")
	}
}

// An Address represents a Bitcoin Cash address.
type Address interface {
	btcutil.Address
	BitcoinAddress() btcutil.Address
}

// AddressLegacy represents a legacy Bitcoin address.
type AddressLegacy struct {
	btcutil.Address
}

// BitcoinAddress returns the address as if it was a Bitcoin address.
func (addr AddressLegacy) BitcoinAddress() btcutil.Address {
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
	encoded, err := encodeAddress(0x00, hash[:], addr.params)
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

// BitcoinAddress returns the address as if it was a Bitcoin address.
func (addr AddressPubKeyHash) BitcoinAddress() btcutil.Address {
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
func NewAddressScriptHash(script []byte, params *chaincfg.Params) (AddressScriptHash, error) {
	addr, err := btcutil.NewAddressScriptHash(script, params)
	return AddressScriptHash{AddressScriptHash: addr, params: params}, err
}

// NewAddressScriptHashFromHash returns a new AddressScriptHash
// that is compatible with the Bitcoin-compat API.
func NewAddressScriptHashFromHash(scriptHash []byte, params *chaincfg.Params) (AddressScriptHash, error) {
	addr, err := btcutil.NewAddressScriptHashFromHash(scriptHash, params)
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
	encoded, err := encodeAddress(8, hash[:], addr.params)
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

// BitcoinAddress returns the address as if it was a Bitcoin address.
func (addr AddressScriptHash) BitcoinAddress() btcutil.Address {
	return addr.AddressScriptHash
}

// encodeAddress using Bitcoin Cash address encoding, assuming that the hash
// data has no prefix or checksum.
func encodeAddress(version byte, hash []byte, params *chaincfg.Params) (string, error) {
	if (len(hash)-20)/4 != int(version)%8 {
		return "", fmt.Errorf("invalid version: %d", version)
	}
	data, err := bech32.ConvertBits(append([]byte{version}, hash...), 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("invalid bech32 encoding: %v", err)
	}
	return EncodeToString(AppendChecksum(AddressPrefix(params), data)), nil
}

// addressFromRawBytes consumes raw bytes representation of a bitcoincash
// address and returns a type that implements the bitcoincash.Address interface.
func addressFromRawBytes(addrBytes []byte, params *chaincfg.Params) (Address, error) {
	switch len(addrBytes) - 1 {
	case ripemd160.Size: // P2PKH or P2SH
		switch addrBytes[0] {
		case 0: // P2PKH
			return NewAddressPubKeyHash(addrBytes[1:21], params)
		case 8: // P2SH
			return NewAddressScriptHashFromHash(addrBytes[1:21], params)
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

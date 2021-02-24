package zcash

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"
	"golang.org/x/crypto/ripemd160"
)

// AddressEncodeDecoder implements the address.EncodeDecoder interface
type AddressEncodeDecoder struct {
	AddressEncoder
	AddressDecoder
}

// AddressEncoder encapsulates the chain specific configurations and implements
// the address.Encoder interface
type AddressEncoder struct {
	params *Params
}

// AddressDecoder encapsulates the chain specific configurations and implements
// the address.Decoder interface
type AddressDecoder struct {
	params *Params
}

// NewAddressEncoder constructs a new AddressEncoder with the chain specific
// configurations
func NewAddressEncoder(params *Params) AddressEncoder {
	return AddressEncoder{params: params}
}

// NewAddressDecoder constructs a new AddressDecoder with the chain specific
// configurations
func NewAddressDecoder(params *Params) AddressDecoder {
	return AddressDecoder{params: params}
}

// NewAddressEncodeDecoder constructs a new AddressEncodeDecoder with the
// chain specific configurations
func NewAddressEncodeDecoder(params *Params) AddressEncodeDecoder {
	return AddressEncodeDecoder{
		AddressEncoder: NewAddressEncoder(params),
		AddressDecoder: NewAddressDecoder(params),
	}
}

// EncodeAddress implements the address.Encoder interface
func (encoder AddressEncoder) EncodeAddress(rawAddr address.RawAddress) (address.Address, error) {
	var addrType uint8
	var err error
	var hash [20]byte
	var prefix []byte

	switch len(rawAddr) {
	case ripemd160.Size + 5:
		prefix = rawAddr[:1]
		addrType, err = addressType(prefix, encoder.params)
		copy(hash[:], rawAddr[1:21])
	case ripemd160.Size + 6:
		prefix = rawAddr[:2]
		addrType, err = addressType(prefix, encoder.params)
		copy(hash[:], rawAddr[2:22])
	default:
		return address.Address(""), fmt.Errorf("validating address length: expected %v or %v, got %v", ripemd160.Size+5, ripemd160.Size+6, len(rawAddr))
	}

	if err != nil {
		return address.Address(""), fmt.Errorf("parsing address type: %v", err)
	}

	switch addrType {
	case 0, 1: // P2PKH or P2SH
		return address.Address(pack.String(encodeAddress(hash[:], prefix))), nil
	default:
		return address.Address(""), fmt.Errorf("unexpected address type: %v", addrType)
	}
}

// DecodeAddress implements the address.Decoder interface
func (decoder AddressDecoder) DecodeAddress(addr address.Address) (address.RawAddress, error) {
	var decoded = base58.Decode(string(addr))
	var addrType uint8
	var err error
	var hash [20]byte

	switch len(decoded) {
	case ripemd160.Size + 5:
		addrType, err = addressType(decoded[:1], decoder.params)
		copy(hash[:], decoded[1:21])
	case ripemd160.Size + 6:
		addrType, err = addressType(decoded[:2], decoder.params)
		copy(hash[:], decoded[2:22])
	default:
		return nil, fmt.Errorf("validating address length: expected %v or %v, got %v", ripemd160.Size+5, ripemd160.Size+6, len(decoded))
	}

	if err != nil {
		return nil, fmt.Errorf("parsing address type: %v", err)
	}

	var cksum [4]byte
	copy(cksum[:], decoded[len(decoded)-4:])
	if checksum(decoded[:len(decoded)-4]) != cksum {
		return nil, fmt.Errorf("validating checksum: %v", base58.ErrChecksum)
	}

	switch addrType {
	case 0, 1: // P2PKH or P2SH
		return address.RawAddress(pack.Bytes(decoded)), nil
	default:
		return nil, fmt.Errorf("unexpected address type: %v", addrType)
	}
}

func addressType(prefix []byte, params *Params) (uint8, error) {
	if bytes.Equal(prefix, params.P2PKHPrefix) {
		return 0, nil
	}
	if bytes.Equal(prefix, params.P2SHPrefix) {
		return 1, nil
	}
	return 0, btcutil.ErrUnknownAddressType
}

// An Address represents a Zcash address.
type Address interface {
	btcutil.Address
	BitcoinAddress() btcutil.Address
}

// AddressPubKeyHash represents an address for P2PKH transactions for Zcash that
// is compatible with the Bitcoin Compat API.
type AddressPubKeyHash struct {
	*btcutil.AddressPubKeyHash
	params *Params
}

// NewAddressPubKeyHash returns a new AddressPubKeyHash that is compatible with
// the Bitcoin Compat API.
func NewAddressPubKeyHash(pkh []byte, params *Params) (AddressPubKeyHash, error) {
	addr, err := btcutil.NewAddressPubKeyHash(pkh, params.Params)
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
	return encodeAddress(hash[:], addr.params.P2PKHPrefix)
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

// BitcoinAddress returns the address as if it was a Bitcoin address.
func (addr AddressPubKeyHash) BitcoinAddress() btcutil.Address {
	return addr.AddressPubKeyHash
}

// AddressScriptHash represents an address for P2SH transactions for Zcash that
// is compatible with the Bitcoin Compat API.
type AddressScriptHash struct {
	*btcutil.AddressScriptHash
	params *Params
}

// NewAddressScriptHash returns a new AddressScriptHash that is compatible with
// the Bitcoin Compat API.
func NewAddressScriptHash(script []byte, params *Params) (AddressScriptHash, error) {
	addr, err := btcutil.NewAddressScriptHash(script, params.Params)
	return AddressScriptHash{AddressScriptHash: addr, params: params}, err
}

// NewAddressScriptHashFromHash returns a new AddressScriptHash that is compatible with
// the Bitcoin Compat API.
func NewAddressScriptHashFromHash(scriptHash []byte, params *Params) (AddressScriptHash, error) {
	addr, err := btcutil.NewAddressScriptHashFromHash(scriptHash, params.Params)
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

// BitcoinAddress returns the address as if it was a Bitcoin address.
func (addr AddressScriptHash) BitcoinAddress() btcutil.Address {
	return addr.AddressScriptHash
}

// EncodeAddress returns the string encoding of the payment address associated
// with the Address value. See the comment on String for how this method differs
// from String.
func (addr AddressScriptHash) EncodeAddress() string {
	hash := *addr.AddressScriptHash.Hash160()
	return encodeAddress(hash[:], addr.params.P2SHPrefix)
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

// addressFromRawBytes decodes a string-representation of an address to an address
// type that implements the zcash.Address interface
func addressFromRawBytes(addrBytes []byte, params *Params) (Address, error) {
	var addrType uint8
	var err error
	var hash [20]byte

	switch len(addrBytes) {
	case ripemd160.Size + 5:
		addrType, err = addressType(addrBytes[:1], params)
		copy(hash[:], addrBytes[1:21])
	case ripemd160.Size + 6:
		addrType, err = addressType(addrBytes[:2], params)
		copy(hash[:], addrBytes[2:22])
	default:
		return nil, fmt.Errorf("validating address length: expected %v or %v, got %v", ripemd160.Size+5, ripemd160.Size+6, len(addrBytes))
	}

	if err != nil {
		return nil, fmt.Errorf("parsing address type: %v", err)
	}

	switch addrType {
	case 0: // P2PKH
		return NewAddressPubKeyHash(hash[:], params)
	case 1: // P2SH
		return NewAddressScriptHashFromHash(hash[:], params)
	default:
		return nil, fmt.Errorf("unexpected address type: %v", addrType)
	}
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

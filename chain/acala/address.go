package acala

import (
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/renproject/multichain/api/address"

	"golang.org/x/crypto/blake2b"
)

const (
	// The default address type byte used for a substrate chain.
	DefaultSubstrateWildcard = byte(42)
)

var (
	// Prefix used before hashing the address bytes for calculating checksum
	Prefix = []byte("SS58PRE")
)

// AddressDecoder implements the address.Decoder interface.
type AddressDecoder struct {
	addressType byte
}

// AddressEncoder implements the address.Encoder interface.
type AddressEncoder struct {
	addressType byte
}

// AddressEncodeDecoder implements the address.EncodeDecoder interface.
type AddressEncodeDecoder struct {
	AddressEncoder
	AddressDecoder
}

// NewAddressEncodeDecoder constructs a new AddressEncodeDecoder.
func NewAddressEncodeDecoder(addressType byte) AddressEncodeDecoder {
	return AddressEncodeDecoder{
		AddressEncoder: AddressEncoder{addressType},
		AddressDecoder: AddressDecoder{addressType},
	}
}

// DecodeAddress the string using the Bitcoin base58 alphabet. The substrate
// address is decoded and only the 32-byte public key is returned as the raw
// address.
func (decoder AddressDecoder) DecodeAddress(addr address.Address) (address.RawAddress, error) {
	data := base58.Decode(string(addr))
	if len(data) != 35 {
		return address.RawAddress([]byte{}), fmt.Errorf("expected 35 bytes, got %v bytes", len(data))
	}
	return address.RawAddress(data[1:33]), nil
}

// EncodeAddress the raw bytes using the Bitcoin base58 alphabet. We expect a
// 32-byte substrate public key as the address in its raw bytes representation.
// A checksum encoded key is then encoded in the base58 format.
func (encoder AddressEncoder) EncodeAddress(rawAddr address.RawAddress) (address.Address, error) {
	if len(rawAddr) != 32 {
		return address.Address(""), fmt.Errorf("expected 32 bytes, got %v bytes", len(rawAddr))
	}
	checksummedAddr := append([]byte{encoder.addressType}, rawAddr...)
	checksum := blake2b.Sum512(append(Prefix, checksummedAddr...))

	checksummedAddr = append(checksummedAddr, checksum[0:2]...)

	return address.Address(base58.Encode(checksummedAddr)), nil
}

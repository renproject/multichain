package acala

import (
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/renproject/multichain"

	"golang.org/x/crypto/blake2b"
)

const (
	// AddressTypeDefault is the default address type byte for a substrate chain.
	AddressTypeDefault = byte(42)
	// AddressTypeTestnet is the address type used for testnet.
	AddressTypeTestnet = byte(42)
	// AddressTypeCanaryNetwork is the address type used for canary network.
	AddressTypeCanaryNetwork = byte(8)
	// AddressTypeMainnet is the address type used for mainnet.
	AddressTypeMainnet = byte(10)
)

var (
	// Prefix used before hashing the address bytes for calculating checksum
	Prefix = []byte("SS58PRE")
)

// GetAddressType returns the appropriate prefix address type for a network
// type.
func GetAddressType(network multichain.Network) byte {
	switch network {
	case multichain.NetworkLocalnet, multichain.NetworkDevnet:
		return AddressTypeDefault
	case multichain.NetworkTestnet:
		return AddressTypeTestnet
	case multichain.NetworkMainnet:
		return AddressTypeMainnet
	default:
		return AddressTypeDefault
	}
}

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
func (decoder AddressDecoder) DecodeAddress(addr multichain.Address) (multichain.RawAddress, error) {
	data := base58.Decode(string(addr))
	if len(data) != 35 {
		return multichain.RawAddress([]byte{}), fmt.Errorf("expected 35 bytes, got %v bytes", len(data))
	}
	return multichain.RawAddress(data[1:33]), nil
}

// EncodeAddress the raw bytes using the Bitcoin base58 alphabet. We expect a
// 32-byte substrate public key as the address in its raw bytes representation.
// A checksum encoded key is then encoded in the base58 format.
func (encoder AddressEncoder) EncodeAddress(rawAddr multichain.RawAddress) (multichain.Address, error) {
	if len(rawAddr) != 32 {
		return multichain.Address(""), fmt.Errorf("expected 32 bytes, got %v bytes", len(rawAddr))
	}
	checksummedAddr := append([]byte{encoder.addressType}, rawAddr...)
	checksum := blake2b.Sum512(append(Prefix, checksummedAddr...))

	checksummedAddr = append(checksummedAddr, checksum[0:2]...)

	return multichain.Address(base58.Encode(checksummedAddr)), nil
}

package acala

import (
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/renproject/multichain/api/address"
)

// AddressDecoder implements the address.Decoder interface.
type AddressDecoder struct{}

// AddressEncoder implements the address.Encoder interface.
type AddressEncoder struct{}

// AddressEncodeDecoder implements the address.EncodeDecoder interface.
type AddressEncodeDecoder struct {
	AddressEncoder
	AddressDecoder
}

// NewAddressEncodeDecoder constructs a new AddressEncodeDecoder.
func NewAddressEncodeDecoder() AddressEncodeDecoder {
	return AddressEncodeDecoder{
		AddressEncoder: AddressEncoder{},
		AddressDecoder: AddressDecoder{},
	}
}

// DecodeAddress the string using the Bitcoin base58 alphabet. If the string
// does not a 2-byte address type, 32-byte array, and 1-byte checksum, then an
// error is returned.
func (AddressDecoder) DecodeAddress(addr address.Address) (address.RawAddress, error) {
	data := base58.Decode(string(addr))
	if len(data) != 35 {
		return address.RawAddress([]byte{}), fmt.Errorf("expected 35 bytes, got %v bytes", len(data))
	}
	return address.RawAddress(data), nil
}

// EncodeAddress the raw bytes using the Bitcoin base58 alphabet. If the data to
// encode is not a 2-byte address type, 32-byte array, and 1-byte checksum, then
// an error is returned.
func (AddressEncoder) EncodeAddress(rawAddr address.RawAddress) (address.Address, error) {
	if len(rawAddr) != 35 {
		return address.Address(""), fmt.Errorf("expected 35 bytes, got %v bytes", len(rawAddr))
	}
	return address.Address(base58.Encode(rawAddr)), nil
}

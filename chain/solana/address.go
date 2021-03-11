package solana

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

// NewAddressEncodeDecoder constructs and returns a new AddressEncodeDecoder.
func NewAddressEncodeDecoder() AddressEncodeDecoder {
	return AddressEncodeDecoder{
		AddressEncoder: AddressEncoder{},
		AddressDecoder: AddressDecoder{},
	}
}

// EncodeAddress consumes a raw byte-representation of an address and encodes it
// to the human-readable Base58 format.
func (AddressEncoder) EncodeAddress(rawAddress address.RawAddress) (address.Address, error) {
	if len(rawAddress) != 32 {
		return address.Address(""), fmt.Errorf("expected address length 32, got address length %v", len(rawAddress))
	}
	return address.Address(base58.Encode(rawAddress)), nil
}

// DecodeAddress consumes a human-readable Base58 format and decodes it into a
// raw byte-representation.
func (AddressDecoder) DecodeAddress(encoded address.Address) (address.RawAddress, error) {
	decoded := base58.Decode(string(encoded))
	if len(decoded) != 32 {
		return nil, fmt.Errorf("expected address length 32, got address length %v", len(decoded))
	}
	return address.RawAddress(decoded), nil
}

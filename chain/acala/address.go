package acala

import (
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/renproject/pack"
)

// An Address represents a public address on a Substrate blockchain. It can be
// the address of an external account, or the address of a smart contract.
type Address pack.Bytes

// The AddressDecoder defines an interface for decoding string representations
// of Substrate address into the concrete Address type.
type AddressDecoder interface {
	DecodeAddress(pack.String) (Address, error)
}

type addressDecoder struct{}

// NewAddressDecoder returns the default AddressDecoder for Substract chains. It
// uses the Bitcoin base58 alphabet to decode the string, and interprets the
// result as a 2-byte address type, 32-byte array, and 1-byte checksum.
func NewAddressDecoder() AddressDecoder {
	return addressDecoder{}
}

// DecodeAddress the string using the Bitcoin base58 alphabet. If the string
// does not a 2-byte address type, 32-byte array, and 1-byte checksum, then an
// error is returned.
func (addressDecoder) DecodeAddress(encoded pack.String) (Address, error) {
	data := base58.Decode(encoded.String())
	if len(data) != 35 {
		return Address{}, fmt.Errorf("expected 35 bytes, got %v bytes", len(data))
	}
	return Address(data), nil
}

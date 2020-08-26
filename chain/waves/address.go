package waves

import (
	"github.com/renproject/multichain/api/address"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

type AddressDecoder struct{}

// NewAddressDecoder returns the default AddressDecoder for Waves chain. It
// uses the base58 alphabet to decode the string, and interprets the
// result as a 26-byte array.
func NewAddressDecoder() AddressDecoder {
	return AddressDecoder{}
}

// DecodeAddress decodes from string into bytes.
func (AddressDecoder) DecodeAddress(encoded address.Address) (address.RawAddress, error) {
	addr, err := proto.NewAddressFromString(string(encoded))
	if err != nil {
		return nil, err
	}
	return address.RawAddress(addr.Bytes()), nil
}

type AddressEncoder struct{}

// DecodeAddress decodes from bytes into string.
func (AddressEncoder) EncodeAddress(decoded address.RawAddress) (address.Address, error) {
	addr, err := proto.NewAddressFromBytes(decoded)
	if err != nil {
		return "", err
	}
	return address.Address(addr.String()), nil
}

type AddressEncodeDecoder struct {
	AddressDecoder
	AddressEncoder
}

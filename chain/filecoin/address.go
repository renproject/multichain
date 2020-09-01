package filecoin

import (
	filaddress "github.com/filecoin-project/go-address"
	"github.com/renproject/multichain/api/address"
)

// AddressEncodeDecoder implements the address.EncodeDecoder interface
type AddressEncodeDecoder struct {
	AddressEncoder
	AddressDecoder
}

// AddressEncoder implements the address.Encoder interface.
type AddressEncoder struct{}

// AddressDecoder implements the address.Decoder interface.
type AddressDecoder struct{}

// NewAddressEncodeDecoder constructs a new AddressEncodeDecoder.
func NewAddressEncodeDecoder() AddressEncodeDecoder {
	return AddressEncodeDecoder{
		AddressEncoder: AddressEncoder{},
		AddressDecoder: AddressDecoder{},
	}
}

// EncodeAddress implements the address.Encoder interface. It receives a raw
// address and encodes it to a human-readable stringified address.
func (encoder AddressEncoder) EncodeAddress(raw address.RawAddress) (address.Address, error) {
	addr, err := filaddress.NewFromBytes([]byte(raw))
	if err != nil {
		return address.Address(""), err
	}
	return address.Address(addr.String()), nil
}

// DecodeAddress implements the address.Decoder interface. It receives a human
// readable address and decodes it to an address represented by raw bytes.
func (addrDecoder AddressDecoder) DecodeAddress(addr address.Address) (address.RawAddress, error) {
	rawAddr, err := filaddress.NewFromString(string(addr))
	if err != nil {
		return nil, err
	}
	return address.RawAddress(rawAddr.Bytes()), nil
}

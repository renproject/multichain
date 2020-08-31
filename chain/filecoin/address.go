package filecoin

import (
	filaddress "github.com/filecoin-project/go-address"
	"github.com/renproject/multichain/api/address"
)

type AddressEncodeDecoder struct {
	AddressEncoder
	AddressDecoder
}

type AddressEncoder struct{}
type AddressDecoder struct{}

func NewAddressEncodeDecoder() AddressEncodeDecoder {
	return AddressEncodeDecoder{
		AddressEncoder: AddressEncoder{},
		AddressDecoder: AddressDecoder{},
	}
}

func (encoder AddressEncoder) EncodeAddress(raw address.RawAddress) (address.Address, error) {
	addr, err := filaddress.NewFromBytes([]byte(raw))
	if err != nil {
		return address.Address(""), err
	}
	return address.Address(addr.String()), nil
}

func (addrDecoder AddressDecoder) DecodeAddress(addr address.Address) (address.RawAddress, error) {
	rawAddr, err := filaddress.NewFromString(string(addr))
	if err != nil {
		return nil, err
	}
	return address.RawAddress(rawAddr.Bytes()), nil
}

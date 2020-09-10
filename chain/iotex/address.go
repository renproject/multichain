package iotex

import (
	iotexaddr "github.com/iotexproject/iotex-address/address"

	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"
)

type AddressEncodeDecoder struct {
	AddressEncoder
	AddressDecoder
}

type AddressEncoder interface {
	EncodeAddress(address.RawAddress) (address.Address, error)
}

type addressEncoder struct{}

func NewAddressEncodeDecoder() address.EncodeDecoder {
	return AddressEncodeDecoder{
		AddressEncoder: NewAddressEncoder(),
		AddressDecoder: NewAddressDecoder(),
	}
}

type AddressDecoder interface {
	DecodeAddress(address.Address) (address.RawAddress, error)
}

type addressDecoder struct{}

func NewAddressDecoder() AddressDecoder {
	return addressDecoder{}
}

func NewAddressEncoder() AddressEncoder {
	return addressEncoder{}
}

func (addressDecoder) DecodeAddress(encoded address.Address) (address.RawAddress, error) {
	addr, err := iotexaddr.FromString(string(encoded))
	if err != nil {
		return nil, err
	}
	return address.RawAddress(pack.Bytes(addr.Bytes())), nil
}

func (addressEncoder) EncodeAddress(rawAddr address.RawAddress) (address.Address, error) {
	addr, err := iotexaddr.FromBytes(rawAddr)
	if err != nil {
		return address.Address(pack.NewString("")), err
	}
	return address.Address(pack.NewString(addr.String())), nil
}

package peercoin

import (
	"github.com/ppcsuite/btcutil"
	"github.com/ppcsuite/btcutil/base58"
	"github.com/ppcsuite/ppcd/chaincfg"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"
)

type AddressEncodeDecoder struct {
	AddressEncoder
	AddressDecoder
}

func NewAddressEncodeDecoder(params *chaincfg.Params) AddressEncodeDecoder {
	return AddressEncodeDecoder{
		AddressEncoder: NewAddressEncoder(params),
		AddressDecoder: NewAddressDecoder(params),
	}
}

type AddressEncoder struct {
	params *chaincfg.Params
}

func NewAddressEncoder(params *chaincfg.Params) AddressEncoder {
	return AddressEncoder{params: params}
}

func (encoder AddressEncoder) EncodeAddress(rawAddr address.RawAddress) (address.Address, error) {
	encodedAddr := base58.Encode([]byte(rawAddr))
	if _, err := btcutil.DecodeAddress(encodedAddr, encoder.params); err != nil {
		// Check that the address is valid.
		return address.Address(""), err
	}
	return address.Address(encodedAddr), nil
}

type AddressDecoder struct {
	params *chaincfg.Params
}

func NewAddressDecoder(params *chaincfg.Params) AddressDecoder {
	return AddressDecoder{params: params}
}

func (decoder AddressDecoder) DecodeAddress(addr address.Address) (pack.Bytes, error) {
	if _, err := btcutil.DecodeAddress(string(addr), decoder.params); err != nil {
		// Check that the address is valid.
		return nil, err
	}
	return pack.NewBytes(base58.Decode(string(addr))), nil
}

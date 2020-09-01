package bitcoin

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"
)

// AddressEncodeDecoder implements the address.EncodeDecoder interface
type AddressEncodeDecoder struct {
	AddressEncoder
	AddressDecoder
}

// NewAddressEncodeDecoder constructs a new AddressEncodeDecoder with the
// chain specific configurations
func NewAddressEncodeDecoder(params *chaincfg.Params) AddressEncodeDecoder {
	return AddressEncodeDecoder{
		AddressEncoder: NewAddressEncoder(params),
		AddressDecoder: NewAddressDecoder(params),
	}
}

// AddressEncoder encapsulates the chain specific configurations and implements
// the address.Encoder interface
type AddressEncoder struct {
	params *chaincfg.Params
}

// NewAddressEncoder constructs a new AddressEncoder with the chain specific
// configurations
func NewAddressEncoder(params *chaincfg.Params) AddressEncoder {
	return AddressEncoder{params: params}
}

// EncodeAddress implements the address.Encoder interface
func (encoder AddressEncoder) EncodeAddress(rawAddr address.RawAddress) (address.Address, error) {
	encodedAddr := base58.Encode([]byte(rawAddr))
	if _, err := btcutil.DecodeAddress(encodedAddr, encoder.params); err != nil {
		// Check that the address is valid.
		return address.Address(""), err
	}
	return address.Address(encodedAddr), nil
}

// AddressDecoder encapsulates the chain specific configurations and implements
// the address.Decoder interface
type AddressDecoder struct {
	params *chaincfg.Params
}

// NewAddressDecoder constructs a new AddressDecoder with the chain specific
// configurations
func NewAddressDecoder(params *chaincfg.Params) AddressDecoder {
	return AddressDecoder{params: params}
}

// DecodeAddress implements the address.Decoder interface
func (decoder AddressDecoder) DecodeAddress(addr address.Address) (pack.Bytes, error) {
	if _, err := btcutil.DecodeAddress(string(addr), decoder.params); err != nil {
		// Check that the address is valid.
		return nil, err
	}
	return pack.NewBytes(base58.Decode(string(addr))), nil
}

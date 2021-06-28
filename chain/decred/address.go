package decred

import (
	"fmt"

	"github.com/decred/dcrd/chaincfg/v3"
	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/decred/base58"
	"github.com/renproject/multichain/api/address"
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

	// Validate that the base58 address is in fact in correct format.
	encodedAddr := base58.Encode([]byte(rawAddr))
	if _, err := dcrutil.DecodeAddress(encodedAddr, encoder.params); err != nil {
		return address.Address(""), err
	}

	return address.Address(encodedAddr), nil
	//return address.Address(""), nil
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
func (decoder AddressDecoder) DecodeAddress(addr address.Address) (address.RawAddress, error) {
	decodedAddr, err := dcrutil.DecodeAddress(string(addr), decoder.params)
	if err != nil {
		return nil, fmt.Errorf("decode address: %v", err)
	}

	switch a := decodedAddr.(type) {
	case *dcrutil.AddressPubKeyHash, *dcrutil.AddressScriptHash:
		return address.RawAddress(base58.Decode(string(addr))), nil
	default:
		return nil, fmt.Errorf("non-exhaustive pattern: address %T", a)
	}
}

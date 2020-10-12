package bitcoin

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
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
	if _, err := btcutil.DecodeAddress(encodedAddr, encoder.params); err != nil {
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
func (decoder AddressDecoder) DecodeAddress(addr address.Address) (address.RawAddress, error) {
	// Decode the checksummed base58 format address.
	decoded, ver, err := base58.CheckDecode(string(addr))
	if err != nil {
		return nil, fmt.Errorf("checking: %v", err)
	}
	if len(decoded) != 20 {
		return nil, fmt.Errorf("expected len 20, got len %v", len(decoded))
	}

	// Validate the address format.
	switch ver {
	case decoder.params.PubKeyHashAddrID, decoder.params.ScriptHashAddrID:
		return address.RawAddress(base58.Decode(string(addr))), nil
	default:
		return nil, fmt.Errorf("unexpected address prefix")
	}
}

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
	switch len(rawAddr) {
	case 25:
		return encoder.encodeBase58(rawAddr)
	case 21, 33:
		return encoder.encodeBech32(rawAddr)
	default:
		return address.Address(""), fmt.Errorf("non-exhaustive pattern: address length %v", len(rawAddr))
	}
}

func (encoder AddressEncoder) encodeBase58(rawAddr address.RawAddress) (address.Address, error) {
	// Validate that the base58 address is in fact in correct format.
	encodedAddr := base58.Encode([]byte(rawAddr))
	if _, err := btcutil.DecodeAddress(encodedAddr, encoder.params); err != nil {
		return address.Address(""), err
	}

	return address.Address(encodedAddr), nil
}

func (encoder AddressEncoder) encodeBech32(rawAddr address.RawAddress) (address.Address, error) {
	switch len(rawAddr) {
	case 21:
		addr, err := btcutil.NewAddressWitnessPubKeyHash(rawAddr[1:], encoder.params)
		if err != nil {
			return address.Address(""), fmt.Errorf("new address witness pubkey hash: %v", err)
		}
		return address.Address(addr.EncodeAddress()), nil
	case 33:
		addr, err := btcutil.NewAddressWitnessScriptHash(rawAddr[1:], encoder.params)
		if err != nil {
			return address.Address(""), fmt.Errorf("new address witness script hash: %v", err)
		}
		return address.Address(addr.EncodeAddress()), nil
	default:
		return address.Address(""), fmt.Errorf("non-exhaustive pattern: bech32 address length %v", len(rawAddr))
	}
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
	decodedAddr, err := btcutil.DecodeAddress(string(addr), decoder.params)
	if err != nil {
		return nil, fmt.Errorf("decode address: %v", err)
	}

	switch a := decodedAddr.(type) {
	case *btcutil.AddressPubKeyHash, *btcutil.AddressScriptHash:
		return address.RawAddress(base58.Decode(string(addr))), nil
	case *btcutil.AddressWitnessPubKeyHash:
		rawAddr := append([]byte{a.WitnessVersion()}, a.WitnessProgram()...)
		return address.RawAddress(rawAddr), nil
	case *btcutil.AddressWitnessScriptHash:
		rawAddr := append([]byte{a.WitnessVersion()}, a.WitnessProgram()...)
		return address.RawAddress(rawAddr), nil
	default:
		return nil, fmt.Errorf("non-exhaustive pattern: address %T", a)
	}
}

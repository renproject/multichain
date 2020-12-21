package bitcoin

import (
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/renproject/multichain/api/address"
)

// AddressEncodeDecoder implements the address.EncodeDecoder interface
type AddressEncodeDecoder struct {
	AddressEncoder
	AddressDecoder
}

// NewAddressEncodeDecoder constructs a new AddressEncodeDecoder with the
// chain specific configurations
func NewAddressEncodeDecoder(params *chaincfg.Params, hrp string) AddressEncodeDecoder {
	return AddressEncodeDecoder{
		AddressEncoder: NewAddressEncoder(params, hrp),
		AddressDecoder: NewAddressDecoder(params),
	}
}

// AddressEncoder encapsulates the chain specific configurations and implements
// the address.Encoder interface
type AddressEncoder struct {
	params *chaincfg.Params
	hrp    string
}

// NewAddressEncoder constructs a new AddressEncoder with the chain specific
// configurations
func NewAddressEncoder(params *chaincfg.Params, hrp string) AddressEncoder {
	return AddressEncoder{params: params, hrp: hrp}
}

// EncodeAddress implements the address.Encoder interface
func (encoder AddressEncoder) EncodeAddress(rawAddr address.RawAddress) (address.Address, error) {
	switch len(rawAddr) {
	case 25:
		return encoder.encodeBase58(rawAddr)
	case 21:
		return encoder.encodeBech32(rawAddr)
	default:
		return address.Address(""), fmt.Errorf("non-exhaustive pattern: raw address length %v", len(rawAddr))
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
	addrBase32, err := bech32.ConvertBits([]byte(rawAddr), 8, 5, false)
	if err != nil {
		return address.Address(""), fmt.Errorf("convert base: %v", err)
	}

	addr, err := bech32.Encode(encoder.hrp, addrBase32)
	if err != nil {
		return address.Address(""), fmt.Errorf("encode bech32: %v", err)
	}

	return address.Address(addr), nil
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
	if strings.HasPrefix(string(addr), "1") ||
		strings.HasPrefix(string(addr), "2") ||
		strings.HasPrefix(string(addr), "3") ||
		strings.HasPrefix(string(addr), "m") ||
		strings.HasPrefix(string(addr), "n") {
		return decoder.decodeBase58(addr)
	} else if strings.HasPrefix(string(addr), "bc") ||
		strings.HasPrefix(string(addr), "tb") {
		return decoder.decodeBech32(addr)
	} else {
		return nil, fmt.Errorf("non-exhaustive pattern: address %v", addr)
	}
}

func (decoder AddressDecoder) decodeBase58(addr address.Address) (address.RawAddress, error) {
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

func (decoder AddressDecoder) decodeBech32(addr address.Address) (address.RawAddress, error) {
	_, dataBase32, err := bech32.Decode(string(addr))
	if err != nil {
		return nil, fmt.Errorf("decoding: %v", err)
	}

	rawAddr, err := bech32.ConvertBits(dataBase32, 5, 8, true)
	if err != nil {
		return nil, fmt.Errorf("convert base: %v", err)
	}

	return address.RawAddress(rawAddr), nil
}

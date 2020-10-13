package cosmos

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/renproject/multichain/api/address"
)

// An Address is a public address that can be encoded/decoded to/from strings.
// Addresses are usually formatted different between different network
// configurations.
type Address sdk.AccAddress

// AccAddress convert Address to sdk.AccAddress
func (addr Address) AccAddress() sdk.AccAddress {
	return sdk.AccAddress(addr)
}

// String implements the Stringer interface
func (addr Address) String() string {
	return sdk.AccAddress(addr).String()
}

// AddressEncodeDecoder encapsulates fields that implement the
// address.EncodeDecoder interface
type AddressEncodeDecoder struct {
	AddressEncoder
	AddressDecoder
}

// NewAddressEncodeDecoder creates a new address encoder-decoder
func NewAddressEncodeDecoder(hrp string) AddressEncodeDecoder {
	return AddressEncodeDecoder{
		AddressEncoder: NewAddressEncoder(hrp),
		AddressDecoder: NewAddressDecoder(hrp),
	}
}

// AddressEncoder implements the address.Encoder interface
type AddressEncoder struct {
	hrp string
}

// AddressDecoder implements the address.Decoder interface
type AddressDecoder struct {
	hrp string
}

// NewAddressDecoder creates a new address decoder
func NewAddressDecoder(hrp string) AddressDecoder {
	return AddressDecoder{hrp: hrp}
}

// NewAddressEncoder creates a new address encoder
func NewAddressEncoder(hrp string) AddressEncoder {
	return AddressEncoder{hrp: hrp}
}

// DecodeAddress consumes a human-readable representation of a cosmos
// compatible address and decodes it to its raw bytes representation.
func (decoder AddressDecoder) DecodeAddress(addr address.Address) (address.RawAddress, error) {
	sdk.GetConfig().SetBech32PrefixForAccount(decoder.hrp, decoder.hrp+"pub")
	rawAddr, err := sdk.AccAddressFromBech32(string(addr))
	if err != nil {
		return nil, err
	}
	return address.RawAddress(rawAddr), nil
}

// EncodeAddress consumes raw bytes and encodes them to a human-readable
// address format.
func (encoder AddressEncoder) EncodeAddress(rawAddr address.RawAddress) (address.Address, error) {
	sdk.GetConfig().SetBech32PrefixForAccount(encoder.hrp, encoder.hrp+"pub")
	bech32Addr := sdk.AccAddress(rawAddr)
	return address.Address(bech32Addr.String()), nil
}

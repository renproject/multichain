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
func NewAddressEncodeDecoder() AddressEncodeDecoder {
	return AddressEncodeDecoder{
		AddressEncoder: NewAddressEncoder(),
		AddressDecoder: NewAddressDecoder(),
	}
}

// AddressEncoder implements the address.Encoder interface
type AddressEncoder struct {
}

// AddressDecoder implements the address.Decoder interface
type AddressDecoder struct {
}

// NewAddressDecoder creates a new address decoder
func NewAddressDecoder() AddressDecoder {
	return AddressDecoder{}
}

// NewAddressEncoder creates a new address encoder
func NewAddressEncoder() AddressEncoder {
	return AddressEncoder{}
}

// DecodeAddress consumes a human-readable representation of a cosmos
// compatible address and decodes it to its raw bytes representation.
func (decoder AddressDecoder) DecodeAddress(addr address.Address) (address.RawAddress, error) {
	rawAddr, err := sdk.AccAddressFromBech32(string(addr))
	if err != nil {
		return nil, err
	}
	return address.RawAddress(rawAddr), nil
}

// EncodeAddress consumes raw bytes and encodes them to a human-readable
// address format.
func (encoder AddressEncoder) EncodeAddress(rawAddr address.RawAddress) (address.Address, error) {
	if err := sdk.VerifyAddressFormat(rawAddr); err != nil {
		return address.Address(""), err
	}
	bech32Addr := sdk.AccAddress(rawAddr)
	return address.Address(bech32Addr.String()), nil
}

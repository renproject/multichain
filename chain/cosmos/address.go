package cosmos

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/renproject/multichain/compat/cosmoscompat"
	"github.com/renproject/pack"
)

type addressDecoder struct {
	HRP string
}

// NewAddressDecoder returns an implementation of the address decoder interface
// from the Cosmos Compat API, and exposes the functionality to decode strings
// into addresses.
func NewAddressDecoder(HRP string) cosmoscompat.AddressDecoder {
	return addressDecoder{HRP: HRP}
}

func (decoder addressDecoder) DecodeAddress(encoded pack.String) (cosmoscompat.Address, error) {
	sdk.GetConfig().SetBech32PrefixForAccount(decoder.HRP, decoder.HRP+"pub")
	addr, err := sdk.AccAddressFromBech32(encoded.String())
	if err != nil {
		return nil, err
	}

	return cosmoscompat.Address(addr), nil
}

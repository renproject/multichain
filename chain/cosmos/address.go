package cosmos

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/renproject/multichain/api/address"
)

type AddressDecoder struct {
	hrp string
}

func NewAddressDecoder(hrp string) AddressDecoder {
	return AddressDecoder{hrp: hrp}
}

func (decoder AddressDecoder) DecodeAddress(addr address.Address) (address.RawAddress, error) {
	types.GetConfig().SetBech32PrefixForAccount(decoder.hrp, decoder.hrp+"pub")
	rawAddr, err := types.AccAddressFromBech32(string(addr))
	if err != nil {
		return nil, err
	}
	return address.RawAddress(rawAddr), nil
}

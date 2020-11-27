package solana

import (
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/renproject/multichain/api/address"
)

type AddressDecoder struct{}

func NewAddressDecoder() AddressDecoder {
	return AddressDecoder{}
}

func (AddressDecoder) DecodeAddress(encoded address.Address) (address.RawAddress, error) {
	decoded := base58.Decode(string(encoded))
	if len(decoded) != 32 {
		return nil, fmt.Errorf("expected address length 32, got address length %v", len(decoded))
	}
	return address.RawAddress(decoded), nil
}

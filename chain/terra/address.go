package terra

import (
	"github.com/renproject/multichain/chain/cosmos"
	"github.com/renproject/multichain/compat/cosmoscompat"
)

// NewAddressDecoder returns an implementation of the address decoder interface
// from the Cosmos Compat API, and exposes the functionality to decode strings
// into addresses.
func NewAddressDecoder() cosmoscompat.AddressDecoder {
	return cosmos.NewAddressDecoder("terra")
}

package bitcoin

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/renproject/multichain/compat/bitcoincompat"
	"github.com/renproject/pack"
)

type addressDecoder struct {
	defaultNet *chaincfg.Params
}

// NewAddressDecoder returns an implementation of the address decoder interface
// from the Bitcoin Compat API, and exposes the functionality to decode strings
// into addresses.
func NewAddressDecoder(defaultNet *chaincfg.Params) bitcoincompat.AddressDecoder {
	return addressDecoder{defaultNet: defaultNet}
}

func (addressDecoder addressDecoder) DecodeAddress(encoded pack.String) (bitcoincompat.Address, error) {
	return btcutil.DecodeAddress(encoded.String(), addressDecoder.defaultNet)
}

package digibyte

import (
	"fmt"
	
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/renproject/multichain/compat/bitcoincompat"
	"github.com/renproject/pack"
)

type addressDecoder struct {
	defaultNet *chaincfg.Params
}

func DigiByteConfig(params *chaincfg.Params) *chaincfg.Params {
	if params == nil {
		panic(fmt.Errorf("non-exhaustive pattern: params %v", params))
	}

	switch params {
	case &chaincfg.MainNetParams:
		return DigiByteMainNetParams
	default:
		panic(fmt.Errorf("non-exhaustive pattern: params %v", params.Name))
	}
}

// NewAddressDecoder returns an implementation of the address decoder interface
// from the Bitcoin Compat API, and exposes the functionality to decode strings
// into addresses.
func NewAddressDecoder(defaultNet *chaincfg.Params) bitcoincompat.AddressDecoder {
	var coinConfig = DigiByteConfig(defaultNet)
	return addressDecoder{defaultNet: coinConfig}
}

func (addressDecoder addressDecoder) DecodeAddress(encoded pack.String) (bitcoincompat.Address, error) {
	return btcutil.DecodeAddress(encoded.String(), addressDecoder.defaultNet)
}

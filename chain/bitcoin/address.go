package bitcoin

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/renproject/pack"
)

type AddressDecoder struct {
	params *chaincfg.Params
}

func NewAddressDecoder(params *chaincfg.Params) AddressDecoder {
	return AddressDecoder{params: params}
}

func (addrDecoder AddressDecoder) DecodeAddress(encoded pack.String) (pack.Bytes, error) {
	addr, err := btcutil.DecodeAddress(string(encoded), addrDecoder.params)
	if err != nil {
		return nil, err
	}
	decoded := base58.Decode(addr.EncodeAddress())
	return pack.Bytes(decoded), nil
}

package solana

import (
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/renproject/pack"
)

type AddressDecoder struct{}

func NewAddressDecoder() AddressDecoder {
	return AddressDecoder{}
}

func (AddressDecoder) DecodeAddress(encoded pack.String) (pack.Bytes, error) {
	decoded := base58.Decode(string(encoded))
	if len(decoded) != 32 {
		return nil, fmt.Errorf("expected address length 32, got address length %v", len(decoded))
	}
	return pack.Bytes(decoded), nil
}

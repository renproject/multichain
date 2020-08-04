package ontology

import (
	"fmt"

	"github.com/ontio/ontology/common"
	"github.com/renproject/pack"
)

type AddressDecoder struct{}

func NewAddressDecoder() AddressDecoder {
	return AddressDecoder{}
}

func (AddressDecoder) DecodeAddress(encoded pack.String) (pack.Bytes, error) {
	address, err := common.AddressFromBase58(string(encoded))
	if err != nil {
		return nil, fmt.Errorf("decode ontology base58 adress error %v", err)
	}
	return pack.Bytes(address[:]), nil
}

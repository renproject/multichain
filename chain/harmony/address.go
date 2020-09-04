package harmony

import (
	"github.com/btcsuite/btcutil/bech32"
	"github.com/renproject/multichain/api/address"
)

type EncoderDecoder struct {
	address.Encoder
	address.Decoder
}

func NewEncoderDecoder() address.EncodeDecoder {
	return EncoderDecoder{
		Encoder: NewEncoder(),
		Decoder: NewDecoder(),
	}
}

type Encoder struct {}

func (Encoder) EncodeAddress(addr address.RawAddress) (address.Address, error) {
	encodedAddr, err := bech32.ConvertBits(addr, 8, 5, true)
	if err != nil {
		return nil, err
	}
	return address.Address(encodedAddr), nil
}

type Decoder struct {}

func (Decoder) DecodeAddress(addr address.Address) (address.RawAddress, error) {
	_, decodedAddr, err := bech32.Decode(string(addr))
	if err != nil {
		return nil, err
	}
	return decodedAddr[:], nil
}

func NewEncoder() address.Encoder {
	return Encoder{}
}

func NewDecoder() address.Decoder {
	return Decoder{}
}
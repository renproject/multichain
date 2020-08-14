package address

import "github.com/renproject/pack"

type Address pack.String

func (addr Address) SizeHint() int {
	return pack.String(addr).SizeHint()
}

func (addr Address) Marshal(buf []byte, rem int) ([]byte, int, error) {
	return pack.String(addr).Marshal(buf, rem)
}

func (addr *Address) Unmarshal(buf []byte, rem int) ([]byte, int, error) {
	return (*pack.String)(addr).Unmarshal(buf, rem)
}

type RawAddress pack.Bytes

func (addr RawAddress) SizeHint() int {
	return pack.Bytes(addr).SizeHint()
}

func (addr RawAddress) Marshal(buf []byte, rem int) ([]byte, int, error) {
	return pack.Bytes(addr).Marshal(buf, rem)
}

func (addr *RawAddress) Unmarshal(buf []byte, rem int) ([]byte, int, error) {
	return (*pack.Bytes)(addr).Unmarshal(buf, rem)
}

type Encoder interface {
	EncodeAddress(RawAddress) (Address, error)
}

type Decoder interface {
	DecodeAddress(Address) (RawAddress, error)
}

type EncodeDecoder interface {
	Encoder
	Decoder
}

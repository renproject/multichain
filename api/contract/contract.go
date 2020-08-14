package contract

import (
	"context"

	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"
)

type CallData pack.Bytes

func (data CallData) SizeHint() int {
	return pack.Bytes(data).SizeHint()
}

func (data CallData) Marshal(buf []byte, rem int) ([]byte, int, error) {
	return pack.Bytes(data).Marshal(buf, rem)
}

func (data *CallData) Unmarshal(buf []byte, rem int) ([]byte, int, error) {
	return (*pack.Bytes)(data).Unmarshal(buf, rem)
}

type Caller interface {
	CallContract(context.Context, address.Address, CallData) (pack.Bytes, error)
}

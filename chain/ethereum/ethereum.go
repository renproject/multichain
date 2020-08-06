package ethereum

import (
	"context"

	"github.com/renproject/pack"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (client *Client) CallContract(ctx context.Context, contract pack.String, input pack.Value, outputType pack.Type) (out pack.Value, err error) {
	// inputEthargs, inputEthvals, err := Encode(input)
	// if err != nil {
	// 	return pack.Bytes(nil), fmt.Errorf("bad input: %v", err)
	// }

	// inputEncoded, err := inputEthargs.Pack(inputEthvals...)
	// if err != nil {
	// 	return pack.Bytes(nil), fmt.Errorf("bad encoding: %v", err)
	// }

	// contractAddress, err := ethereumcompat.NewAddressFromHex(contract.String())
	// if err != nil {
	// 	return pack.Bytes(nil), fmt.Errorf("bad contract address %v: %v", contract.String(), err)
	// }

	panic("unimplemented")
}

package harmony

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/harmony-one/harmony/rpc"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"
)

type Params struct {
	CallArgs rpc.CallArgs
	Block    uint64
}

func (c *Client) CallContract(ctx context.Context, addr address.Address, callData contract.CallData) (pack.Bytes, error) {
	const method = "hmyv2_call"
	// Unmarshal required to get the block number parameter for the call
	var callParams Params
	err := json.Unmarshal(callData, &callParams)
	if err != nil {
		return nil, err
	}
	args, err := json.Marshal(callParams.CallArgs)
	if err != nil {
		return nil, err
	}
	data := []byte(fmt.Sprintf("[%s, %d]", string(args), callParams.Block))
	response, err := SendData(method, data, c.opts.Host)
	if err != nil {
		return nil, err
	}
	return pack.NewBytes(*response.Result), nil
}
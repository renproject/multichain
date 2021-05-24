package harmony

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"
)

type CallArgs struct {
	From     *common.Address `json:"from"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Data     *hexutil.Bytes  `json:"data"`
}
type Params struct {
	CallArgs CallArgs
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
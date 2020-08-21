package solana

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"
	"go.uber.org/zap"
)

const DefaultClientRPCURL = "http://localhost:8899"

type ClientOptions struct {
	Logger *zap.Logger
	RPCURL string
}

func DefaultClientOptions() ClientOptions {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return ClientOptions{
		Logger: logger,
		RPCURL: DefaultClientRPCURL,
	}
}

type Client struct {
	opts ClientOptions
}

func NewClient(opts ClientOptions) *Client {
	return &Client{opts: opts}
}

func (client *Client) CallContract(ctx context.Context, contract address.Address, input pack.Bytes) (output pack.Bytes, err error) {
	if input != nil && len(input) != 0 {
		return nil, fmt.Errorf("expected nil input, got %v input", input)
	}

	// Make an RPC call to "getAccountInfo" to get the data associated with the
	// account (we interpret the contract address as the account identifier).
	params, err := json.Marshal(string(pack.String(contract)))
	if err != nil {
		return pack.Bytes(nil), fmt.Errorf("encoding params: %v", err)
	}
	res, err := SendDataWithRetry("getAccountInfo", params, client.opts.RPCURL)
	if err != nil {
		return pack.Bytes(nil), fmt.Errorf("calling rpc method \"getAccountInfo\": %v", err)
	}
	if res.Result == nil {
		return pack.Bytes(nil), fmt.Errorf("decoding result: empty")
	}

	// Decode the data associated with the account into pack-encoded bytes.
	info := ResponseGetAccountInfo{}
	if err := json.Unmarshal(*res.Result, &info); err != nil {
		return pack.Bytes(nil), fmt.Errorf("decoding result: %v", err)
	}
	fmt.Printf("account data: %v", info.Value.Data)

	data := base58.Decode(info.Value.Data)
	// data, err := base64.StdEncoding.DecodeString()
	if err != nil {
		return pack.Bytes(nil), fmt.Errorf("decoding result from base58: %v", err)
	}
	return pack.NewBytes(data), nil
}

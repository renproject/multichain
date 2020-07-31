package solana

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/renproject/pack"
	"go.uber.org/zap"
)

type ClientOptions struct {
	Logger *zap.Logger
	URL    string
}

type Client struct {
	opts ClientOptions
}

func NewClient(opts ClientOptions) *Client {
	return &Client{opts: opts}
}

func (client *Client) CallContract(ctx context.Context, contract pack.String, input pack.Value, outputType pack.Type) (output pack.Value, err error) {
	if input.Type().Kind() != pack.KindBytes {
		return pack.Bytes(nil), fmt.Errorf("unsupported input kind %v", input.Type().Kind())
	}
	if outputType.Kind() != pack.KindBytes {
		return pack.Bytes(nil), fmt.Errorf("unsupported output kind %v", outputType.Kind())
	}

	// Make an RPC call to "getAccountInfo" to get the data associated with the
	// account (we interpret the contract address as the account identifier).
	params, err := json.Marshal([]string{contract.String()})
	if err != nil {
		return pack.Bytes(nil), fmt.Errorf("encoding params: %v", err)
	}
	res, err := SendDataWithRetry("getAccountInfo", params, client.opts.URL)
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

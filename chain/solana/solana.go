package solana

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"
	"go.uber.org/zap"
)

// DefaultClientRPCURL is the default RPC URL for the Solana cluster.
const DefaultClientRPCURL = "http://localhost:8899"

// ClientOptions define the options to instantiate a new Solana client.
type ClientOptions struct {
	Logger *zap.Logger
	RPCURL string
}

// DefaultClientOptions return the client options used to instantiate a Solana
// client by default.
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

// WithRPCURL returns a modified version of the options with the given API
// rpc-url
func (opts ClientOptions) WithRPCURL(rpcURL pack.String) ClientOptions {
	opts.RPCURL = string(rpcURL)
	return opts
}

// Client represents a Solana client that implements the multichain Contract API.
type Client struct {
	opts ClientOptions
}

// NewClient returns a new solana.Client interface that implements the
// multichain Contract API.
func NewClient(opts ClientOptions) *Client {
	return &Client{opts: opts}
}

// FindProgramAddress is a wrapper function that calls the Solana FFI to find
// the deterministic program-derived address using the program and seeds.
func FindProgramAddress(seeds []byte, program address.RawAddress) (address.Address, error) {
	addrEncodeDecoder := NewAddressEncodeDecoder()
	encoded, err := addrEncodeDecoder.EncodeAddress(program)
	if err != nil {
		return address.Address(""), err
	}

	return ProgramDerivedAddress(seeds, encoded), nil
}

// GetAccountInfo fetches and returns the account data.
func (client *Client) GetAccountData(account address.Address) (pack.Bytes, error) {
	// Fetch account info with base64 encoding. The default base58 encoding does
	// not support account data that is larger than 128 bytes, hence base64.
	params := json.RawMessage(fmt.Sprintf(`["%v", {"encoding":"base64"}]`, string(account)))
	res, err := SendDataWithRetry("getAccountInfo", params, client.opts.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("calling rpc method \"getAccountInfo\": %v", err)
	}
	if res.Result == nil {
		return nil, fmt.Errorf("decoding result: empty")
	}

	// Deserialise the account's info into the appropriate struct.
	info := ResponseGetAccountInfo{}
	if err := json.Unmarshal(*res.Result, &info); err != nil {
		return nil, fmt.Errorf("decoding result: %v", err)
	}

	// Decode the Base58 encoded account data into raw byte-representation. Since
	// this holds the burn log's data.
	data, err := base64.RawStdEncoding.DecodeString(info.Value.Data[0])
	if err != nil {
		return nil, fmt.Errorf("decoding base64 value: %v", err)
	}

	return pack.Bytes(data), nil
}

// CallContract implements the multichain Contract API. In the case of Solana,
// it is used to fetch burn logs associated with a particular burn nonce.
func (client *Client) CallContract(
	ctx context.Context,
	program address.Address,
	calldata contract.CallData,
) (pack.Bytes, error) {
	addrEncodeDecoder := NewAddressEncodeDecoder()
	decodedProgram, err := addrEncodeDecoder.DecodeAddress(program)
	if err != nil {
		return pack.Bytes(nil), fmt.Errorf("decode address: %v", err)
	}

	// Find the program-derived address that will have persisted the burn log.
	burnLogAccount, err := FindProgramAddress([]byte(calldata), decodedProgram)
	if err != nil {
		return pack.Bytes(nil), fmt.Errorf("find program-derived address: %v", err)
	}

	// Make an RPC call to "getAccountInfo" to get the data associated with the
	// account (we interpret the contract address as the account identifier).
	params := json.RawMessage(fmt.Sprintf(`["%v", {"encoding":"base58"}]`, string(burnLogAccount)))
	res, err := SendDataWithRetry("getAccountInfo", params, client.opts.RPCURL)
	if err != nil {
		return pack.Bytes(nil), fmt.Errorf("calling rpc method \"getAccountInfo\": %v", err)
	}
	if res.Result == nil {
		return pack.Bytes(nil), fmt.Errorf("decoding result: empty")
	}

	// Deserialise the account's info into the appropriate struct.
	info := ResponseGetAccountInfo{}
	if err := json.Unmarshal(*res.Result, &info); err != nil {
		return pack.Bytes(nil), fmt.Errorf("decoding result: %v", err)
	}

	// Decode the Base58 encoded account data into raw byte-representation. Since
	// this holds the burn log's data.
	data := base58.Decode(info.Value.Data[0])
	if err != nil {
		return pack.Bytes(nil), fmt.Errorf("decoding result from base58: %v", err)
	}

	return pack.NewBytes(data), nil
}

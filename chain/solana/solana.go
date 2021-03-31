package solana

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/near/borsh-go"
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

// GetGatewayBySelectorHash queries the gateway registry account on Solana for
// the deployed gateway address of the specific RenVM selector hash.
func (client *Client) GetGatewayBySelectorHash(registryProgram address.Address, shash pack.Bytes32) (address.Address, error) {
	gateways, err := client.GetGateways(registryProgram)
	if err != nil {
		return address.Address(""), err
	}
	gateway, ok := gateways[shash]
	if !ok {
		return address.Address(""), fmt.Errorf("gateway registry does not contain selector=%v", shash)
	}
	return gateway, nil
}

// GetGateways queries the gateway registry account on Solana to fetch all
// the gateway addresses and returns a map of RenVM selector hashes to their
// gateway addresses.
func (client *Client) GetGateways(registryProgram address.Address) (map[pack.Bytes32]address.Address, error) {
	registry, err := client.getGatewayRegistry(registryProgram)
	if err != nil {
		return nil, fmt.Errorf("get gateway registry: %v", err)
	}

	gateways := make(map[pack.Bytes32]address.Address)
	addrEncodeDecoder := NewAddressEncodeDecoder()
	for i := 0; i < int(registry.Count); i++ {
		selector := registry.Selectors[i]
		gateway, err := addrEncodeDecoder.EncodeAddress(address.RawAddress(registry.Gateways[i][:]))
		if err != nil {
			return nil, fmt.Errorf("encode address: %v", err)
		}
		gateways[selector] = gateway
	}

	return gateways, nil
}

func (client *Client) getGatewayRegistry(registryProgram address.Address) (GatewayRegistry, error) {
	seeds := []byte("GatewayRegistryState")
	programDerivedAddress := ProgramDerivedAddress(pack.Bytes(seeds), registryProgram)

	// Fetch account info with base64 encoding. The default base58 encoding does
	// not support account data that is larger than 128 bytes, hence base64.
	params := json.RawMessage(fmt.Sprintf(`["%v", {"encoding":"base64"}]`, string(programDerivedAddress)))
	res, err := SendDataWithRetry("getAccountInfo", params, client.opts.RPCURL)
	if err != nil {
		return GatewayRegistry{}, fmt.Errorf("calling rpc method \"getAccountInfo\": %v", err)
	}
	if res.Result == nil {
		return GatewayRegistry{}, fmt.Errorf("decoding result: empty")
	}

	// Deserialise the account's info into the appropriate struct.
	info := ResponseGetAccountInfo{}
	if err := json.Unmarshal(*res.Result, &info); err != nil {
		return GatewayRegistry{}, fmt.Errorf("decoding result: %v", err)
	}

	// Decode the Base58 encoded account data into raw byte-representation. Since
	// this holds the burn log's data.
	data, err := base64.RawStdEncoding.DecodeString(info.Value.Data[0])
	if err != nil {
		return GatewayRegistry{}, fmt.Errorf("decoding base64 value: %v", err)
	}

	registry := GatewayRegistry{}
	err = borsh.Deserialize(&registry, data)
	if err != nil {
		return GatewayRegistry{}, fmt.Errorf("deserializing gateway registry data: %v", err)
	}

	return registry, nil
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

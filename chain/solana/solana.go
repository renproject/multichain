package solana

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"
	"github.com/renproject/surge"
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

// BurnCallContractInput represents the data structure of the expected calldata
// for Solana's Contract API implementation.
type BurnCallContractInput struct {
	Nonce pack.U64
}

// BurnCallContractOutput represents the data structure for Solana Contract
// implementation's output bytes.
type BurnCallContractOutput struct {
	Amount    pack.U64
	Recipient address.RawAddress
	Confs     pack.U64
	Payload   pack.Bytes
}

// CallContract implements the multichain Contract API. In the case of Solana,
// it is used to fetch burn logs associated with a particular burn nonce.
func (client *Client) CallContract(
	ctx context.Context,
	program address.Address,
	calldata contract.CallData,
) (pack.Bytes, error) {
	// Deserialise the calldata bytes.
	input := BurnCallContractInput{}
	if err := surge.FromBinary(&input, calldata); err != nil {
		return pack.Bytes{}, fmt.Errorf("deserialise calldata: %v", err)
	}

	addrEncodeDecoder := NewAddressEncodeDecoder()
	decodedProgram, err := addrEncodeDecoder.DecodeAddress(program)
	if err != nil {
		return pack.Bytes(nil), fmt.Errorf("decode address: %v", err)
	}

	// little endian serialization of burn count (nonce)
	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, uint64(input.Nonce))

	// Find the program-derived address that will have persisted the burn log.
	burnLogAccount, err := FindProgramAddress(nonceBytes, decodedProgram)
	if err != nil {
		return pack.Bytes(nil), fmt.Errorf("find program-derived address: %v", err)
	}

	// Make an RPC call to "getAccountInfo" to get the data associated with the
	// account (we interpret the contract address as the account identifier).
	params, err := json.Marshal([]string{string(burnLogAccount)})
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

	// Deserialise the account's info into the appropriate struct.
	info := ResponseGetAccountInfo{}
	if err := json.Unmarshal(*res.Result, &info); err != nil {
		return pack.Bytes(nil), fmt.Errorf("decoding result: %v", err)
	}

	// Decode the Base58 encoded account data into raw byte-representation. Since
	// this holds the burn log's data, it should be 33 bytes in size. 8 bytes for
	// the burn amount and 25 bytes for the recipient of the released assets.
	data := base58.Decode(info.Value.Data)
	if err != nil {
		return pack.Bytes(nil), fmt.Errorf("decoding result from base58: %v", err)
	}
	if len(data) != 33 {
		return pack.Bytes(nil), fmt.Errorf("unexpected data length. Expected 33, got: %v", len(data))
	}

	// Serialize the burn log's data into raw bytes.
	amount := binary.LittleEndian.Uint64(data[:8])
	recipient := pack.Bytes(data[8:])
	burnOut := BurnCallContractOutput{
		Amount:    pack.NewU64(amount),
		Recipient: address.RawAddress(recipient),
		Confs:     pack.NewU64(0),
		Payload:   pack.Bytes(nil),
	}
	burnOutBytes, err := surge.ToBinary(burnOut)
	if err != nil {
		return pack.Bytes(nil), fmt.Errorf("serialize burn log output: %v", err)
	}

	return pack.NewBytes(burnOutBytes), nil
}

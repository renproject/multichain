package acala

import (
	"context"
	"fmt"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"
	"github.com/renproject/surge"
	"go.uber.org/zap"
)

const (
	DefaultClientRPCURL = "ws://127.0.0.1:9944"
)

type ClientOptions struct {
	Logger *zap.Logger
	rpcURL pack.String
}

func DefaultClientOptions() ClientOptions {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return ClientOptions{
		Logger: logger,
		rpcURL: DefaultClientRPCURL,
	}
}

func (opts ClientOptions) WithRPCURL(rpcURL pack.String) ClientOptions {
	opts.rpcURL = rpcURL
	return opts
}

type Client struct {
	opts ClientOptions
	api  gsrpc.SubstrateAPI
}

func NewClient(opts ClientOptions) (*Client, error) {
	substrateAPI, err := gsrpc.NewSubstrateAPI(string(opts.rpcURL))
	if err != nil {
		return nil, err
	}

	return &Client{
		opts: opts,
		api:  *substrateAPI,
	}, nil
}

type BurnLogInput struct {
	Blockhash pack.Bytes32
	ExtSign   pack.Bytes
}

type BurnLogOutput struct {
	Amount    pack.U256
	Recipient address.RawAddress
	Confs     pack.U64
}

func (client *Client) CallContractSystemEvents(_ context.Context, _ address.Address, calldata contract.CallData) (pack.Bytes, error) {
	// Deserialise the calldata bytes.
	input := BurnLogInput{}
	if err := surge.FromBinary(&input, calldata); err != nil {
		return pack.Bytes{}, fmt.Errorf("deserialise calldata: %v\n", err)
	}

	// Get chain metadata.
	meta, err := client.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return pack.Bytes{}, fmt.Errorf("get metadata: %v", err)
	}

	// This key is used to read the state storage at the block of interest.
	key, err := types.CreateStorageKey(meta, "System", "Events", nil, nil)
	if err != nil {
		return pack.Bytes{}, fmt.Errorf("create storage key: %v", err)
	}

	// Get the block in which the burn event was logged.
	block, err := client.api.RPC.Chain.GetBlock(types.Hash(input.Blockhash))
	if err != nil {
		return pack.Bytes{}, fmt.Errorf("get block: %v", err)
	}

	// Get the latest block header. This will be used to calculate number of block
	// confirmations of the burn log of interest.
	header, err := client.api.RPC.Chain.GetHeaderLatest()
	if err != nil {
		return pack.Bytes{}, fmt.Errorf("get header: %v", err)
	}

	// Retrieve raw bytes from storage at the block and storage key of interest.
	data, err := client.api.RPC.State.GetStorageRaw(key, types.Hash(input.Blockhash))
	if err != nil {
		return pack.Bytes{}, fmt.Errorf("get storage: %v", err)
	}

	// Fetch the extrinsic's index in the block.
	extID := -1
	for i, ext := range block.Block.Extrinsics {
		if input.ExtSign.Equal(pack.Bytes(ext.Signature.Signature.AsSr25519[:])) {
			extID = i
			break
		}
	}
	if extID == -1 {
		return pack.Bytes{}, fmt.Errorf("extrinsic not found in block")
	}

	// Decode the event data to get the burn log.
	burnEvent, err := decodeEventData(meta, data, uint32(extID))
	if err != nil {
		return pack.Bytes{}, err
	}

	// Calculate block confirmations for the event.
	confs := header.Number - block.Block.Header.Number + 1

	burnLogOutput := BurnLogOutput{
		Amount:    pack.NewU256FromInt(burnEvent.Amount.Int),
		Recipient: address.RawAddress(burnEvent.Dest[:]),
		Confs:     pack.NewU64(uint64(confs)),
	}

	out, err := surge.ToBinary(burnLogOutput)
	if err != nil {
		return pack.Bytes{}, fmt.Errorf("serialise output: %v", err)
	}

	return pack.Bytes(out), nil
}

func decodeEventData(meta *types.Metadata, data *types.StorageDataRaw, id uint32) (eventBurnt, error) {
	events := RenVmBridgeEvents{}
	if err := types.EventRecordsRaw(*data).DecodeEventRecords(meta, &events); err != nil {
		return eventBurnt{}, fmt.Errorf("decode event data: %v", err)
	}

	// Match the event to the appropriate extrinsic index.
	for _, event := range events.RenVmBridge_Burnt {
		if event.Phase.AsApplyExtrinsic == id {
			return event, nil
		}
	}

	return eventBurnt{}, fmt.Errorf("burn event not found")
}

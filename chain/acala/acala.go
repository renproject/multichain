package acala

import (
	"fmt"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"
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

func (client *Client) BurnEvent(blockhash pack.Bytes32) (pack.U256, address.RawAddress, pack.U64, error) {
	// Get chain metadata.
	meta, err := client.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return pack.U256{}, nil, pack.U64(uint64(0)), fmt.Errorf("get metadata: %v", err)
	}

	// This key is used to read the state storage at the block of interest.
	key, err := types.CreateStorageKey(meta, "System", "Events", nil, nil)
	if err != nil {
		return pack.U256{}, nil, pack.U64(uint64(0)), fmt.Errorf("create storage key: %v", err)
	}

	// Get the block in which the burn event was logged.
	block, err := client.api.RPC.Chain.GetBlock(types.Hash(blockhash))
	if err != nil {
		return pack.U256{}, nil, pack.U64(uint64(0)), fmt.Errorf("get block: %v", err)
	}

	// Get the latest block header. This will be used to calculate number of block
	// confirmations of the burn log of interest.
	header, err := client.api.RPC.Chain.GetHeaderLatest()
	if err != nil {
		return pack.U256{}, nil, pack.U64(uint64(0)), fmt.Errorf("get header: %v", err)
	}

	// Retrieve raw bytes from storage at the block and storage key of interest.
	data, err := client.api.RPC.State.GetStorageRaw(key, types.Hash(blockhash))
	if err != nil {
		return pack.U256{}, nil, pack.U64(uint64(0)), fmt.Errorf("get storage: %v", err)
	}

	// Decode the event data to get the burn log.
	burnEvent, err := decodeEventData(meta, data)
	if err != nil {
		return pack.U256{}, nil, pack.U64(uint64(0)), err
	}

	// Calculate block confirmations for the event.
	confs := header.Number - block.Block.Header.Number + 1

	return pack.NewU256FromInt(burnEvent.Amount.Int), address.RawAddress(burnEvent.Dest[:]), pack.NewU64(uint64(confs)), nil
}

func decodeEventData(meta *types.Metadata, data *types.StorageDataRaw) (eventBurnt, error) {
	events := RenVmBridgeEvents{}
	if err := types.EventRecordsRaw(*data).DecodeEventRecords(meta, &events); err != nil {
		return eventBurnt{}, fmt.Errorf("decode event data: %v", err)
	}

	if len(events.RenVmBridge_Burnt) != 1 {
		return eventBurnt{}, fmt.Errorf("expected burn events: %v, got: %v", 1, len(events.RenVmBridge_Burnt))
	}

	return events.RenVmBridge_Burnt[0], nil
}

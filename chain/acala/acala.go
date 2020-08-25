package acala

import (
	"context"
	"fmt"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/renproject/multichain"
	"github.com/renproject/pack"
	"go.uber.org/zap"
)

const DefaultClientRPCURL = "http://127.0.0.1:9944"

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

func printEvents(meta *types.Metadata, data *types.StorageDataRaw, nhBlock types.Hash) error {
	er := types.EventRecordsRaw(*data)
	e := EventsWithMint{}
	err := DecodeEvents(&er, meta, &e)
	if err != nil {
		return err
	}

	for _, e := range e.RenToken_AssetsMinted {
		fmt.Printf("[EVENT] RenToken::AssetsMinted:: (phase=%#v)\n", e.Phase)
		fmt.Printf("[EVENT] RenToken::AssetsMinted:: (phase=%#v)\n", e.Who)
		fmt.Printf("[EVENT] RenToken::AssetsMinted:: (phase=%#v)\n", e.Currency)
		fmt.Printf("[EVENT] RenToken::AssetsMinted:: (phase=%#v)\n", e.Amount)
		fmt.Printf("[EVENT] RenToken::AssetsMinted:: (phase=%#v)\n", e.Topics)
	}

	return nil
}

func (client *Client) BurnEvent(ctx context.Context, asset multichain.Asset, nonce pack.Bytes32, blockheight pack.U64) (amount pack.U256, to pack.String, confs int64, err error) {

	meta, err := client.api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%#v\n", meta)

	// Subscribe to system events via storage
	key, err := types.CreateStorageKey(meta, "System", "Events", nil, nil)
	if err != nil {
		panic(err)
	}

	blockhash, err := client.api.RPC.Chain.GetBlockHash(blockheight.Uint64())
	if err != nil {
		panic(err)
	}

	fmt.Printf("blockhash: %#v\n", blockhash.Hex())

	data, err := client.api.RPC.State.GetStorageRaw(key, blockhash)

	if err = printEvents(meta, data, blockhash); err != nil {
		panic(err)
	}

	fmt.Printf("\nLive events:\n")

	// fmt.Printf("data: %#v\n", data)

	// panic("unimplemented")

	sub, err := client.api.RPC.State.SubscribeStorageRaw([]types.StorageKey{key})
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	// outer for loop for subscription notifications
	for {
		set := <-sub.Chan()
		// inner loop for the changes within one of those notifications
		for _, chng := range set.Changes {
			if !types.Eq(chng.StorageKey, key) || !chng.HasStorageData {
				// skip, we are only interested in events with content
				continue
			}

			// printEvents(meta, &chng.StorageData)
		}
	}

	// return pack.U256{}, pack.String(""), 0, fmt.Errorf("unimplemented")
}

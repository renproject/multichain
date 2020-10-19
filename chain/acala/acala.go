package acala

import (
	"bytes"
	"context"
	"fmt"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"
	"go.uber.org/zap"
)

const (
	DefaultClientRPCURL = "http://127.0.0.1:9944"
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

func decodeEventData(meta *types.Metadata, data *types.StorageDataRaw) (eventBurnt, error) {
	eventRecords := types.EventRecordsRaw(*data)

	events := EventsWithBurn{}
	if err := ParseEvents(&eventRecords, meta, &events); err != nil {
		return eventBurnt{}, err
	}
	if len(events.RenToken_Burnt) != 1 {
		return eventBurnt{}, fmt.Errorf("expected burn events: %v, got: %v", 1, len(events.RenToken_Burnt))
	}

	return events.RenToken_Burnt[0], nil
}

func (client *Client) BurnEvent(
	ctx context.Context,
	asset multichain.Asset,
	nonce pack.Bytes32,
	blockheight pack.U64,
) (pack.U256, pack.String, pack.U64, error) {
	// get metadata
	meta, err := client.api.RPC.State.GetMetadataLatest()
	if err != nil {
		// FIXME: return err instead of panicking
		panic(err)
	}

	// subscribe to system events via storage
	key, err := types.CreateStorageKey(meta, "System", "Events", nil, nil)
	if err != nil {
		// FIXME: return err instead of panicking
		panic(err)
	}

	// get the block's hash
	blockhash, err := client.api.RPC.Chain.GetBlockHash(blockheight.Uint64())
	if err != nil {
		// FIXME: return err instead of panicking
		panic(err)
	}

	header, err := client.api.RPC.Chain.GetHeaderLatest()
	if err != nil {
		// FIXME: return err instead of panicking
		panic(err)
	}

	// retrieve raw bytes for the stored data
	data, err := client.api.RPC.State.GetStorageRaw(key, blockhash)
	if err != nil {
		// FIXME: return err instead of panicking
		panic(err)
	}

	// decode event data to a burn event
	burnEvent, err := decodeEventData(meta, data)
	if err != nil {
		// FIXME: return err instead of panicking
		panic(err)
	}

	// calculate block confirmations for the event
	confs := uint64(header.Number) - blockheight.Uint64() + 1

	// get and encode destination address
	dest := types.NewAddressFromAccountID(burnEvent.Dest[:])
	buf := new(bytes.Buffer)
	encoder := scale.NewEncoder(buf)
	if err := dest.Encode(*encoder); err != nil {
		// FIXME: return err instead of panicking
		panic(err)
	}
	addrEncodeDecoder := NewAddressEncodeDecoder()
	to, err := addrEncodeDecoder.EncodeAddress(address.RawAddress(pack.NewBytes(buf.Bytes())))
	if err != nil {
		// FIXME: return err instead of panicking
		panic(err)
	}

	return pack.NewU256FromU128(pack.NewU128FromInt(burnEvent.Amount.Int)), pack.String(to), pack.NewU64(confs), nil
}

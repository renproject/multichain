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

const DefaultClientRPCURL = "ws://127.0.0.1:9933"

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
	// key, err := types.CreateStorageKey(meta, "RenToken", "Events", nil, nil)
	// if err != nil {
	// 	return err
	// }

	fmt.Printf("Meta: %#v\n", meta)

	// var er EventRecordsRaw
	// err = est.getStorage(key, &er, nhBlock)
	// if err != nil {
	// 	return err
	// }

	fmt.Printf("data: %#v\n", data)

	er := types.EventRecordsRaw(*data)

	e := EventsWithMint{}
	err := DecodeEvents(&er, meta, &e)
	if err != nil {
		return err
	}

	// decoder := scale.NewDecoder(bytes.NewReader(er))

	// // determine number of events
	// n, err := decoder.DecodeUintCompact()
	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("found %v events", n)

	// // iterate over events
	// for i := uint64(0); i < n.Uint64(); i++ {
	// 	fmt.Printf("decoding event #%v\n", i)

	// 	// decode Phase
	// 	phase := types.Phase{}
	// 	err := decoder.Decode(&phase)
	// 	if err != nil {
	// 		return fmt.Errorf("unable to decode Phase for event #%v: %v", i, err)
	// 	}

	// 	// decode EventID
	// 	id := types.EventID{}
	// 	err = decoder.Decode(&id)
	// 	if err != nil {
	// 		return fmt.Errorf("unable to decode EventID for event #%v: %v", i, err)
	// 	}

	// 	fmt.Printf("event #%v has EventID %v\n", i, id)

	// }

	// events := types.EventRecords{}
	// err := EventRecordsRaw(*data).DecodeEventRecords(meta, &events)
	// if err != nil {
	// 	panic(err)
	// }

	// Show what we are busy with
	for _, e := range e.RenToken_AssetsMinted {
		fmt.Printf("[EVENT] RenToken::AssetsMinted:: (phase=%#v)\n", e.Phase)
		fmt.Printf("[EVENT] RenToken::AssetsMinted:: (phase=%#v)\n", e.Who)
		fmt.Printf("[EVENT] RenToken::AssetsMinted:: (phase=%#v)\n", e.Currency)
		fmt.Printf("[EVENT] RenToken::AssetsMinted:: (phase=%#v)\n", e.Amount)
		fmt.Printf("[EVENT] RenToken::AssetsMinted:: (phase=%#v)\n", e.Topics)
	}
	// for _, e := range e.Balances_Endowed {
	// 	fmt.Printf("[EVENT] Balances:Endowed:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%#x, %v\n", e.Who, e.Balance)
	// }
	// for _, e := range e.Balances_DustLost {
	// 	fmt.Printf("[EVENT] Balances:DustLost:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%#x, %v\n", e.Who, e.Balance)
	// }
	// for _, e := range e.Balances_Transfer {
	// 	fmt.Printf("[EVENT] Balances:Transfer:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%v, %v, %v\n", e.From, e.To, e.Value)
	// }
	// for _, e := range e.Balances_BalanceSet {
	// 	fmt.Printf("[EVENT] Balances:BalanceSet:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%v, %v, %v\n", e.Who, e.Free, e.Reserved)
	// }
	// for _, e := range e.Balances_Deposit {
	// 	fmt.Printf("[EVENT] Balances:Deposit:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%v, %v\n", e.Who, e.Balance)
	// }
	// for _, e := range e.Grandpa_NewAuthorities {
	// 	fmt.Printf("[EVENT] Grandpa:NewAuthorities:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%v\n", e.NewAuthorities)
	// }
	// for _, e := range e.Grandpa_Paused {
	// 	fmt.Printf("[EVENT] Grandpa:Paused:: (phase=%#v)\n", e.Phase)
	// }
	// for _, e := range e.Grandpa_Resumed {
	// 	fmt.Printf("[EVENT] Grandpa:Resumed:: (phase=%#v)\n", e.Phase)
	// }
	// for _, e := range e.ImOnline_HeartbeatReceived {
	// 	fmt.Printf("[EVENT] ImOnline:HeartbeatReceived:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%#x\n", e.AuthorityID)
	// }
	// for _, e := range e.ImOnline_AllGood {
	// 	fmt.Printf("[EVENT] ImOnline:AllGood:: (phase=%#v)\n", e.Phase)
	// }
	// for _, e := range e.ImOnline_SomeOffline {
	// 	fmt.Printf("[EVENT] ImOnline:SomeOffline:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%v\n", e.IdentificationTuples)
	// }
	// for _, e := range e.Indices_IndexAssigned {
	// 	fmt.Printf("[EVENT] Indices:IndexAssigned:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%#x%v\n", e.AccountID, e.AccountIndex)
	// }
	// for _, e := range e.Indices_IndexFreed {
	// 	fmt.Printf("[EVENT] Indices:IndexFreed:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%v\n", e.AccountIndex)
	// }
	// for _, e := range e.Offences_Offence {
	// 	fmt.Printf("[EVENT] Offences:Offence:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%v%v\n", e.Kind, e.OpaqueTimeSlot)
	// }
	// for _, e := range e.Session_NewSession {
	// 	fmt.Printf("[EVENT] Session:NewSession:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%v\n", e.SessionIndex)
	// }
	// for _, e := range e.Staking_Reward {
	// 	fmt.Printf("[EVENT] Staking:Reward:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%v\n", e.Balance)
	// }
	// for _, e := range e.Staking_Slash {
	// 	fmt.Printf("[EVENT] Staking:Slash:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%#x%v\n", e.AccountID, e.Balance)
	// }
	// for _, e := range e.Staking_OldSlashingReportDiscarded {
	// 	fmt.Printf("[EVENT] Staking:OldSlashingReportDiscarded:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%v\n", e.SessionIndex)
	// }
	// for _, e := range e.System_ExtrinsicSuccess {
	// 	fmt.Printf("[EVENT] System:ExtrinsicSuccess:: (phase=%#v)\n", e.Phase)
	// }
	// for _, e := range e.System_ExtrinsicFailed {
	// 	fmt.Printf("[EVENT] System:ErtrinsicFailed:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%v\n", e.DispatchError)
	// }
	// for _, e := range e.System_CodeUpdated {
	// 	fmt.Printf("[EVENT] System:CodeUpdated:: (phase=%#v)\n", e.Phase)
	// }
	// for _, e := range e.System_NewAccount {
	// 	fmt.Printf("[EVENT] System:NewAccount:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%#x\n", e.Who)
	// }
	// for _, e := range e.System_KilledAccount {
	// 	fmt.Printf("[EVENT] System:KilledAccount:: (phase=%#v)\n", e.Phase)
	// 	fmt.Printf("\t%#X\n", e.Who)
	// }

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

package cosmos

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	// DefaultClientTimeout used by the Client.
	DefaultClientTimeout = time.Minute
	// DefaultClientTimeoutRetry used by the Client.
	DefaultClientTimeoutRetry = time.Second
	// DefaultClientHost used by the Client. This should only be used for local
	// deployments of the multichain.
	DefaultClientHost = "http://0.0.0.0:26657"
	// DefaultBroadcastMode configures the behaviour of a cosmos client while it
	// interacts with the cosmos node. Allowed broadcast modes can be async, sync
	// and block. "async" returns immediately after broadcasting, "sync" returns
	// after the transaction has been checked and "block" waits until the
	// transaction is committed to the chain.
	DefaultBroadcastMode = "sync"
)

// ClientOptions are used to parameterise the behaviour of the Client.
type ClientOptions struct {
	Timeout       time.Duration
	TimeoutRetry  time.Duration
	Host          pack.String
	BroadcastMode pack.String
}

// DefaultClientOptions returns ClientOptions with the default settings. These
// settings are valid for use with the default local deployment of the
// multichain. In production, the host, user, and password should be changed.
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		Timeout:       DefaultClientTimeout,
		TimeoutRetry:  DefaultClientTimeoutRetry,
		Host:          pack.String(DefaultClientHost),
		BroadcastMode: pack.String(DefaultBroadcastMode),
	}
}

// WithHost sets the URL of the Bitcoin node.
func (opts ClientOptions) WithHost(host pack.String) ClientOptions {
	opts.Host = host
	return opts
}

// Client interacts with an instance of the Cosmos based network using the REST
// interface exposed by a lightclient node.
type Client struct {
	opts   ClientOptions
	cliCtx cliContext.CLIContext
	cdc    *codec.Codec
}

// NewClient returns a new Client.
func NewClient(opts ClientOptions, cdc *codec.Codec) *Client {
	httpClient, err := rpchttp.NewWithTimeout(opts.Host.String(), "websocket", uint(opts.Timeout/time.Second))
	if err != nil {
		panic(err)
	}

	cliCtx := cliContext.NewCLIContext().WithCodec(cdc).WithClient(httpClient).WithTrustNode(true)

	return &Client{
		opts:   opts,
		cliCtx: cliCtx,
		cdc:    cdc,
	}
}

// Tx query transaction with txHash
func (client *Client) Tx(ctx context.Context, txHash pack.Bytes) (account.Tx, pack.U64, error) {
	res, err := utils.QueryTx(client.cliCtx, hex.EncodeToString(txHash[:]))
	if err != nil {
		return &StdTx{}, pack.NewU64(0), fmt.Errorf("query fail: %v", err)
	}

	authStdTx := res.Tx.(auth.StdTx)
	if res.Code != 0 {
		return &StdTx{}, pack.NewU64(0), fmt.Errorf("tx failed code: %v, log: %v", res.Code, res.RawLog)
	}

	stdTx, err := parseStdTx(authStdTx)
	if err != nil {
		return &StdTx{}, pack.NewU64(0), fmt.Errorf("parse tx failed: %v", err)
	}

	return &stdTx, pack.NewU64(1), nil
}

// SubmitTx to the Cosmos based network.
func (client *Client) SubmitTx(ctx context.Context, tx account.Tx) error {
	txBytes, err := tx.Serialize()
	if err != nil {
		return fmt.Errorf("bad \"submittx\": %v", err)
	}

	res, err := client.cliCtx.WithBroadcastMode(client.opts.BroadcastMode.String()).BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	if res.Code != 0 {
		return fmt.Errorf("tx failed code: %v, log: %v", res.Code, res.RawLog)
	}

	return nil
}

// AccountNonce returns the current nonce of the account. This is the nonce to
// be used while building a new transaction.
func (client *Client) AccountNonce(_ context.Context, addr address.Address) (pack.U256, error) {
	cosmosAddr, err := types.AccAddressFromBech32(string(addr))
	if err != nil {
		return pack.U256{}, fmt.Errorf("bad address: '%v': %v", addr, err)
	}

	accGetter := auth.NewAccountRetriever(client.cliCtx)
	acc, err := accGetter.GetAccount(Address(cosmosAddr).AccAddress())
	if err != nil {
		return pack.U256{}, err
	}

	return pack.NewU256FromU64(pack.NewU64(acc.GetSequence())), nil
}

// AccountNumber returns the account number for a given address.
func (client *Client) AccountNumber(_ context.Context, addr address.Address) (pack.U64, error) {
	cosmosAddr, err := types.AccAddressFromBech32(string(addr))
	if err != nil {
		return 0, fmt.Errorf("bad address: '%v': %v", addr, err)
	}

	accGetter := auth.NewAccountRetriever(client.cliCtx)
	acc, err := accGetter.GetAccount(Address(cosmosAddr).AccAddress())
	if err != nil {
		return 0, err
	}

	return pack.U64(acc.GetAccountNumber()), nil
}

// AccountBalance returns the account balancee for a given address.
func (client *Client) AccountBalance(_ context.Context, addr address.Address, denom string) (pack.U256, error) {
	cosmosAddr, err := types.AccAddressFromBech32(string(addr))
	if err != nil {
		return pack.U256{}, fmt.Errorf("bad address: '%v': %v", addr, err)
	}

	accGetter := auth.NewAccountRetriever(client.cliCtx)
	acc, err := accGetter.GetAccount(Address(cosmosAddr).AccAddress())
	if err != nil {
		return pack.U256{}, err
	}

	balance := acc.GetCoins().AmountOf(denom).BigInt()

	// If the balance exceeds `MaxU256`, return an error.
	if pack.MaxU256.Int().Cmp(balance) == -1 {
		return pack.U256{}, fmt.Errorf("balance %v for %v exceeds MaxU256", balance.String(), addr)
	}

	return pack.NewU256FromInt(balance), nil
}

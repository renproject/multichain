package cosmos

import (
	"fmt"
	"time"

	"github.com/renproject/pack"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
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
)

// ClientOptions are used to parameterise the behaviour of the Client.
type ClientOptions struct {
	Timeout      time.Duration
	TimeoutRetry time.Duration
	Host         string
}

// DefaultClientOptions returns ClientOptions with the default settings. These
// settings are valid for use with the default local deployment of the
// multichain. In production, the host, user, and password should be changed.
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		Timeout:      DefaultClientTimeout,
		TimeoutRetry: DefaultClientTimeoutRetry,
		Host:         DefaultClientHost,
	}
}

// WithHost sets the URL of the Bitcoin node.
func (opts ClientOptions) WithHost(host string) ClientOptions {
	opts.Host = host
	return opts
}

// A Client interacts with an instance of the Cosmos based network using the REST
// interface exposed by a lightclient node.
type Client interface {
	// Account query account with address
	Account(address Address) (Account, error)
	// Tx query transaction with txHash
	Tx(txHash pack.String) (StdTx, error)
	// SubmitTx to the Cosmos based network.
	SubmitTx(tx Tx, broadcastMode pack.String) (pack.String, error)
}

type client struct {
	opts   ClientOptions
	cliCtx context.CLIContext
}

// NewClient returns a new Client.
func NewClient(opts ClientOptions, cdc *codec.Codec) Client {
	httpClient, err := rpchttp.NewWithTimeout(opts.Host, "websocket", uint(opts.Timeout/time.Second))
	if err != nil {
		panic(err)
	}

	cliCtx := context.NewCLIContext().WithCodec(cdc).WithClient(httpClient).WithTrustNode(true)

	return &client{
		opts:   opts,
		cliCtx: cliCtx,
	}
}

// Account contains necessary info for sdk.Account
type Account struct {
	Address        Address  `json:"address"`
	AccountNumber  pack.U64 `json:"account_number"`
	SequenceNumber pack.U64 `json:"sequence_number"`
	Coins          Coins    `json:"coins"`
}

// Account query account with address
func (client *client) Account(addr Address) (Account, error) {
	accGetter := auth.NewAccountRetriever(client.cliCtx)
	acc, err := accGetter.GetAccount(addr.AccAddress())
	if err != nil {
		return Account{}, err
	}

	return Account{
		Address:        addr,
		AccountNumber:  pack.U64(acc.GetAccountNumber()),
		SequenceNumber: pack.U64(acc.GetSequence()),
		Coins:          parseCoins(acc.GetCoins()),
	}, nil
}

// Tx query transaction with txHash
func (client *client) Tx(txHash pack.String) (StdTx, error) {
	res, err := utils.QueryTx(client.cliCtx, txHash.String())
	if err != nil {
		return StdTx{}, err
	}

	stdTx := res.Tx.(auth.StdTx)
	if res.Code != 0 {
		return StdTx{}, fmt.Errorf("Tx Failed Code: %v, Log: %v", res.Code, res.RawLog)
	}

	return parseStdTx(stdTx)
}

// SubmitTx to the Cosmos based network.
func (client *client) SubmitTx(tx Tx, broadcastMode pack.String) (pack.String, error) {
	txBytes, err := tx.Serialize()
	if err != nil {
		return pack.String(""), fmt.Errorf("bad \"submittx\": %v", err)
	}

	res, err := client.cliCtx.WithBroadcastMode(broadcastMode.String()).BroadcastTx(txBytes)
	if err != nil {
		return pack.String(""), err
	}

	if res.Code != 0 {
		return pack.String(""), fmt.Errorf("Tx Failed Code: %v, Log: %v", res.Code, res.RawLog)
	}

	return pack.NewString(res.TxHash), nil
}

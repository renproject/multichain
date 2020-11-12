package starname

import (
	"github.com/iov-one/iovns/app"
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/chain/cosmos"
)

type (
	// Client re-exports cosmos.Client
	Client = cosmos.Client

	// ClientOptions re-exports cosmos.ClientOptions
	ClientOptions = cosmos.ClientOptions

	// TxBuilderOptions re-exports cosmos.TxBuilderOptions
	TxBuilderOptions = cosmos.TxBuilderOptions
)

var (
	// DefaultClientOptions re-exports default cosmos-compatible client options
	DefaultClientOptions = cosmos.DefaultClientOptions

	// NewGasEstimator re-exports cosmos.NewGasEstimator
	NewGasEstimator = cosmos.NewGasEstimator
)

// NewClient returns returns a new Client with starname codec
func NewClient(opts ClientOptions) *Client {
	return cosmos.NewClient(opts, app.MakeCodec(), "star")
}

// NewTxBuilder returns an implementation of the transaction builder interface
// from the Cosmos Compat API, and exposes the functionality to build simple
// Starname transactions.
func NewTxBuilder(opts TxBuilderOptions, client *Client) account.TxBuilder {
	return cosmos.NewTxBuilder(opts, client)
}

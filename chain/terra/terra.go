package terra

import (
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/chain/cosmos"
	"github.com/terra-project/core/app"
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
	// DefaultClientOptions re-exports cosmos.DefaultClientOptions
	DefaultClientOptions = cosmos.DefaultClientOptions

	// DefaultTxBuilderOptions re-exports cosmos.DefaultTxBuilderOptions
	DefaultTxBuilderOptions = cosmos.DefaultTxBuilderOptions

	// NewGasEstimator re-exports cosmos.NewGasEstimator
	NewGasEstimator = cosmos.NewGasEstimator
)

// NewClient returns returns a new Client with Terra codec.
func NewClient(opts ClientOptions) *Client {
	return cosmos.NewClient(opts, app.MakeCodec(), "terra")
}

// NewTxBuilder returns an implementation of the transaction builder interface
// from the Cosmos Compat API, and exposes the functionality to build simple
// Terra transactions.
func NewTxBuilder(opts TxBuilderOptions, client *Client) account.TxBuilder {
	return cosmos.NewTxBuilder(opts, client)
}

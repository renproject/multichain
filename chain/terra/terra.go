package terra

import (
	"github.com/cosmos/cosmos-sdk/types"
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

// Set the Bech32 address prefix for the globally-defined config variable inside
// Cosmos SDK. This is required as there are a number of functions inside the
// SDK that make use of this global config directly, instead of allowing us to
// provide a custom config.
func init() {
	// TODO: This will prevent us from being able to support multiple
	// Cosmos-compatible chains in the Multichain. This is expected to be
	// resolved before v1.0 of the Cosmos SDK (issue being tracked here:
	// https://github.com/cosmos/cosmos-sdk/issues/7448).
	types.GetConfig().SetBech32PrefixForAccount("terra", "terrapub")
	types.GetConfig().Seal()
}

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

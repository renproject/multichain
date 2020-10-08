package starname

import (
	"github.com/iov-one/iovns/app"
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
)

// NewClient returns returns a new Client with Starname (iovns) codec
func NewClient(opts ClientOptions) cosmos.CompositeClient {
	return cosmos.NewClient(opts, app.MakeCodec())
}

// NewTxBuilder returns an implementation of the transaction builder interface
// from the Cosmos Compat API, and exposes the functionality to build simple
// Terra transactions.
func NewTxBuilder(opts TxBuilderOptions) cosmos.TxBuilder {
	return cosmos.NewTxBuilder(opts)
}

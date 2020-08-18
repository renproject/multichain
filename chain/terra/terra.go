package terra

import (
	"github.com/renproject/multichain/chain/cosmos"
	"github.com/terra-project/core/app"
)

type (
	Client           = cosmos.Client
	ClientOptions    = cosmos.ClientOptions
	Tx               = cosmos.Tx
	TxBuilder        = cosmos.TxBuilder
	TxBuilderOptions = cosmos.TxBuilderOptions
)

// NewClient returns returns a new Client with terra codec
func NewClient(opts ClientOptions) Client {
	return cosmos.NewClient(opts, app.MakeCodec())
}

// NewTxBuilder returns an implementation of the transaction builder interface
// from the Cosmos Compat API, and exposes the functionality to build simple
// Terra transactions.
func NewTxBuilder(opts TxBuilderOptions) TxBuilder {
	return cosmos.NewTxBuilder(opts).WithCodec(app.MakeCodec())
}

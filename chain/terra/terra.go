package terra

import (
	"github.com/renproject/multichain/chain/cosmos"
	"github.com/renproject/multichain/compat/cosmoscompat"

	"github.com/terra-project/core/app"
)

// NewClient returns returns a new Client with terra codec
func NewClient(opts cosmoscompat.ClientOptions) cosmoscompat.Client {
	return cosmoscompat.NewClient(opts, app.MakeCodec())
}

// NewTxBuilder returns an implementation of the transaction builder interface
// from the Cosmos Compat API, and exposes the functionality to build simple
// Terra transactions.
func NewTxBuilder(opts cosmoscompat.TxOptions) cosmoscompat.TxBuilder {
	return cosmos.NewTxBuilder(opts).WithCodec(app.MakeCodec())
}

// The Tx type is copied from Cosmos.
type Tx = cosmos.Tx

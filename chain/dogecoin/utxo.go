package dogecoin

import "github.com/renproject/multichain/chain/utxochain"

type (
	// Tx re-exports utxo.Tx.
	Tx = utxochain.Tx

	// TxBuilder re-exports utxo.TxBuilder.
	TxBuilder = utxochain.TxBuilder

	// Client re-exports utxo.Client.
	Client = utxochain.Client

	// ClientOptions re-exports utxochain.ClientOptions.
	ClientOptions = utxochain.ClientOptions
)

var (
	// NewTxBuilder re-exports utxochain.NewTxBuilder.
	NewTxBuilder = utxochain.NewTxBuilder

	// NewClient re-exports utxochain.NewClient.
	NewClient = utxochain.NewClient
)

// DefaultClientOptions returns ClientOptions with the default settings. These
// settings are valid for use with the default local deployment of the
// multichain. In production, the host, user, and password should be changed.
func DefaultClientOptions() ClientOptions {
	return utxochain.DefaultClientOptions().WithHost("http://0.0.0.0:18332")
}

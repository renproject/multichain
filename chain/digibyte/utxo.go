package digibyte

import "github.com/renproject/multichain/chain/utxochain"

type (
	// Tx represents a simple Bitcoin transaction that implements the Bitcoin Compat
	// API.
	Tx = utxochain.Tx

	// The TxBuilder is an implementation of a UTXO-compatible transaction builder
	// for Bitcoin.
	TxBuilder = utxochain.TxBuilder

	// A Client interacts with an instance of the Bitcoin network using the RPC
	// interface exposed by a Bitcoin node.
	Client = utxochain.Client

	// ClientOptions are used to parameterise the behaviour of the Client.
	ClientOptions = utxochain.ClientOptions
)

var (
	// NewTxBuilder re-exports bitoin.NewTxBuilder
	NewTxBuilder = utxochain.NewTxBuilder

	// NewClient re-exports utxochain.NewClient
	NewClient = utxochain.NewClient
)

// DefaultClientOptions returns ClientOptions with the default settings. These
// settings are valid for use with the default local deployment of the
// multichain. In production, the host, user, and password should be changed.
func DefaultClientOptions() ClientOptions {
	return utxochain.DefaultClientOptions().WithHost("http://0.0.0.0:20443")
}

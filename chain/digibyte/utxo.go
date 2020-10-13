package digibyte

import "github.com/renproject/multichain/chain/bitcoin"

type (
	// Tx represents a simple Bitcoin transaction that implements the Bitcoin Compat
	// API.
	Tx = bitcoin.Tx

	// The TxBuilder is an implementation of a UTXO-compatible transaction builder
	// for Bitcoin.
	TxBuilder = bitcoin.TxBuilder

	// A Client interacts with an instance of the Bitcoin network using the RPC
	// interface exposed by a Bitcoin node.
	Client = bitcoin.Client

	// ClientOptions are used to parameterise the behaviour of the Client.
	ClientOptions = bitcoin.ClientOptions
)

var (
	// NewTxBuilder re-exports bitoin.NewTxBuilder
	NewTxBuilder = bitcoin.NewTxBuilder

	// NewClient re-exports bitcoin.NewClient
	NewClient = bitcoin.NewClient
)

// DefaultClientOptions returns ClientOptions with the default settings. These
// settings are valid for use with the default local deployment of the
// multichain. In production, the host, user, and password should be changed.
func DefaultClientOptions() ClientOptions {
	return bitcoin.DefaultClientOptions().WithHost("http://0.0.0.0:20443")
}

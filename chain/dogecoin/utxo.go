package dogecoin

import "github.com/renproject/multichain/chain/bitcoin"

type (
	// Tx re-exports bitcoin.Tx.
	Tx = bitcoin.Tx

	// TxBuilder re-exports bitcoin.TxBuilder.
	TxBuilder = bitcoin.TxBuilder

	// Client re-exports bitcoin.Client.
	Client = bitcoin.Client

	// ClientOptions re-exports bitcoin.ClientOptions.
	ClientOptions = bitcoin.ClientOptions
)

var (
	// NewTxBuilder re-exports bitcoin.NewTxBuilder.
	NewTxBuilder = bitcoin.NewTxBuilder

	// NewClient re-exports bitcoin.NewClient.
	NewClient = bitcoin.NewClient
)

// DefaultClientOptions returns ClientOptions with the default settings. These
// settings are valid for use with the default local deployment of the
// multichain. In production, the host, user, and password should be changed.
func DefaultClientOptions() ClientOptions {
	return bitcoin.DefaultClientOptions().WithHost("http://0.0.0.0:18332")
}

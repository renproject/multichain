package dogecoin

import "github.com/renproject/multichain/chain/bitcoin"

type (
	Tx            = bitcoin.Tx
	TxBuilder     = bitcoin.TxBuilder
	Client        = bitcoin.Client
	ClientOptions = bitcoin.ClientOptions
)

var (
	NewTxBuilder = bitcoin.NewTxBuilder
	NewClient    = bitcoin.NewClient
)

// DefaultClientOptions returns ClientOptions with the default settings. These
// settings are valid for use with the default local deployment of the
// multichain. In production, the host, user, and password should be changed.
func DefaultClientOptions() ClientOptions {
	return bitcoin.DefaultClientOptions().WithHost("http://0.0.0.0:18332")
}

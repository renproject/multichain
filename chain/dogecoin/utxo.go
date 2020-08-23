package dogecoin

import "github.com/renproject/multichain/chain/bitcoin"

type (
	Tx            = bitcoin.Tx
	TxBuilder     = bitcoin.TxBuilder
	Client        = bitcoin.Client
	ClientOptions = bitcoin.ClientOptions
)

var (
	NewTxBuilder         = bitcoin.NewTxBuilder
	NewClient            = bitcoin.NewClient
	DefaultClientOptions = bitcoin.DefaultClientOptions
)

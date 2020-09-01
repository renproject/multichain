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

	// DefaultClientOptions re-exports bitcoin.DefaultClientOptions
	DefaultClientOptions = bitcoin.DefaultClientOptions
)

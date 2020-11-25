package qtum

import "github.com/renproject/multichain/chain/bitcoin"

type (
	// Tx re-exports bitcoin.Tx.
	Tx = bitcoin.Tx

	// TxBuilder re-exports bitcoin.TxBuilder.
	TxBuilder = bitcoin.TxBuilder
)

var (
	// NewTxBuilder re-exports bitcoin.NewTxBuilder.
	NewTxBuilder = bitcoin.NewTxBuilder
)

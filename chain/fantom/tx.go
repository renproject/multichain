package fantom

import (
	"github.com/renproject/multichain/chain/ethereum"
)

type (
	// TxBuilder re-exports ethereum.TxBuilder.
	TxBuilder = ethereum.TxBuilder

	// Tx re-exports ethereum.Tx.
	Tx = ethereum.Tx
)

// NewTxBuilder re-exports ethereum.NewTxBuilder.
var NewTxBuilder = ethereum.NewTxBuilder

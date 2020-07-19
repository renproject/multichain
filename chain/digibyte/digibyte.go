package digibyte

import (
	"github.com/renproject/multichain/chain/bitcoin"
	"github.com/renproject/multichain/compat/bitcoincompat"
)

// NewTxBuilder returns an implementation of the transaction builder interface
// from the Bitcoin Compat API, and exposes the functionality to build simple
// Dogecoin transactions.
func NewTxBuilder() bitcoincompat.TxBuilder {
	return bitcoin.NewTxBuilder()
}

// The Tx type is copied from Bitcoin.
type Tx = bitcoin.Tx

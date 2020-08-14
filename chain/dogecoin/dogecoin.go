package dogecoin

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/renproject/multichain/chain/bitcoin"
)

// NewTxBuilder returns an implementation of the transaction builder interface
// from the Bitcoin Compat API, and exposes the functionality to build simple
// Dogecoin transactions.
func NewTxBuilder(params *chaincfg.Params) bitcoin.TxBuilder {
	return bitcoin.NewTxBuilder(params)
}

// The Tx type is copied from Bitcoin.
type Tx = bitcoin.Tx

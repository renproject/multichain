// Package gas defines the Gas API. All chains that support transactions (either
// account-based or utxo-based) should implement this API. This API is used to
// understand the current recommended gas costs required to get confirmations in
// a reasonable amount of time.
package gas

import (
	"context"

	"github.com/renproject/pack"
)

type TxType uint8

const (
	ETHTransfer = TxType(0)
)

// The Estimator interface defines the functionality required to know the
// current recommended gas prices.
type Estimator interface {
	// EstimateGasPrice that should be used to get confirmation within a
	// reasonable amount of time. The precise definition of "reasonable amount
	// of time" varies from chain to chain, and so is left open to
	// interpretation by the implementation. For example, in Bitcoin, it would
	// be the recommended SATs-per-byte required to get a transaction into the
	// next block. In Ethereum, it would be the recommended GWEI-per-gas
	// required to get a transaction into one of the next few blocks (because
	// blocks happen a lot faster).
	EstimateGasPrice(context.Context) (pack.U256, error)

	// EstimateGasLimit ...
	EstimateGasLimit(TxType) (pack.U256, error)
}

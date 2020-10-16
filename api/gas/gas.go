// Package gas defines the Gas API. All chains that support transactions (either
// account-based or utxo-based) should implement this API. This API is used to
// understand the current recommended gas costs required to get confirmations in
// a reasonable amount of time.
package gas

import (
	"context"

	"github.com/renproject/pack"
)

// The Estimator interface defines the functionality required to know the
// current recommended gas price per gas unit and gas cap per gas unit. Not all
// chains have the concept of a gas cap, in which case it should be set to be
// equal to the gas price.
type Estimator interface {
	// EstimateGas base/price that should be used in order to get a transaction
	// confirmed within a reasonable amount of time. The precise definition of
	// "reasonable amount of time" varies from chain to chain, and so is left
	// open to interpretation by the implementation.

	// For example, in Bitcoin, the gas price (and gas cap) would be the
	// recommended SATs-per-byte required to get a transaction into the next
	// block.

	// In Ethereum without EIP-1559, the gas price (and gas cap) would be the
	// recommended GWEI-per-gas required to get a transaction into one of the
	// next few blocks (because blocks happen a lot faster). In Ethereum with
	// EIP1559, the gas price and gas cap would be back calculated given the gas
	// cap and an estimate of the current gas base.
	//
	// The assumption is that the total gas cost will be gasLimit * gasPrice <
	// gasLimit * gasCap. If chain nodes give back values based on different
	// assumptions, then the values must be normalised as needed.
	EstimateGas(context.Context) (gasPrice, gasCap pack.U256, err error)
}

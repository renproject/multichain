package cosmos

import (
	"context"

	"github.com/renproject/multichain/api/gas"
	"github.com/renproject/pack"
)

// A GasEstimator returns the gas-per-byte that is needed in order to confirm
// transactions with an estimated maximum delay of one block. In distributed
// networks that collectively build, sign, and submit transactions, it is
// important that all nodes in the network have reached consensus on the
// gas-per-byte.
type GasEstimator struct {
	gasPerByte pack.U256
}

// NewGasEstimator returns a simple gas estimator that always returns the same
// amount of gas-per-byte.
func NewGasEstimator(gasPerByte pack.U256) gas.Estimator {
	return &GasEstimator{
		gasPerByte: gasPerByte,
	}
}

// EstimateGas returns gas required per byte for Cosmos-compatible chains. This
// value is used for both the price and cap, because Cosmos-compatible chains do
// not have a distinct concept of cap.
func (gasEstimator *GasEstimator) EstimateGas(ctx context.Context) (pack.U256, pack.U256, error) {
	return gasEstimator.gasPerByte, gasEstimator.gasPerByte, nil
}

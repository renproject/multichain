package cosmos

import (
	"context"

	"github.com/renproject/multichain/api/gas"
	"github.com/renproject/pack"
)

const (
	microToPico = 1000000
)

// A GasEstimator returns the gas price that is needed in order to confirm
// transactions with an estimated maximum delay of one block. As of now, Cosmos
// compatible chains do not support transaction prioritisation in the mempool.
// Hence we use constant gas price set in the micro-denomination (1e-6) of the
// underlying token. However, the gas price as returned by the gas estimation
// API is a 32-byte unsigned integer that represents the gas price in the
// pico-denomination (1e-12).
type GasEstimator struct {
	gasPrice float64
}

// NewGasEstimator returns a simple gas estimator that always returns the same
// amount of gas-per-byte.
func NewGasEstimator(gasPrice float64) gas.Estimator {
	return &GasEstimator{
		gasPrice: gasPrice,
	}
}

// EstimateGas returns gas required per byte for Cosmos-compatible chains. This
// value is used for both the price and cap, because Cosmos-compatible chains do
// not have a distinct concept of cap. As of now, Cosmos compatible chains do
// not support transaction prioritisation in the mempool. Hence we use constant
// gas price set in the micro-denomination (1e-6) of the underlying token.
// However, the gas price as returned by the gas estimation API is a 32-byte
// unsigned integer representing gas price in the pico-denomination (1e-12).
func (gasEstimator *GasEstimator) EstimateGas(ctx context.Context) (pack.U256, pack.U256, error) {
	gasPrice := pack.NewU256FromUint64(uint64(gasEstimator.gasPrice * microToPico))
	return gasPrice, gasPrice, nil
}

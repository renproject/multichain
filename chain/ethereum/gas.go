package ethereum

import (
	"context"
	"fmt"

	"github.com/renproject/multichain/api/gas"
	"github.com/renproject/pack"
)

// A GasEstimator returns the gas price (in wei) that is needed in order to
// confirm transactions with an estimated maximum delay of one block. In
// distributed networks that collectively build, sign, and submit transactions,
// it is important that all nodes in the network have reached consensus on the
// gas price.
type GasEstimator struct {
	wei pack.U256
}

// NewGasEstimator returns a simple gas estimator that always returns the given
// gas price (in wei) to be used for broadcasting an Ethereum transaction.
func NewGasEstimator(wei pack.U256) GasEstimator {
	return GasEstimator{
		wei: wei,
	}
}

// EstimateGasPrice returns the number of wei that is needed in order to confirm
// transactions with an estimated maximum delay of one block. It is the
// responsibility of the caller to know the number of bytes in their
// transaction.
func (gasEstimator GasEstimator) EstimateGasPrice(_ context.Context) (pack.U256, error) {
	return gasEstimator.wei, nil
}

// EstimateGasLimit returns the gas limit depending on what type of transaction we wish to do
func (gasEstimator GasEstimator) EstimateGasLimit(txType gas.TxType) (pack.U256, error) {
	switch txType {
	case gas.ETHTransfer:
		return pack.NewU256FromU64(pack.NewU64(21000)), nil
	default:
		return pack.NewU256([32]byte{}), fmt.Errorf("non-exhaustive transaction type: %v", txType)
	}
}

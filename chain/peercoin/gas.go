package peercoin

import (
	"context"

	"github.com/renproject/pack"
)

// A GasEstimator returns the SATs-per-byte that is needed in order to confirm
// transactions with an estimated maximum delay of one block. In distributed
// networks that collectively build, sign, and submit transactions, it is
// important that all nodes in the network have reached consensus on the
// SATs-per-byte.
type GasEstimator struct {
	satsPerByte pack.U256
}

// NewGasEstimator returns a simple gas estimator that always returns the given
// number of SATs-per-byte.
func NewGasEstimator(satsPerByte pack.U256) GasEstimator {
	return GasEstimator{
		satsPerByte: satsPerByte,
	}
}

// GasPrice returns the number of SATs-per-byte that is needed in order to
// confirm transactions with an estimated maximum delay of one block. It is the
// responsibility of the caller to know the number of bytes in their
// transaction.
func (gasEstimator GasEstimator) GasPrice(_ context.Context) (pack.U256, error) {
	return gasEstimator.satsPerByte, nil
}

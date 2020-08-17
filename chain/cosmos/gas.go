package cosmos

import (
	"context"

	"github.com/renproject/pack"
)

// A GasEstimator returns the gas-per-byte that is needed in order to confirm
// transactions with relatively high speed. In distributed networks that work to
// collectively build transactions, it is important that all nodes return the
// same values from this interface.
type GasEstimator interface {
	GasPerByte(ctx context.Context) (pack.U64, error)
	GasPerSignature(ctx context.Context) (pack.U64, error)
}

type gasEstimator struct {
	gasPerByte      pack.U64
	gasPerSignature pack.U64
}

// NewGasEstimator returns a simple gas estimator that always returns the same
// amount of gas-per-byte.
func NewGasEstimator(gasPerByte, gasPerSignature pack.U64) GasEstimator {
	return &gasEstimator{
		gasPerByte:      gasPerByte,
		gasPerSignature: gasPerSignature,
	}
}

func (gasEstimator *gasEstimator) GasPerByte(ctx context.Context) (pack.U64, error) {
	return gasEstimator.gasPerByte, nil
}

func (gasEstimator *gasEstimator) GasPerSignature(ctx context.Context) (pack.U64, error) {
	return gasEstimator.gasPerSignature, nil
}

package bitcoincompat

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
}

type gasEstimator struct {
	satsPerByte pack.U64
}

// NewGasEstimator returns a simple gas estimator that always returns the same
// amount of gas-per-byte.
func NewGasEstimator(satsPerByte pack.U64) GasEstimator {
	return &gasEstimator{
		satsPerByte: satsPerByte,
	}
}

func (gasEstimator *gasEstimator) GasPerByte(ctx context.Context) (pack.U64, error) {
	return gasEstimator.satsPerByte, nil
}

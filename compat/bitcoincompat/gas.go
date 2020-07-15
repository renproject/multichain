package bitcoincompat

import (
	"context"

	"github.com/renproject/pack"
)

type GasEstimator interface {
	GasPerByte(ctx context.Context) (pack.U64, error)
}

type gasEstimator struct {
	satsPerByte pack.U64
}

func NewGasEstimator(satsPerByte pack.U64) GasEstimator {
	return &gasEstimator{
		satsPerByte: satsPerByte,
	}
}

func (gasEstimator *gasEstimator) GasPerByte(ctx context.Context) (pack.U64, error) {
	return gasEstimator.satsPerByte, nil
}

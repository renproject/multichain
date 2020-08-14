package gas

import (
	"context"

	"github.com/renproject/pack"
)

type Estimator interface {
	EstimateGasPrice(context.Context) (pack.U256, error)
}

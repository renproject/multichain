package bitcoincash

import (
	"context"
	"fmt"
	"math"

	"github.com/renproject/pack"
)

const (
	bchToSatoshis  = 1e8
	kilobyteToByte = 1024
)

// A GasEstimator returns the SATs-per-byte that is needed in order to confirm
// transactions with an estimated maximum delay of one block. In distributed
// networks that collectively build, sign, and submit transactions, it is
// important that all nodes in the network have reached consensus on the
// SATs-per-byte.
type GasEstimator struct {
	client      Client
	fallbackGas pack.U256
}

// NewGasEstimator returns a simple gas estimator that always returns the given
// number of SATs-per-byte.
func NewGasEstimator(client Client, fallbackGas pack.U256) GasEstimator {
	return GasEstimator{
		client:      client,
		fallbackGas: fallbackGas,
	}
}

// EstimateGas returns the number of SATs-per-byte (for both price and cap) that
// is needed in order to confirm transactions with a minimal delay. It is the
// responsibility of the caller to know the number of bytes in their
// transaction. This method calls the `estimatefee` RPC call to the node, which
// based on a conservative (considering longer history) strategy returns the
// estimated BCH per kilobyte of data in the transaction. An error will be
// returned if the node hasn't observed enough blocks to make an estimate.
func (gasEstimator GasEstimator) EstimateGas(ctx context.Context) (pack.U256, pack.U256, error) {
	feeRate, err := gasEstimator.client.EstimateFeeLegacy(ctx, int64(0))
	if err != nil {
		return gasEstimator.fallbackGas, gasEstimator.fallbackGas, err
	}

	if feeRate <= 0.0 {
		return gasEstimator.fallbackGas, gasEstimator.fallbackGas, fmt.Errorf("invalid fee rate: %v", feeRate)
	}

	satsPerByte := uint64(math.Ceil(feeRate * bchToSatoshis / kilobyteToByte))
	return pack.NewU256FromUint64(satsPerByte), pack.NewU256FromUint64(satsPerByte), nil
}

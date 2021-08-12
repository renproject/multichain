package evm

import (
	"context"
	"fmt"

	"github.com/renproject/pack"
)

// A GasEstimator returns the gas price and the provide gas limit that is needed in
// order to confirm transactions with an estimated maximum delay of one block.
type GasEstimator struct {
	client *Client
}

// NewGasEstimator returns a simple gas estimator that fetches the ideal gas
// price for a ethereum transaction to be included in a block
// with minimal delay.
func NewGasEstimator(client *Client) *GasEstimator {
	return &GasEstimator{
		client: client,
	}
}

// EstimateGas returns an estimate of the current gas price
// and returns the gas limit provided. These numbers change with congestion. These estimates
// are often a little bit off, and this should be considered when using them.
func (gasEstimator *GasEstimator) EstimateGas(ctx context.Context) (pack.U256, pack.U256, error) {
	gasPrice, err := gasEstimator.client.EthClient.SuggestGasPrice(ctx)
	if err != nil {
		return pack.NewU256([32]byte{}), pack.NewU256([32]byte{}), fmt.Errorf("failed to get eth suggested gas price: %v", err)
	}
	return pack.NewU256FromInt(gasPrice), pack.NewU256FromInt(gasPrice), nil
}

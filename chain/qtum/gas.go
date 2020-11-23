package qtum
// DEZU: TODO: This should probably be replaced with something from Qtum

import "github.com/renproject/multichain/chain/bitcoin"

// GasEstimator re-exports bitcoin.GasEstimator.
type GasEstimator = bitcoin.GasEstimator

// NewGasEstimator re-exports bitcoin.NewGasEstimator.
var NewGasEstimator = bitcoin.NewGasEstimator

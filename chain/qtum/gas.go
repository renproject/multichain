package qtum

// I *Think* it is fine to simply use Bitcoin's estimator here

import "github.com/renproject/multichain/chain/bitcoin"

// GasEstimator re-exports bitcoin.GasEstimator.
type GasEstimator = bitcoin.GasEstimator

// NewGasEstimator re-exports bitcoin.NewGasEstimator.
var NewGasEstimator = bitcoin.NewGasEstimator

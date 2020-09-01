package bitcoincash

import "github.com/renproject/multichain/chain/bitcoin"

// A GasEstimator returns the SATs-per-byte that is needed in order to confirm
// transactions with an estimated maximum delay of one block. In distributed
// networks that collectively build, sign, and submit transactions, it is
// important that all nodes in the network have reached consensus on the
// SATs-per-byte.
type GasEstimator = bitcoin.GasEstimator

// NewGasEstimator returns a simple gas estimator that always returns the given
// number of SATs-per-byte.
var NewGasEstimator = bitcoin.NewGasEstimator

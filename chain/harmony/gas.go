package harmony

import(
	"context"
	"math/big"

	"github.com/renproject/pack"
)

var (
	defaultGas = big.NewInt(1)
)

type Estimator struct {}

// The average block time on Harmony is 5 seconds & each block has a max
// gas limit of 80 million. There is currently no need to estimate gas for
// regular transactions & we do not have the RPC for it.
func (Estimator) EstimateGasPrice(ctx context.Context) (pack.U256, error) {
	return pack.NewU256FromInt(defaultGas), nil
}

package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"sort"

	"github.com/renproject/pack"
)

const (
	// FeeHistoryBlocks specifies how many blocks to consider for priority fee estimation
	FeeHistoryBlocks = 10
	// FeeHistoryPercentile specifies the percentile of effective priority fees to include
	FeeHistoryPercentile = 5
	// FallbackMaxFeePerGas is the fallback value used when MaxFeePerGas cannot be calculated
	FallbackMaxFeePerGas = 20000000000
)

var (
	// PriorityFeeEstimationTrigger specifies which base fee to trigger priority fee estimation at
	PriorityFeeEstimationTrigger = big.NewInt(100000000000) // WEI
	// DefaultPriorityFee is returned if above trigger is not met
	DefaultPriorityFee = big.NewInt(3000000000)
	// PriorityFeeIncreaseBoundary signifies a big bump in fee history priority reward, due to which we choose
	// not to consider values under it while calculating the median priority fee.
	PriorityFeeIncreaseBoundary = big.NewInt(200)
)

type feeHistoryResult struct {
	Reward      [][]string `json:"reward"`
}

// GasOptions allow a user to configure the parameters used while heuristically recommending
// fees for EIP-1559 compatible transactions.
type GasOptions struct {
	FeeHistoryBlocks             uint64
	FeeHistoryPercentile         uint64
	FallbackMaxFeePerGas         uint64
	PriorityFeeEstimationTrigger *big.Int
	DefaultPriorityFee           *big.Int
	PriorityFeeIncreaseBoundary  *big.Int
}

// A GasEstimator returns the gas price and the provide gas limit that is needed in
// order to confirm transactions with an estimated maximum delay of one block.
type GasEstimator struct {
	client  *Client
	options *GasOptions
}

// NewGasEstimator returns a simple gas estimator that fetches the ideal gas
// price for an ethereum transaction to be included in a block
// with minimal delay.
func NewGasEstimator(client *Client, opts GasOptions) *GasEstimator {
	return &GasEstimator{
		client:  client,
		options: &opts,
	}
}

// NewDefaultGasEstimator returns a simple gas estimator with default gas options
// that fetches the ideal gas price for an ethereum transaction to be included
// in a block with minimal delay.
func NewDefaultGasEstimator(client *Client) *GasEstimator {
	return &GasEstimator{
		client: client,
		options: &GasOptions{
			FeeHistoryBlocks,
			FeeHistoryPercentile,
			FallbackMaxFeePerGas,
			PriorityFeeEstimationTrigger,
			DefaultPriorityFee,
			PriorityFeeIncreaseBoundary,
		},
	}
}

// EstimateGas returns an estimate of the current gas price
// and returns the gas limit provided. These numbers change with congestion. These estimates
// are often a little off, and this should be considered when using them.
func (gasEstimator *GasEstimator) EstimateGas(ctx context.Context) (pack.U256, pack.U256, error) {
	latest, err := gasEstimator.client.EthClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return pack.NewU256([32]byte{}), pack.NewU256([32]byte{}), fmt.Errorf("failed to get eth suggested gas price: %v", err)
	}
	// base fee is in wei
	if latest.BaseFee == nil {
		// fallback values
		return pack.NewU256FromInt(gasEstimator.options.DefaultPriorityFee), pack.NewU256FromUint64(gasEstimator.options.FallbackMaxFeePerGas), nil
	}

	baseFee := new(big.Int).Set(latest.BaseFee)
	estimatedPriorityFee, err := gasEstimator.estimatePriorityFee(ctx, baseFee, latest.Number)
	if err != nil {
		return pack.NewU256([32]byte{}), pack.NewU256([32]byte{}), err
	}
	if estimatedPriorityFee == nil {
		// fallback values
		return pack.NewU256FromInt(gasEstimator.options.DefaultPriorityFee), pack.NewU256FromUint64(gasEstimator.options.FallbackMaxFeePerGas), nil
	}

	maxPriorityFeePerGas := gasEstimator.options.DefaultPriorityFee
	if estimatedPriorityFee.Cmp(maxPriorityFeePerGas) == 1 {
		maxPriorityFeePerGas = estimatedPriorityFee
	}

	potentialMaxFee := new(big.Int).Mul(baseFee, big.NewInt(12))
	if baseFee.Cmp(big.NewInt(40000000000)) == -1 {
		potentialMaxFee = new(big.Int).Mul(baseFee, big.NewInt(20))
	} else if baseFee.Cmp(big.NewInt(100000000000)) == -1 {
		potentialMaxFee = new(big.Int).Mul(baseFee, big.NewInt(16))
	} else if baseFee.Cmp(big.NewInt(200000000000)) == -1 {
		potentialMaxFee = new(big.Int).Mul(baseFee, big.NewInt(14))
	}
	potentialMaxFee.Div(potentialMaxFee, big.NewInt(10))

	maxFeePerGas := potentialMaxFee
	if maxPriorityFeePerGas.Cmp(potentialMaxFee) == 1 {
		maxFeePerGas = potentialMaxFee.Add(potentialMaxFee, maxPriorityFeePerGas)
	}
	return pack.NewU256FromInt(maxPriorityFeePerGas), pack.NewU256FromInt(maxFeePerGas), nil
}

func (gasEstimator *GasEstimator) estimatePriorityFee(ctx context.Context, baseFee *big.Int, blockNumber *big.Int) (*big.Int, error) {
	if baseFee.Cmp(gasEstimator.options.PriorityFeeEstimationTrigger) == -1 {
		return gasEstimator.options.DefaultPriorityFee, nil
	}
	var feeHistory feeHistoryResult

	if err := gasEstimator.client.RpcClient.CallContext(ctx, &feeHistory, "eth_feeHistory", gasEstimator.options.FeeHistoryBlocks, "0x"+blockNumber.Text(16), []int{int(gasEstimator.options.FeeHistoryPercentile)}); err != nil {
		return nil, fmt.Errorf("failed to get eth fee history: %v", err)
	}
	rewards := []*big.Int{}

	// filter and remove outliers
	for _, r := range feeHistory.Reward {
		if res, success := new(big.Int).SetString(r[0], 0); success && res.Cmp(big.NewInt(0)) != 0 {
			rewards = append(rewards, res)
		}
	}
	// sort in ascending order
	sort.Slice(rewards, func(i, j int) bool { return rewards[j].Cmp(rewards[i]) >= 0 })

	// if len <=1 percentage increase cannot be calculated
	if len(rewards) <= 1 {
		return nil, nil
	}

	percentageIncreases := []*big.Int{}
	for i, r := range rewards {
		if i == (len(rewards) - 1) {
			continue
		}
		next := new(big.Int).Set(rewards[i+1])
		temp := next.Sub(next, r)
		temp = temp.Div(temp, r)
		temp = temp.Mul(temp, big.NewInt(100))
		percentageIncreases = append(percentageIncreases, temp)
	}
	highestIncrease := percentageIncreases[0]
	highestIncreaseIndex := 0
	for i := 1; i < len(percentageIncreases); i++ {
		if highestIncrease.Cmp(percentageIncreases[i]) == -1 {
			highestIncrease = percentageIncreases[i]
			highestIncreaseIndex = i
		}
	}
	if highestIncrease.Cmp(gasEstimator.options.PriorityFeeIncreaseBoundary) == 1 && highestIncreaseIndex >= (len(rewards)/2) {
		rewards = rewards[highestIncreaseIndex:]
	}
	return rewards[len(rewards)/2], nil
}

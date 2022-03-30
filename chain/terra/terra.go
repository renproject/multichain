package terra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/chain/cosmos"
	"github.com/renproject/pack"
	"github.com/terra-money/core/app"
)

type (
	// Client re-exports cosmos.Client
	Client = cosmos.Client

	// ClientOptions re-exports cosmos.ClientOptions
	ClientOptions = cosmos.ClientOptions

	// TxBuilderOptions re-exports cosmos.TxBuilderOptions
	TxBuilderOptions = cosmos.TxBuilderOptions
)

var (
	// DefaultClientOptions re-exports cosmos.DefaultClientOptions
	DefaultClientOptions = cosmos.DefaultClientOptions

	// DefaultTxBuilderOptions re-exports cosmos.DefaultTxBuilderOptions
	DefaultTxBuilderOptions = cosmos.DefaultTxBuilderOptions

	// NewGasEstimator re-exports cosmos.NewGasEstimator
	NewGasEstimator = cosmos.NewGasEstimator
)

// Set the Bech32 address prefix for the globally-defined config variable inside
// Cosmos SDK. This is required as there are a number of functions inside the
// SDK that make use of this global config directly, instead of allowing us to
// provide a custom config.
func init() {
	// TODO: This will prevent us from being able to support multiple
	// Cosmos-compatible chains in the Multichain. This is expected to be
	// resolved before v1.0 of the Cosmos SDK (issue being tracked here:
	// https://github.com/cosmos/cosmos-sdk/issues/7448).
	types.GetConfig().SetBech32PrefixForAccount("terra", "terrapub")
	types.GetConfig().Seal()
}

// NewClient returns returns a new Client with Terra codec.
func NewClient(opts ClientOptions) *Client {
	cfg := app.MakeEncodingConfig()
	return cosmos.NewClient(opts, cfg.Marshaler, cfg.TxConfig, cfg.InterfaceRegistry, cfg.Amino, "terra")
}

// NewTxBuilder returns an implementation of the transaction builder interface
// from the Cosmos Compat API, and exposes the functionality to build simple
// Terra transactions.
func NewTxBuilder(opts TxBuilderOptions, client *Client) account.TxBuilder {
	return cosmos.NewTxBuilder(opts, client)
}

type GasEstimator struct {
	url      string
	decimals int
}

func NewHttpGasEstimator(url string, decimals int) GasEstimator {
	return GasEstimator{
		url:      url,
		decimals: decimals,
	}
}

func (gasEstimator GasEstimator) EstimateGas(ctx context.Context) (pack.U256, pack.U256, error) {
	response, err := http.Get(gasEstimator.url)
	if err != nil {
		return pack.U256{}, pack.U256{}, err
	}
	defer response.Body.Close()

	var results map[string]string
	if err := json.NewDecoder(response.Body).Decode(&results); err != nil {
		return pack.U256{}, pack.U256{}, err
	}
	gasPriceStr, ok := results["uluna"]
	if !ok {
		return pack.U256{}, pack.U256{}, fmt.Errorf("no uluna in response")
	}
	gasPriceFloat, err := strconv.ParseFloat(gasPriceStr, 64)
	if err != nil {
		panic(err)
	}
	gasPrice := uint64(gasPriceFloat * float64(gasEstimator.decimals))
	return pack.NewU256FromUint64(gasPrice), pack.NewU256FromUint64(gasPrice), nil
}

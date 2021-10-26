package terra

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/chain/cosmos"
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
	// This config construction is based on app.MakeEncodingConfig in
	// https://github.com/terra-money/core. We do not import this method
	// directly as the repo is dependent on CosmWasm which fails to compile on
	// M1 architecture.
	amino := codec.NewLegacyAmino()
	interfaceRegistry := codecTypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txConfig := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

	std.RegisterLegacyAminoCodec(amino)
	std.RegisterInterfaces(interfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(amino)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)

	return cosmos.NewClient(opts, marshaler, txConfig, interfaceRegistry, amino, "terra")
}

// NewTxBuilder returns an implementation of the transaction builder interface
// from the Cosmos Compat API, and exposes the functionality to build simple
// Terra transactions.
func NewTxBuilder(opts TxBuilderOptions, client *Client) account.TxBuilder {
	return cosmos.NewTxBuilder(opts, client)
}

var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
)

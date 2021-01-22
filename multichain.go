// Package multichain defines all supported assets and chains. It also
// re-exports the individual multichain APIs.
package multichain

import (
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/multichain/api/gas"
	"github.com/renproject/multichain/api/utxo"
	"github.com/renproject/surge"
)

type (
	// An Address is a human-readable representation of a public identity. It can
	// be the address of an external account, contract, or script.
	Address = address.Address

	// The AddressEncodeDecoder interfaces combines encoding and decoding
	// functionality into one interface.
	AddressEncodeDecoder = address.EncodeDecoder

	// RawAddress is an address that has been decoded into its binary form.
	RawAddress = address.RawAddress
)

type (
	// The AccountTx interface defines the functionality that must be exposed by
	// account-based transactions.
	AccountTx = account.Tx

	// The AccountTxBuilder interface defines the functionality required to build
	// account-based transactions. Most chain implementations require additional
	// information, and this should be accepted during the construction of the
	// chain-specific transaction builder.
	AccountTxBuilder = account.TxBuilder

	// The AccountClient interface defines the functionality required to interact
	// with a chain over RPC.
	AccountClient = account.Client
)

type (
	// A UTXOutpoint identifies a specific output produced by a transaction.
	UTXOutpoint = utxo.Outpoint

	// A UTXOutput is produced by a transaction. It includes the conditions
	// required to spend the output (called the pubkey script, based on Bitcoin).
	UTXOutput = utxo.Output

	// A UTXOInput specifies an existing output, produced by a previous
	// transaction, to be consumed by another transaction. It includes the script
	// that meets the conditions specified by the consumed output (called the sig
	// script, based on Bitcoin).
	UTXOInput = utxo.Input

	// A UTXORecipient specifies an address, and an amount, for which a
	// transaction will produce an output. Depending on the output, the address
	// can take on different formats (e.g. in Bitcoin, addresses can be P2PK,
	// P2PKH, or P2SH).
	UTXORecipient = utxo.Recipient

	// A UTXOTx interfaces defines the functionality that must be exposed by
	// utxo-based transactions.
	UTXOTx = utxo.Tx

	// A UTXOTxBuilder interface defines the functionality required to build
	// account-based transactions. Most chain implementations require additional
	// information, and this should be accepted during the construction of the
	// chain-specific transaction builder.
	UTXOTxBuilder = utxo.TxBuilder

	// A UTXOClient interface defines the functionality required to interact with
	// a chain over RPC.
	UTXOClient = utxo.Client
)

type (
	// ContractCallData is used to specify a function and its parameters when
	// invoking business logic on a contract.
	ContractCallData = contract.CallData

	// The ContractCaller interface defines the functionality required to call
	// readonly functions on a contract. Calling functions that mutate contract
	// state should be done using the Account API.
	ContractCaller = contract.Caller
)

type (
	// The GasEstimator interface defines the functionality required to know the
	// current recommended gas prices.
	GasEstimator = gas.Estimator
)

// An Asset uniquely identifies assets using human-readable strings.
type Asset string

// Enumeration of supported assets. When introducing a new chain, or new asset
// from an existing chain, you must add a human-readable string to this set of
// enumerated values. Assets must be listed in alphabetical order.
const (
	BCH  = Asset("BCH")  // Bitcoin Cash
	BNB  = Asset("BNB")  // Binance Coin
	BTC  = Asset("BTC")  // Bitcoin
	CELO = Asset("CELO") // Celo
	DGB  = Asset("DGB")  // DigiByte
	DOGE = Asset("DOGE") // Dogecoin
	ETH  = Asset("ETH")  // Ether
	FIL  = Asset("FIL")  // Filecoin
	FTM  = Asset("FTM")  // Fantom
	LBC  = Asset("LBC")  // LBRY
	SOL  = Asset("SOL")  // Solana
	LUNA = Asset("LUNA") // Luna
	ZEC  = Asset("ZEC")  // Zcash

	// These assets are defined separately because they are mock assets. These
	// assets should only be used for testing.

	AMOCK1 = Asset("AMOCK1") // Account-based mock asset
	AMOCK2 = Asset("AMOCK2") // Account-based mock asset
	UMOCK  = Asset("UMOCK")  // UTXO-based mock asset
)

// OriginChain returns the chain upon which the asset originates. For example,
// the origin chain of BTC is Bitcoin.
func (asset Asset) OriginChain() Chain {
	switch asset {
	case BCH:
		return BitcoinCash
	case BNB:
		return BinanceSmartChain
	case BTC:
		return Bitcoin
	case CELO:
		return Celo
	case DGB:
		return DigiByte
	case DOGE:
		return Dogecoin
	case ETH:
		return Ethereum
	case FIL:
		return Filecoin
	case FTM:
		return Fantom
	case LBC:
		return LBRY
	case LUNA:
		return Terra
	case SOL:
		return Solana
	case ZEC:
		return Zcash

	// These assets are handled separately because they are mock assets. These
	// assets should only be used for testing.

	case AMOCK1:
		return AccountMocker1
	case AMOCK2:
		return AccountMocker2
	case UMOCK:
		return UTXOMocker

	default:
		return Chain("")
	}
}

// ChainType returns the chain-type (Account or UTXO) for the given asset
func (asset Asset) ChainType() ChainType {
	switch asset {
	case BCH, BTC, DGB, DOGE, LBC, ZEC:
		return ChainTypeUTXOBased
	case BNB, ETH, FIL, LUNA:
		return ChainTypeAccountBased

	// These assets are handled separately because they are mock assets. These
	// assets should only be used for testing.

	case AMOCK1, AMOCK2:
		return ChainTypeAccountBased
	case UMOCK:
		return ChainTypeUTXOBased

	default:
		return ChainType("")
	}
}

// SizeHint returns the number of bytes required to represent the asset in
// binary.
func (asset Asset) SizeHint() int {
	return surge.SizeHintString(string(asset))
}

// Marshal the asset to binary.
func (asset Asset) Marshal(buf []byte, rem int) ([]byte, int, error) {
	return surge.MarshalString(string(asset), buf, rem)
}

// Unmarshal the asset from binary.
func (asset *Asset) Unmarshal(buf []byte, rem int) ([]byte, int, error) {
	return surge.UnmarshalString((*string)(asset), buf, rem)
}

// A Chain uniquely identifies a blockchain using a human-readable string.
type Chain string

// Enumeration of supported chains. When introducing a new chain, you must add a
// human-readable string to this set of enumerated values. Chains must be listed
// in alphabetical order.
const (
	Acala             = Chain("Acala")
	BinanceSmartChain = Chain("BinanceSmartChain")
	Bitcoin           = Chain("Bitcoin")
	BitcoinCash       = Chain("BitcoinCash")
	Celo              = Chain("Celo")
	DigiByte          = Chain("DigiByte")
	Dogecoin          = Chain("Dogecoin")
	Ethereum          = Chain("Ethereum")
	Fantom            = Chain("Fantom")
	Filecoin          = Chain("Filecoin")
	LBRY              = Chain("LBRY")
	Solana            = Chain("Solana")
	Terra             = Chain("Terra")
	Zcash             = Chain("Zcash")

	// These chains are defined separately because they are mock chains. These
	// chains should only be used for testing.

	AccountMocker1 = Chain("AccountMocker1")
	AccountMocker2 = Chain("AccountMocker2")
	UTXOMocker     = Chain("UTXOMocker")
)

// SizeHint returns the number of bytes required to represent the chain in
// binary.
func (chain Chain) SizeHint() int {
	return surge.SizeHintString(string(chain))
}

// Marshal the chain to binary. You should not call this function directly,
// unless you are implementing marshalling for a container type.
func (chain Chain) Marshal(buf []byte, rem int) ([]byte, int, error) {
	return surge.MarshalString(string(chain), buf, rem)
}

// Unmarshal the chain from binary. You should not call this function directly,
// unless you are implementing unmarshalling for a container type.
func (chain *Chain) Unmarshal(buf []byte, rem int) ([]byte, int, error) {
	return surge.UnmarshalString((*string)(chain), buf, rem)
}

// ChainType returns the chain type (whether account-based or utxo-based chain)
// for the chain.
func (chain Chain) ChainType() ChainType {
	switch chain {
	case Bitcoin, BitcoinCash, DigiByte, Dogecoin, LBRY, Zcash:
		return ChainTypeUTXOBased
	case BinanceSmartChain, Ethereum, Filecoin, Terra:
		return ChainTypeAccountBased

	// These chains are handled separately because they are mock chains. These
	// chains should only be used for testing.

	case AccountMocker1, AccountMocker2:
		return ChainTypeAccountBased
	case UTXOMocker:
		return ChainTypeUTXOBased

	default:
		return ChainType("")
	}
}

// IsAccountBased returns true when invoked on an account-based chain, otherwise
// returns false.
func (chain Chain) IsAccountBased() bool {
	return chain.ChainType() == ChainTypeAccountBased
}

// IsUTXOBased returns true when invoked on a utxo-based chain, otherwise
// returns false.
func (chain Chain) IsUTXOBased() bool {
	return chain.ChainType() == ChainTypeUTXOBased
}

// NativeAsset returns the underlying native asset for a chain. For example, the
// root asset of Bitcoin chain is BTC.
func (chain Chain) NativeAsset() Asset {
	switch chain {
	case BinanceSmartChain:
		return BNB
	case BitcoinCash:
		return BCH
	case Bitcoin:
		return BTC
	case DigiByte:
		return DGB
	case Dogecoin:
		return DOGE
	case Ethereum:
		return ETH
	case Filecoin:
		return FIL
	case LBRY:
		return LBC
	case Terra:
		return LUNA
	case Zcash:
		return ZEC

	// These chains are handled separately because they are mock chains. These
	// chains should only be used for testing.

	case AccountMocker1:
		return AMOCK1
	case AccountMocker2:
		return AMOCK2
	case UTXOMocker:
		return UMOCK

	default:
		return Asset("")
	}
}

// ChainType represents the type of chain (whether account-based or utxo-based)
type ChainType string

const (
	// ChainTypeAccountBased is an identifier for all account-based chains,
	// namely, BinanceSmartChain, Ethereum, Filecoin, and so on.
	ChainTypeAccountBased = ChainType("Account")

	// ChainTypeUTXOBased is an identifier for all utxo-based chains, namely,
	// Bitcoin, BitcoinCash, DigiByte, and so on.
	ChainTypeUTXOBased = ChainType("UTXO")
)

// SizeHint returns the number of bytes required to represent the chain type in
// binary.
func (chainType ChainType) SizeHint() int {
	return surge.SizeHintString(string(chainType))
}

// Marshal the chain type to binary. You should not call this function directly,
// unless you are implementing marshalling for a container type.
func (chainType ChainType) Marshal(buf []byte, rem int) ([]byte, int, error) {
	return surge.MarshalString(string(chainType), buf, rem)
}

// Unmarshal the chain type from binary. You should not call this function
// directly, unless you are implementing unmarshalling for a container type.
func (chainType *ChainType) Unmarshal(buf []byte, rem int) ([]byte, int, error) {
	return surge.UnmarshalString((*string)(chainType), buf, rem)
}

// Network identifies the network type for the multichain deployment
type Network string

const (
	// NetworkLocalnet represents a local network for chains. It is usually only
	// accessible from the device running the network, and is not accessible
	// over the Internet.  Chain rules are often slightly different to allow for
	// faster block times and easier access to testing funds. This is also
	// sometimes referred to as "regnet" or "regression network". It should only
	// be used for local testing.
	NetworkLocalnet = Network("localnet")

	// NetworkDevnet represents the development network for chains. This network
	// is typically a deployed version of the localnet. Chain rules are often
	// slightly different to allow for faster block times and easier access to
	// testing funds.
	NetworkDevnet = Network("devnet")

	// NetworkTestnet represents the test network for chains. This network is
	// typically a publicly accessible network that has the same, or very
	// similar, chain rules compared to mainnet. Assets on this type of network
	// are usually not considered to have value.
	NetworkTestnet = Network("testnet")

	// NetworkMainnet represents the main network for chains.
	NetworkMainnet = Network("mainnet")
)

// SizeHint returns the number of bytes required to represent the network in
// binary.
func (net Network) SizeHint() int {
	return surge.SizeHintString(string(net))
}

// Marshal the network to binary. You should not call this function directly,
// unless you are implementing marshalling for a container type.
func (net Network) Marshal(buf []byte, rem int) ([]byte, int, error) {
	return surge.MarshalString(string(net), buf, rem)
}

// Unmarshal the network from binary. You should not call this function
// directly, unless you are implementing unmarshalling for a container type.
func (net *Network) Unmarshal(buf []byte, rem int) ([]byte, int, error) {
	return surge.UnmarshalString((*string)(net), buf, rem)
}

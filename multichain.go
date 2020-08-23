package multichain

import (
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/multichain/api/gas"
	"github.com/renproject/multichain/api/utxo"
	"github.com/renproject/multichain/chain/ethereum"
	"github.com/renproject/surge"
)

type (
	Address               = address.Address
	AddressEncodeDecoder  = address.EncodeDecoder
	EthereumCompatAddress = ethereum.Address
	RawAddress            = address.RawAddress
)

type (
	AccountTx        = account.Tx
	AccountTxBuilder = account.TxBuilder
	AccountClient    = account.Client
)

type (
	UTXOutpoint   = utxo.Outpoint
	UTXOutput     = utxo.Output
	UTXOInput     = utxo.Input
	UTXORecipient = utxo.Recipient
	UTXOTx        = utxo.Tx
	UTXOTxBuilder = utxo.TxBuilder
	UTXOClient    = utxo.Client
)

type (
	ContractCallData = contract.CallData
	ContractCaller   = contract.Caller
)

type (
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
	SOL  = Asset("SOL")  // Solana
	LUNA = Asset("LUNA") // Luna
	ZEC  = Asset("ZEC")  // Zcash
	PPC = Asset("PPC")   // Peercoin
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
	case LUNA:
		return Terra
	case PPC:
		return Peercoin
	case SOL:
		return Solana
	case ZEC:
		return Zcash
	default:
		return Chain("")
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
	Filecoin          = Chain("Filecoin")
	Peercoin    	  = Chain("Peercoin")
	Solana            = Chain("Solana")
	Terra             = Chain("Terra")
	Zcash             = Chain("Zcash")
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

package multichain

import "github.com/renproject/surge"

// An Asset uniquely identifies assets using human-readable strings.
type Asset string

// Enumeration of supported assets. When introducing a new chain, or new asset
// from an existing chain, you must add a human-readable string to this set of
// enumerated values. Assets must be listed in alphabetical order.
const (
	BCH  = Asset("BCH")  // Bitcoin Cash
	BTC  = Asset("BTC")  // Bitcoin
	DGB  = Asset("DGB")  // DigiByte
	DOGE = Asset("DOGE") // Dogecoin
	ETH  = Asset("ETH")  // Ether
	FIL  = Asset("FIL")  // Filecoin
	SOL  = Asset("SOL")  // Solana
	LUNA = Asset("LUNA") // Luna
	ONT  = Asset("ONT")  // ONT
	ZEC  = Asset("ZEC")  // Zcash
)

// OriginChain returns the chain upon which the asset originates. For example,
// the origin chain of BTC is Bitcoin.
func (asset Asset) OriginChain() Chain {
	switch asset {
	case BCH:
		return BitcoinCash
	case BTC:
		return Bitcoin
	case DGB:
		return DigiByte
	case DOGE:
		return Dogecoin
	case ETH:
		return Ethereum
	case SOL:
		return Solana
	case ONT:
		return Ontology
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
	Acala       = Chain("Acala")
	Bitcoin     = Chain("Bitcoin")
	BitcoinCash = Chain("BitcoinCash")
	DigiByte    = Chain("DigiByte")
	Dogecoin    = Chain("Dogecoin")
	Ethereum    = Chain("Ethereum")
	Filecoin    = Asset("Filecoin")
	Solana      = Chain("Solana")
	Terra       = Chain("Terra")
	Ontology    = Chain("Ontology")
	Zcash       = Chain("Zcash")
)

// SizeHint returns the number of bytes required to represent the chain in
// binary.
func (chain Chain) SizeHint() int {
	return surge.SizeHintString(string(chain))
}

// Marshal the chain to binary.
func (chain Chain) Marshal(buf []byte, rem int) ([]byte, int, error) {
	return surge.MarshalString(string(chain), buf, rem)
}

// Unmarshal the chain from binary.
func (chain *Chain) Unmarshal(buf []byte, rem int) ([]byte, int, error) {
	return surge.UnmarshalString((*string)(chain), buf, rem)
}

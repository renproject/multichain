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
	LUNA = Asset("LUNA") // LUNA
	ZEC  = Asset("ZEC")  // Zcash
)

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
	Ethereum    = Chain("Ethereum")
	Terra       = Chain("Terra")
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

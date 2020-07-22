package bitcoincompat

import (
	"github.com/btcsuite/btcutil"
	"github.com/renproject/pack"
)

// TxBuilder defines an interface that can be used to build simple Bitcoin
// transactions.
type TxBuilder interface {
	// BuildTx returns a simple Bitcoin transaction that consumes a set of
	// Bitcoin outputs and uses the funds to make payments to a set of Bitcoin
	// recipients. The sum value of the inputs must be greater than the sum
	// value of the outputs, and the difference is paid as a fee to the Bitcoin
	// network.
	BuildTx(inputs []Output, recipients []Recipient) (Tx, error)
}

// Tx defines an interface that must be implemented by all types of Bitcoin
// transactions.
type Tx interface {
	// Hash of the transaction.
	Hash() pack.Bytes32

	// Sighashes that need to be signed before this transaction can be
	// submitted.
	Sighashes() ([]pack.Bytes32, error)

	// Sign the transaction by injecting signatures and the serialized pubkey of
	// the signer.
	Sign([]pack.Bytes65, pack.Bytes) error

	// Serialize the transaction.
	Serialize() (pack.Bytes, error)
}

// An Address is a public address that can be encoded/decoded to/from strings.
// Addresses are usually formatted different between different network
// configurations.
type Address btcutil.Address

// An Outpoint is a transaction outpoint, identifying a specific part of a
// transaction output.
//
// https://developer.bitcoin.org/reference/transactions.html#outpoint-the-specific-part-of-a-specific-output
type Outpoint struct {
	Hash  pack.Bytes32 `json:"hash"`
	Index pack.U32     `json:"index"`
}

// An Input to a Bitcoin transaction. It consumes an outpoint, and contains a
// signature script that satisfies the pubkey script of the output that it is
// consuming.
//
// https://developer.bitcoin.org/reference/transactions.html#txin-a-transaction-input-non-coinbase
type Input struct {
	Outpoint  Outpoint   `json:"outpoint"`
	SigScript pack.Bytes `json:"sigScript"`
}

// An Output of a Bitcoin transaction. It contains a pubkey script that defines
// the conditions required to spend the output.
//
// https://developer.bitcoin.org/reference/transactions.html#txout-a-transaction-output
type Output struct {
	Outpoint     Outpoint   `json:"outpoint"`
	Value        pack.U64   `json:"value"`
	PubKeyScript pack.Bytes `json:"pubKeyScript"`
}

// A Recipient of funds from a Bitcoin transaction. This is useful for buidling
// simple pay-to-address Bitcoin transactions.
type Recipient struct {
	Address pack.String `json:"address"`
	Value   pack.U64    `json:"value"`
}

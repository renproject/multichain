// Package utxo defines the UTXO API. All chains that use a utxo-based model
// should implement this API. The UTXO API is used to send and confirm
// transactions between addresses.
package utxo

import (
	"context"

	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"
)

// An Outpoint identifies a specific output produced by a transaction.
type Outpoint struct {
	Hash  pack.Bytes `json:"hash"`
	Index pack.U32   `json:"index"`
}

// An Output is produced by a transaction. It includes the conditions required
// to spend the output (called the pubkey script, based on Bitcoin).
type Output struct {
	Outpoint     `json:"outpoint"`
	Value        pack.U256  `json:"value"`
	PubKeyScript pack.Bytes `json:"pubKeyScript"`
}

// An Input specifies an existing output, produced by a previous transaction, to
// be consumed by another transaction. It includes the script that meets the
// conditions specified by the consumed output (called the sig script, based on
// Bitcoin).
type Input struct {
	Output    `json:"output"`
	SigScript pack.Bytes `json:"sigScript"`
}

// A Recipient specifies an address, and an amount, for which a transaction will
// produce an output. Depending on the output, the address can take on different
// formats (e.g. in Bitcoin, addresses can be P2PK, P2PKH, or P2SH).
type Recipient struct {
	To    address.Address `json:"to"`
	Value pack.U256       `json:"value"`
}

// The Tx interfaces defines the functionality that must be exposed by
// utxo-based transactions.
type Tx interface {
	// Hash returns the hash that uniquely identifies the transaction.
	// Generally, hashes are irreversible hash functions that consume the
	// content of the transaction.
	Hash() (pack.Bytes, error)

	// Inputs consumed by the transaction.
	Inputs() ([]Input, error)

	// Outputs produced by the transaction.
	Outputs() ([]Output, error)

	// Sighashes that must be signed before the transaction can be submitted by
	// the client.
	Sighashes() ([]pack.Bytes32, error)

	// Sign the transaction by injecting signatures for the required sighashes.
	// The serialized public key used to sign the sighashes should also be
	// specified whenever it is available.
	Sign([]pack.Bytes65, pack.Bytes) error

	// Serialize the transaction into bytes. This is the format in which the
	// transaction will be submitted by the client.
	Serialize() (pack.Bytes, error)
}

// The TxBuilder interface defines the functionality required to build
// account-based transactions. Most chain implementations require additional
// information, and this should be accepted during the construction of the
// chain-specific transaction builder.
type TxBuilder interface {
	BuildTx([]Input, []Recipient) (Tx, error)
}

// The Client interface defines the functionality required to interact with a
// chain over RPC.
type Client interface {
	// Output returns the transaction output identified by the given outpoint.
	// It also returns the number of confirmations for the output. If the output
	// cannot be found before the context is done, or the output is invalid,
	// then an error should be returned. This method will not error, even if the
	// output has been spent.
	Output(context.Context, Outpoint) (Output, pack.U64, error)

	// UnspentOutput returns the unspent transaction output identified by the
	// given outpoint. It also returns the number of confirmations for the
	// output. If the output cannot be found before the context is done, the
	// output is invalid, or the output has been spent, then an error should be
	// returned.
	UnspentOutput(context.Context, Outpoint) (Output, pack.U64, error)

	// SubmitTx to the underlying chain. If the transaction cannot be found
	// before the context is done, or the transaction is invalid, then an error
	// should be returned.
	SubmitTx(context.Context, Tx) error
}

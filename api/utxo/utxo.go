package utxo

import (
	"context"

	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"
)

type Outpoint struct {
	Hash  pack.Bytes `json:"hash"`
	Index pack.U32   `json:"index"`
}

type Output struct {
	Outpoint
	Value        pack.U256  `json:"value"`
	PubKeyScript pack.Bytes `json:"pubKeyScript"`
}

type Input struct {
	Output
	SigScript pack.Bytes `json:"sigScript"`
}

type Recipient struct {
	To    address.Address `json:"to"`
	Value pack.U256       `json:"value"`
}

type Tx interface {
	// Hash returns the hash that uniquely identifies the transaction.
	// Generally, hashes are irreversible hash functions that consume the
	// content of the transaction.
	Hash() (pack.Bytes, error)

	// Inputs returns the transaction inputs consumed by the transaction.
	Inputs() ([]Input, error)

	// Outputs returns the transaction outputs produced by the transaction.
	Outputs() ([]Output, error)

	// Sighashes returns the digests that must be signed before the transaction
	// can be submitted by the client.
	Sighashes() ([]pack.Bytes32, error)

	// Sign the transaction by injecting signatures for the required sighashes.
	// The serialized public key used to sign the sighashes must also be
	// specified.
	Sign([]pack.Bytes65, pack.Bytes) error

	// Serialize the transaction into bytes. Generally, this is the format in
	// which the transaction will be submitted by the client.
	Serialize() (pack.Bytes, error)
}

type TxBuilder interface {
	BuildTx([]Input, []Recipient) (Tx, error)
}

type Client interface {
	// Output returns the transaction output uniquely identified by the given
	// transaction outpoint. It also returns the number of confirmations for the
	// transaction output.
	Output(context.Context, Outpoint) (Output, pack.U64, error)

	// SubmitTx to the underlying blockchain network.
	SubmitTx(context.Context, Tx) error
}

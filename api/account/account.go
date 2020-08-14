package account

import (
	"context"

	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"
)

type Tx interface {
	// Hash returns the hash that uniquely identifies the transaction.
	// Generally, hashes are irreversible hash functions that consume the
	// content of the transaction.
	Hash() pack.Bytes

	// From returns the address that is sending the transaction. Generally,
	// this is also the address that must sign the transaction.
	From() address.Address

	// To returns the address that is receiving the transaction. This can be the
	// address of an external account, controlled by a private key, or it can be
	// the address of a contract.
	To() address.Address

	// Value being sent from the sender to the receiver.
	Value() pack.U256

	// Nonce returns the nonce used to order the transaction with respect to all
	// other transactions signed and submitted by the sender.
	Nonce() pack.U256

	// Payload returns arbitrary data that is associated with the transaction.
	// Generally, this payload is used to send notes between external accounts,
	// or invoke business logic on a contract.
	Payload() pack.Bytes

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
	BuildTx(from, to address.Address, value, nonce pack.U256, payload pack.Bytes) (Tx, error)
}

type Client interface {
	// Tx returns the transaction uniquely identified by the given transaction
	// hash. It also returns the number of confirmations for the transaction.
	Tx(context.Context, pack.Bytes) (Tx, pack.U64, error)

	// SubmitTx to the underlying blockchain network.
	SubmitTx(context.Context, Tx) error
}

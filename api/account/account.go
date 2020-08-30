// Package account defines the Account API. All chains that use an account-based
// model should implement this API. The Account API is used to send and confirm
// transactions between addresses.
package account

import (
	"context"

	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"
)

// The Tx interfaces defines the functionality that must be exposed by
// account-based transactions.
type Tx interface {
	// Hash that uniquely identifies the transaction. Hashes are usually the
	// result of an irreversible hashing function applied to some serialized
	// representation of the transaction.
	Hash() pack.Bytes

	// From returns the address from which value is being sent.
	From() address.Address

	// To returns the address to which value is being sent.
	To() address.Address

	// Value being sent from one address to another.
	Value() pack.U256

	// Nonce used to order the transaction with respect to all other
	// transactions signed and submitted by the sender of this transaction.
	Nonce() pack.U256

	// Payload returns arbitrary data that is associated with the transaction.
	// This payload is often used to send notes between external accounts, or
	// call functions on a contract.
	Payload() contract.CallData

	// Sighash that must be signed before the transaction can be submitted by
	// the client.
	Sighash() (pack.Bytes32, error)

	// Sign the transaction by injecting signatures for the required sighashes.
	// The serialized public key used to sign the sighashes should also be
	// specified whenever it is available.
	Sign(pack.Bytes65, pack.Bytes) error

	// Serialize the transaction into bytes. This is the format in which the
	// transaction will be submitted by the client.
	Serialize() (pack.Bytes, error)
}

// The TxBuilder interface defines the functionality required to build
// account-based transactions. Most chain implementations require additional
// information, and this should be accepted during the construction of the
// chain-specific transaction builder.
type TxBuilder interface {
	BuildTx(from, to address.Address, value, nonce, gasLimit, gasPrice pack.U256, payload pack.Bytes) (Tx, error)
}

// The Client interface defines the functionality required to interact with a
// chain over RPC.
type Client interface {
	// Tx returns the transaction uniquely identified by the given transaction
	// hash. It also returns the number of confirmations for the transaction. If
	// the transaction cannot be found before the context is done, or the
	// transaction is invalid, then an error should be returned.
	Tx(context.Context, pack.Bytes) (Tx, pack.U64, error)

	// SubmitTx to the underlying chain. If the transaction cannot be found
	// before the context is done, or the transaction is invalid, then an error
	// should be returned.
	SubmitTx(context.Context, Tx) error
}

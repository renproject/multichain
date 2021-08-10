package evm

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"
	"math/big"
)

// TxBuilder represents a transaction builder that builds transactions to be
// broadcasted to the ethereum network. The TxBuilder is configured using a
// chain id.
type TxBuilder struct {
	ChainID *big.Int
}

// NewTxBuilder creates a new transaction builder.
func NewTxBuilder(chainID *big.Int) TxBuilder {
	return TxBuilder{chainID}
}

// BuildTx receives transaction fields and constructs a new transaction.
func (txBuilder TxBuilder) BuildTx(ctx context.Context, from, to address.Address, value, nonce, gasLimit, gasPrice, gasCap pack.U256, payload pack.Bytes) (account.Tx, error) {
	toAddr, err := NewAddressFromHex(string(pack.String(to)))
	if err != nil {
		return nil, fmt.Errorf("bad to address '%v': %v", to, err)
	}
	addr := common.Address(toAddr)
	return &Tx{
		ethTx: types.NewTransaction(nonce.Int().Uint64(),
			addr, value.Int(),
			gasLimit.Int().Uint64(),
			gasPrice.Int(),
			payload,
		),
		signer: types.LatestSignerForChainID(txBuilder.ChainID),
	}, nil
}

// Tx represents a ethereum transaction, encapsulating a payload/data and its
// signer.
type Tx struct {
	ethTx  *types.Transaction
	signer types.Signer
}

// Hash returns the hash that uniquely identifies the transaction.
// Generally, hashes are irreversible hash functions that consume the
// content of the transaction.
func (tx Tx) Hash() pack.Bytes {
	return tx.ethTx.Hash().Bytes()
}

// From returns the address that is sending the transaction. Generally,
// this is also the address that must sign the transaction.
func (tx Tx) From() address.Address {
	addr, err := types.Sender(tx.signer, tx.ethTx)
	if err != nil {
		return address.Address("")
	}
	return address.Address(addr.Hex())
}

// To returns the address that is receiving the transaction. This can be the
// address of an external account, controlled by a private key, or it can be
// the address of a contract.
func (tx Tx) To() address.Address {
	return address.Address(tx.ethTx.To().Hex())
}

// Value being sent from the sender to the receiver.
func (tx Tx) Value() pack.U256 {
	return pack.NewU256FromInt(tx.ethTx.Value())
}

// Nonce returns the nonce used to order the transaction with respect to all
// other transactions signed and submitted by the sender.
func (tx Tx) Nonce() pack.U256 {
	return pack.NewU256FromU64(pack.NewU64(tx.ethTx.Nonce()))
}

// Payload returns arbitrary data that is associated with the transaction.
// Generally, this payload is used to send notes between external accounts,
// or invoke business logic on a contract.
func (tx Tx) Payload() contract.CallData {
	return contract.CallData(tx.ethTx.Data())
}

// Sighashes returns the digests that must be signed before the transaction
// can be submitted by the client.
func (tx Tx) Sighashes() ([]pack.Bytes32, error) {
	return []pack.Bytes32{pack.Bytes32(tx.signer.Hash(tx.ethTx))}, nil
}

// Sign the transaction by injecting signatures for the required sighashes.
// The serialized public key used to sign the sighashes must also be
// specified.
func (tx *Tx) Sign(signatures []pack.Bytes65, pubkey pack.Bytes) error {
	signedtx, err := tx.ethTx.WithSignature(tx.signer, signatures[0].Bytes())
	if err != nil {
		return err
	}
	tx.ethTx = signedtx
	return nil
}

// Serialize the transaction into bytes. Generally, this is the format in
// which the transaction will be submitted by the client.
func (tx Tx) Serialize() (pack.Bytes, error) {
	return tx.ethTx.MarshalBinary()
}

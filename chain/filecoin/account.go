package filecoin

import (
	"bytes"
	"context"
	"fmt"

	filaddress "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/minio/blake2b-simd"
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"
)

// TxBuilder represents a transaction builder that builds transactions to be
// broadcasted to the filecoin network. The TxBuilder is configured using a
// gas price and gas limit.
type TxBuilder struct {
}

// NewTxBuilder creates a new transaction builder.
func NewTxBuilder() TxBuilder {
	return TxBuilder{}
}

// BuildTx receives transaction fields and constructs a new transaction.
func (txBuilder TxBuilder) BuildTx(ctx context.Context, from, to address.Address, value, nonce, gasLimit, gasPrice, gasCap pack.U256, payload pack.Bytes) (account.Tx, error) {
	filfrom, err := filaddress.NewFromString(string(from))
	if err != nil {
		return nil, fmt.Errorf("bad from address '%v': %v", from, err)
	}
	filto, err := filaddress.NewFromString(string(to))
	if err != nil {
		return nil, fmt.Errorf("bad to address '%v': %v", to, err)
	}
	methodNum := abi.MethodNum(0)
	return &Tx{
		msg: types.Message{
			Version:    types.MessageVersion,
			From:       filfrom,
			To:         filto,
			Value:      big.Int{Int: value.Int()},
			Nonce:      nonce.Int().Uint64(),
			GasFeeCap:  big.Int{Int: gasCap.Int()},
			GasLimit:   gasLimit.Int().Int64(),
			GasPremium: big.Int{Int: gasPrice.Int()},
			Method:     methodNum,
			Params:     payload,
		},
		signature: pack.Bytes65{},
	}, nil
}

// Tx represents a filecoin transaction, encapsulating a message and its
// signature.
type Tx struct {
	msg       types.Message
	signature pack.Bytes65
}

// Hash returns the hash that uniquely identifies the transaction.
// Generally, hashes are irreversible hash functions that consume the
// content of the transaction.
func (tx Tx) Hash() pack.Bytes {
	return pack.NewBytes(tx.msg.Cid().Bytes())
}

// From returns the address that is sending the transaction. Generally,
// this is also the address that must sign the transaction.
func (tx Tx) From() address.Address {
	return address.Address(tx.msg.From.String())
}

// To returns the address that is receiving the transaction. This can be the
// address of an external account, controlled by a private key, or it can be
// the address of a contract.
func (tx Tx) To() address.Address {
	return address.Address(tx.msg.To.String())
}

// Value being sent from the sender to the receiver.
func (tx Tx) Value() pack.U256 {
	return pack.NewU256FromInt(tx.msg.Value.Int)
}

// Nonce returns the nonce used to order the transaction with respect to all
// other transactions signed and submitted by the sender.
func (tx Tx) Nonce() pack.U256 {
	return pack.NewU256FromU64(pack.NewU64(tx.msg.Nonce))
}

// Payload returns arbitrary data that is associated with the transaction.
// Generally, this payload is used to send notes between external accounts,
// or invoke business logic on a contract.
func (tx Tx) Payload() contract.CallData {
	return contract.CallData(tx.msg.Params)
}

// Sighashes returns the digests that must be signed before the transaction
// can be submitted by the client.
func (tx Tx) Sighashes() ([]pack.Bytes32, error) {
	return []pack.Bytes32{pack.Bytes32(blake2b.Sum256(tx.Hash()))}, nil
}

// Sign the transaction by injecting signatures for the required sighashes.
// The serialized public key used to sign the sighashes must also be
// specified.
func (tx *Tx) Sign(signatures []pack.Bytes65, pubkey pack.Bytes) error {
	if len(signatures) != 1 {
		return fmt.Errorf("expected 1 signature, got %v signatures", len(signatures))
	}
	tx.signature = signatures[0]
	return nil
}

// Serialize the transaction into bytes. Generally, this is the format in
// which the transaction will be submitted by the client.
func (tx Tx) Serialize() (pack.Bytes, error) {
	buf := new(bytes.Buffer)
	if err := tx.msg.MarshalCBOR(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

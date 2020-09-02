package ethereum

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"
)

// The TxBuilder is an implementation of a Account-based chains compatible
// transaction builder for Ethereum.
type TxBuilder struct {
	config *params.ChainConfig
}

// NewTxBuilder returns a transaction builder that builds Account-compatible
// Ethereum transactions for the given chain configuration.
func NewTxBuilder(config *params.ChainConfig) TxBuilder {
	return TxBuilder{config: config}
}

// BuildTx returns an Ethereum transaction that transfers value from an address
// to another address, with other transaction-specific fields. The returned
// transaction implements the multichain.AccountTx interface.
func (txBuilder TxBuilder) BuildTx(
	from, to address.Address,
	value, nonce pack.U256,
	gasPrice, gasLimit pack.U256,
	payload pack.Bytes,
) (account.Tx, error) {
	toAddr, err := NewAddressFromHex(string(to))
	if err != nil {
		return nil, fmt.Errorf("decoding address: %v", err)
	}
	fromAddr, err := NewAddressFromHex(string(from))
	if err != nil {
		return nil, fmt.Errorf("decoding address: %v", err)
	}

	tx := types.NewTransaction(nonce.Int().Uint64(), common.Address(toAddr), value.Int(), gasLimit.Int().Uint64(), gasPrice.Int(), []byte(payload))

	signer := types.MakeSigner(txBuilder.config, nil)
	signed := false

	return &Tx{fromAddr, tx, signer, signed}, nil
}

// Tx represents an Ethereum transaction that implements the
// multichain.AccountTx interface
type Tx struct {
	from Address

	tx *types.Transaction

	signer types.Signer
	signed bool
}

// Hash returns the transaction hash of the given transaction.
func (tx *Tx) Hash() pack.Bytes {
	return pack.NewBytes(tx.tx.Hash().Bytes())
}

// From returns the sender of the transaction.
func (tx *Tx) From() address.Address {
	return address.Address(tx.from.String())
}

// To returns the recipient of the transaction.
func (tx *Tx) To() address.Address {
	return address.Address(tx.tx.To().String())
}

// Value returns the value (in native tokens) that is being transferred in the
// transaction.
func (tx *Tx) Value() pack.U256 {
	return pack.NewU256FromInt(tx.tx.Value())
}

// Nonce returns the transaction nonce for the transaction sender. This is a
// one-time use incremental identifier to protect against double spending.
func (tx *Tx) Nonce() pack.U256 {
	return pack.NewU256FromU64(pack.NewU64(tx.tx.Nonce()))
}

// Payload returns the data/payload attached in the transaction.
func (tx *Tx) Payload() contract.CallData {
	return contract.CallData(pack.NewBytes(tx.tx.Data()))
}

// Sighashes returns the digest that must be signed by the sender before the
// transaction can be submitted by the client.
func (tx *Tx) Sighashes() ([]pack.Bytes, error) {
	sighash32 := tx.signer.Hash(tx.tx)
	sighash := sighash32[:]
	return []pack.Bytes{pack.NewBytes(sighash)}, nil
}

// Sign consumes a list of signatures, and adds them to the underlying
// Ethereum transaction. In case of Ethereum, we expect only a single signature
// per transaction.
func (tx *Tx) Sign(signatures []pack.Bytes, pubKey pack.Bytes) error {
	if tx.signed {
		return fmt.Errorf("already signed")
	}

	if len(signatures) != 1 {
		return fmt.Errorf("expected 1 signature, got %v signatures", len(signatures))
	}

	if len(signatures[0]) != 65 {
		return fmt.Errorf("expected signature to be 65 bytes, got %v bytes", len(signatures[0]))
	}

	signedTx, err := tx.tx.WithSignature(tx.signer, []byte(signatures[0]))
	if err != nil {
		return err
	}

	tx.tx = signedTx
	tx.signed = true
	return nil
}

// Serialize serializes the transaction to bytes.
func (tx *Tx) Serialize() (pack.Bytes, error) {
	// FIXME: I am pretty sure that this is not the format the Ethereum expects
	// transactions to be serialised to on the network. Although the client
	// might expect to send JSON objects, that is different from serialization,
	// and can better represented by implementing MarshalJSON on this type.
	serialized, err := tx.tx.MarshalJSON()
	if err != nil {
		return pack.Bytes{}, err
	}

	return pack.NewBytes(serialized), nil
}

// EthClient interacts with an instance of the Ethereum network using the RPC
// interface exposed by an Ethereum node.
type EthClient struct {
	client *ethclient.Client
}

// NewClient returns a new Client.
func NewClient(rpcURL pack.String) (account.Client, error) {
	client, err := ethclient.Dial(string(rpcURL))
	if err != nil {
		return nil, fmt.Errorf("dialing RPC URL %v: %v", rpcURL, err)
	}

	return EthClient{client}, nil
}

// Tx queries the Ethereum node to fetch a transaction with the provided tx ID
// and also returns the number of block confirmations for the transaction.
func (client EthClient) Tx(ctx context.Context, txID pack.Bytes) (account.Tx, pack.U64, error) {
	txHash := common.BytesToHash(txID)
	tx, isPending, err := client.client.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, pack.NewU64(0), fmt.Errorf("fetching tx: %v", err)
	}
	if isPending {
		return nil, pack.NewU64(0), fmt.Errorf("tx not confirmed")
	}
	txReceipt, err := client.client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, pack.NewU64(0), fmt.Errorf("fetching tx receipt: %v", err)
	}
	block, err := client.client.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, pack.NewU64(0), fmt.Errorf("fetching current block: %v", err)
	}
	confs := block.NumberU64() - txReceipt.BlockNumber.Uint64() + 1

	return &Tx{tx: tx}, pack.NewU64(confs), nil
}

// SubmitTx submits a signed transaction to the Ethereum network.
func (client EthClient) SubmitTx(ctx context.Context, tx account.Tx) error {
	panic("unimplemented")
}

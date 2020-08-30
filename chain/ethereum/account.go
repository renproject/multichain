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

type TxBuilder struct {
	config *params.ChainConfig
}

func NewTxBuilder(config *params.ChainConfig) TxBuilder {
	return TxBuilder{config: config}
}

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

type Tx struct {
	from Address

	tx *types.Transaction

	signer types.Signer
	signed bool
}

func (tx *Tx) Hash() pack.Bytes {
	return pack.NewBytes(tx.tx.Hash().Bytes())
}

func (tx *Tx) From() address.Address {
	return address.Address(tx.from.String())
}

func (tx *Tx) To() address.Address {
	return address.Address(tx.tx.To().String())
}

func (tx *Tx) Value() pack.U256 {
	return pack.NewU256FromInt(tx.tx.Value())
}

func (tx *Tx) Nonce() pack.U256 {
	return pack.NewU256FromU64(pack.NewU64(tx.tx.Nonce()))
}

func (tx *Tx) Payload() contract.CallData {
	return contract.CallData(pack.NewBytes(tx.tx.Data()))
}

func (tx *Tx) Sighash() (pack.Bytes32, error) {
	sighash := tx.signer.Hash(tx.tx)
	return pack.NewBytes32(sighash), nil
}

func (tx *Tx) Sign(signature pack.Bytes65, pubKey pack.Bytes) error {
	if tx.signed {
		return fmt.Errorf("already signed")
	}

	signedTx, err := tx.tx.WithSignature(tx.signer, signature.Bytes())
	if err != nil {
		return err
	}

	tx.tx = signedTx
	tx.signed = true
	return nil
}

func (tx *Tx) Serialize() (pack.Bytes, error) {
	serialized, err := tx.tx.MarshalJSON()
	if err != nil {
		return pack.Bytes{}, err
	}

	return pack.NewBytes(serialized), nil
}

type EthClient struct {
	client *ethclient.Client
}

func NewClient(rpcURL pack.String) (account.Client, error) {
	client, err := ethclient.Dial(string(rpcURL))
	if err != nil {
		return nil, fmt.Errorf("dialing RPC URL %v: %v", rpcURL, err)
	}

	return EthClient{client}, nil
}

func (client EthClient) Tx(ctx context.Context, txId pack.Bytes) (account.Tx, pack.U64, error) {
	txHash := common.BytesToHash(txId)
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

func (client EthClient) SubmitTx(ctx context.Context, tx account.Tx) error {
	panic("unimplemented")
}

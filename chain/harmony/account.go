package harmony

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/harmony-one/harmony/core/types"
	common2 "github.com/harmony-one/harmony/rpc/common"
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"
	"math/big"
)

const (
	DefaultShardID = 1
	DefaultHost    = "http://127.0.0.1:9598"
)

type TxBuilder struct {
	chainID *big.Int
}

func NewTxBuilder(chainId *big.Int) account.TxBuilder {
	return &TxBuilder{
		chainID: chainId,
	}
}

func (txBuilder *TxBuilder) BuildTx(ctx context.Context, from, to address.Address, value, nonce, gasLimit, gasPrice, gasCap pack.U256, payload pack.Bytes) (account.Tx, error) {
	toAddr, err := NewEncoderDecoder().DecodeAddress(to)
	if err != nil {
		return nil, err
	}
	tx := types.NewTransaction(
		nonce.Int().Uint64(),
		common.BytesToAddress(toAddr),
		DefaultShardID,
		value.Int(),
		gasLimit.Int().Uint64(),
		gasPrice.Int(),
		payload)
	return &Tx{
		harmonyTx: *tx,
		chainId:   txBuilder.chainID,
		sender:    from,
		signed:    false,
	}, nil
}

type TxData struct {
	Blockhash        string   `json:"blockHash"`
	Blocknumber      uint64   `json:"blockNumber"`
	From             string   `json:"from"`
	Gas              uint64   `json:"gas"`
	Gasprice         *big.Int `json:"gasPrice"`
	Hash             string   `json:"hash"`
	Input            string   `json:"input"`
	Nonce            uint64   `json:"nonce"`
	R                string   `json:"r"`
	S                string   `json:"s"`
	Shardid          uint32   `json:"shardID"`
	Timestamp        uint64   `json:"timestamp"`
	To               string   `json:"to"`
	Toshardid        uint32   `json:"toShardID"`
	Transactionindex uint64   `json:"transactionIndex"`
	V                string   `json:"v"`
	Value            *big.Int `json:"value"`
}

type Tx struct {
	harmonyTx types.Transaction
	chainId   *big.Int
	sender    address.Address
	signed    bool
}

func (tx *Tx) Hash() pack.Bytes {
	return pack.NewBytes(tx.harmonyTx.Hash().Bytes())
}

func (tx *Tx) From() address.Address {
	from, err := tx.harmonyTx.SenderAddress()
	if err == nil {
		addr, err := NewEncoderDecoder().EncodeAddress(from.Bytes())
		if err == nil {
			return addr
		}
	}
	return tx.sender
}

func (tx *Tx) To() address.Address {
	to := tx.harmonyTx.To()
	if to != nil {
		addr, err := NewEncoderDecoder().EncodeAddress(to.Bytes())
		if err == nil {
			return addr
		}
	}
	return ""
}

func (tx *Tx) Value() pack.U256 {
	return pack.NewU256FromInt(tx.harmonyTx.Value())
}

func (tx *Tx) Nonce() pack.U256 {
	return pack.NewU256FromU64(pack.NewU64(tx.harmonyTx.Nonce()))
}

func (tx *Tx) Payload() contract.CallData {
	return tx.harmonyTx.Data()
}

func (tx *Tx) Sighashes() ([]pack.Bytes32, error) {
	const digestLength = 32
	var (
		digestHash [32]byte
		sighashes  []pack.Bytes32
	)
	h := types.NewEIP155Signer(tx.chainId).Hash(&tx.harmonyTx).Bytes()
	if len(h) != digestLength {
		return nil, fmt.Errorf("hash is required to be exactly %d bytes (%d)", digestLength, len(h))
	}
	copy(digestHash[:], h[:32])
	sighashes = append(sighashes, digestHash)
	return sighashes, nil
}

func (tx *Tx) Sign(signatures []pack.Bytes65, pubKey pack.Bytes) error {
	if len(signatures) != 1 {
		return fmt.Errorf("expected 1 signature, got %v signatures", len(signatures))
	}
	signedTx, err := tx.harmonyTx.WithSignature(types.NewEIP155Signer(tx.chainId), signatures[0].Bytes())
	if err != nil {
		return err
	}
	tx.harmonyTx = *signedTx
	tx.signed = true
	return nil
}

func (tx *Tx) Serialize() (pack.Bytes, error) {
	serializedTx, err := rlp.EncodeToBytes(&tx.harmonyTx)
	if err != nil {
		return pack.Bytes{}, err
	}
	return pack.NewBytes(serializedTx), nil
}

type ClientOptions struct {
	Host string
}

type Client struct {
	opts ClientOptions
}

func (opts ClientOptions) WithHost(host string) ClientOptions {
	opts.Host = host
	return opts
}

func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		Host: DefaultHost,
	}
}

func NewClient(opts ClientOptions) *Client {
	return &Client{opts: opts}
}

func (c *Client) LatestBlock(ctx context.Context) (pack.U64, error) {
	for {
		select {
		case <-ctx.Done():
			return pack.NewU64(0), ctx.Err()
		default:
		}
		const method = "hmyv2_blockNumber"
		response, err := SendData(method, []byte{}, c.opts.Host)
		if err != nil {
			fmt.Println(err)
			return pack.NewU64(0), err
		}
		var latestBlock uint64
		if err := json.Unmarshal(*response.Result, &latestBlock); err != nil {
			return pack.NewU64(0), fmt.Errorf("decoding result: %v", err)
		}
		return pack.NewU64(latestBlock), nil
	}
}

func (c *Client) AccountBalance(ctx context.Context, addr address.Address) (pack.U256, error) {
	for {
		select {
		case <-ctx.Done():
			return pack.U256{}, ctx.Err()
		default:
		}
		data := []byte(fmt.Sprintf("[\"%s\"]", addr))
		const method = "hmyv2_getBalance"
		response, err := SendData(method, data, c.opts.Host)
		if err != nil {
			fmt.Println(err)
			return pack.U256{}, err
		}
		var balance uint64
		if err := json.Unmarshal(*response.Result, &balance); err != nil {
			return pack.U256{}, fmt.Errorf("decoding result: %v", err)
		}
		return pack.NewU256FromU64(pack.NewU64(balance)), nil
	}
}

func (c *Client) AccountNonce(ctx context.Context, addr address.Address) (pack.U256, error) {
	for {
		select {
		case <-ctx.Done():
			return pack.U256{}, ctx.Err()
		default:
		}
		data := []byte(fmt.Sprintf("[\"%s\", \"%s\"]", addr, "SENT"))
		const method = "hmyv2_getTransactionsCount"
		response, err := SendData(method, data, c.opts.Host)
		if err != nil {
			fmt.Println(err)
			return pack.U256{}, err
		}
		var nonce uint64
		if err := json.Unmarshal(*response.Result, &nonce); err != nil {
			return pack.U256{}, fmt.Errorf("decoding result: %v", err)
		}
		return pack.NewU256FromU64(pack.NewU64(nonce)), nil
	}
}

func (c *Client) Tx(ctx context.Context, hash pack.Bytes) (account.Tx, pack.U64, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, pack.NewU64(0), ctx.Err()
		default:
		}
		data := []byte(fmt.Sprintf("[\"%s\"]", hexutil.Encode(hash)))
		const method = "hmyv2_getTransactionByHash"
		response, err := SendData(method, data, c.opts.Host)
		if err != nil {
			return nil, pack.NewU64(0), err
		}
		var txData TxData
		if response.Result == nil {
			return nil, pack.NewU64(0), fmt.Errorf("decoding result: %v", err)
		}
		if err := json.Unmarshal(*response.Result, &txData); err != nil {
			return nil, pack.NewU64(0), fmt.Errorf("decoding result: %v", err)
		}

		tx, err := buildTxFromTxData(txData)
		return tx, pack.NewU64(1), err
	}
}

func (c *Client) SubmitTx(ctx context.Context, tx account.Tx) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		txSerilized, err := tx.Serialize()
		if err != nil {
			return err
		}
		hexSignature := hexutil.Encode(txSerilized)
		data := []byte(fmt.Sprintf("[\"%s\"]", hexSignature))
		const method = "hmyv2_sendRawTransaction"
		tx1 := new(types.Transaction)
		err = rlp.DecodeBytes(txSerilized, tx1)
		_, err = SendData(method, data, c.opts.Host)
		if err != nil {
			return err
		}
		return nil
	}
}

func (c *Client) ChainId(ctx context.Context) (*big.Int, error) {
	for {
		select {
		case <-ctx.Done():
			return big.NewInt(0), ctx.Err()
		default:
		}
		const method = "hmyv2_getNodeMetadata"
		response, err := SendData(method, []byte{}, c.opts.Host)
		if err != nil {
			fmt.Println(err)
			return big.NewInt(0), err
		}
		var nodeMetadata common2.NodeMetadata
		if err := json.Unmarshal(*response.Result, &nodeMetadata); err != nil {
			return big.NewInt(0), fmt.Errorf("decoding result: %v", err)
		}
		return nodeMetadata.ChainConfig.ChainID, nil
	}
}

func buildTxFromTxData(data TxData) (account.Tx, error) {
	toAddr, err := NewEncoderDecoder().DecodeAddress(address.Address(data.To))
	if err != nil {
		return nil, err
	}
	tx := types.NewTransaction(
		data.Nonce,
		common.BytesToAddress(toAddr),
		data.Shardid,
		data.Value,
		data.Gas,
		data.Gasprice,
		pack.Bytes(nil),
	)
	return &Tx{
		harmonyTx: *tx,
		sender:    address.Address(data.From),
		signed:    true,
	}, nil
}
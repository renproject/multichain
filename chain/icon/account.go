package icon

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/icon-project/goloop/server/jsonrpc"
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/multichain/chain/icon/crypto"
	"github.com/renproject/multichain/chain/icon/intconv"
	"github.com/renproject/multichain/chain/icon/transaction"
	"github.com/renproject/pack"
	"strconv"
	"time"
)

type RawMessage []byte

// Tx ...
type Tx struct {
	Version     jsonrpc.HexInt  `json:"version" validate:"required,t_int"`
	FromAddress address.Address `json:"from" validate:"required,t_addr_eoa"`
	ToAddress   address.Address `json:"to" validate:"required,t_addr"`
	Amount      jsonrpc.HexInt  `json:"value,omitempty" validate:"optional,t_int"`
	StepLimit   jsonrpc.HexInt  `json:"stepLimit" validate:"required,t_int"`
	Timestamp   jsonrpc.HexInt  `json:"timestamp" validate:"required,t_int"`
	NID         jsonrpc.HexInt  `json:"nid" validate:"required,t_int"`
	Salt        jsonrpc.HexInt  `json:"nonce,omitempty" validate:"optional,t_int"`
	Signature   string          `json:"signature" validate:"required,t_sig"`
	DataType    *string         `json:"dataType,omitempty" validate:"optional,call|deploy|message"`
	Data        RawMessage      `json:"data,omitempty"`
}

// PrivateKey is a type representing a private key.
// for both private key and public key
type PrivateKey pack.Bytes

// Hash ...
func (tx *Tx) Hash() pack.Bytes {
	h, err := tx.Sighashes()
	if err != nil {
		return pack.Bytes{}
	} else {
		b := h[0][:]
		return pack.NewBytes(b)
	}
}

// To ...
func (tx *Tx) To() address.Address {
	return tx.ToAddress
}

// From ...
func (tx *Tx) From() address.Address {
	return tx.FromAddress
}

// Value ...
func (tx *Tx) Value() pack.U256 {
	output, _ := strconv.ParseInt(string(tx.Amount)[2:], 16, 64)
	return pack.NewU256FromU64(pack.NewU64(uint64(output)))
}

// Nonce ...
func (tx *Tx) Nonce() pack.U256 {
	output, _ := strconv.ParseInt(string(tx.Salt)[2:], 16, 64)
	return pack.NewU256FromU64(pack.NewU64(uint64(output)))
}

// Serialize ...
func (tx *Tx) Serialize() (pack.Bytes, error) {
	txSerializeExcludes := map[string]bool{"signature": true}
	//tx.Timestamp = jsonrpc.HexInt(intconv.FormatInt(time.Now().UnixNano() / int64(time.Microsecond)))
	js, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	return transaction.SerializeJSON(js, nil, txSerializeExcludes)
}

// Sighashes returns the digests that must be signed before the transaction
// can be submitted by the client.
func (tx *Tx) Sighashes() ([]pack.Bytes32, error) {
	bs, _ := tx.Serialize()
	bs = append([]byte("icx_sendTransaction."), bs...)
	sighashes := make([]pack.Bytes32, 1)
	d := crypto.SHA3Sum256(bs)
	sighash := [32]byte{}
	copy(sighash[:], d)
	sighashes[0] = pack.NewBytes32(sighash)
	return sighashes, nil
}

// Payload returns the memo attached to the transaction.
func (tx *Tx) Payload() contract.CallData {
	return contract.CallData(pack.Bytes(make([]byte, 0)))
}
func (tx *Tx) Sign(signatures []pack.Bytes65, pubKey pack.Bytes) error {
	if len(signatures) != 1 {
		return fmt.Errorf("expected 1 signature, got %v signatures", len(signatures))
	}
	tx.Signature = base64.StdEncoding.EncodeToString(signatures[0].Bytes())
	return nil
}

// TxBuilder represents a transaction builder that builds transactions to be
// broadcasted to the icon network. The TxBuilder is configured using a
// gas price and gas limit.
type TxBuilder struct {
	chainID string
}

// NewTxBuilder creates a new transaction builder.
func NewTxBuilder(chainId string) TxBuilder {
	return TxBuilder{chainID: chainId}
}

// BuildTx ...
func (txBuilder TxBuilder) BuildTx(from, to address.Address, value pack.U256, nonce pack.U256, stepLimit pack.U256, gasPrice pack.U256, payload pack.Bytes) (account.Tx, error) {

	stepHex := fmt.Sprintf("0x%x", stepLimit.Int())
	valueHex := fmt.Sprintf("0x%x", value.Int())
	tx := Tx{
		Version:     "0x3",
		FromAddress: from,
		ToAddress:   to,
		StepLimit:   jsonrpc.HexInt(stepHex),
		Amount:      jsonrpc.HexInt(valueHex),
		Timestamp:   jsonrpc.HexInt(intconv.FormatInt(time.Now().UnixNano() / int64(time.Microsecond))),
		NID:         jsonrpc.HexInt(txBuilder.chainID),
		Salt:        "0x1",
	}
	return &tx, nil
}

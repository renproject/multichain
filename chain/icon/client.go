package icon

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/icon-project/goloop/client"
	v3 "github.com/icon-project/goloop/server/v3"
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"log"
	"net/http"
	"time"

	"github.com/icon-project/goloop/server/jsonrpc"
	"github.com/renproject/pack"
)

const jsonprcVersion = "2.0"

// Request ...
type Request struct {
	Version string          `json:"jsonrpc" validate:"required,version"`
	Method  string          `json:"method" validate:"required"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id"`
}

// Response ...
type Response struct {
	Version string      `json:"jsonrpc" validate:"required,version"`
	ID      interface{} `json:"id"`
	Result  Tx          `json:"result"`
}

// Client ...
type Client struct {
	v3       client.ClientV3
	endpoint string
}

// NewClient returns a new Client.
func NewClient(endpoint pack.String) *Client {

	return &Client{
		v3:       *client.NewClientV3(endpoint.String()),
		endpoint: endpoint.String(),
	}
}

// Tx ...
func (ct *Client) Tx(ctx context.Context, txHash pack.Bytes) (account.Tx, pack.U64, error) {
	txResult, err := ct.v3.GetTransactionByHash(&v3.TransactionHashParam{Hash: jsonrpc.HexBytes(txHash)})
	if err != nil {
		return nil, pack.NewU64(0), err
	}
	return &Tx{
		Version:     txResult.Version,
		Amount:      txResult.Value,
		FromAddress: address.Address(txResult.From.Address().String()),
		ToAddress:   address.Address(txResult.To.Address().String()),
		StepLimit:   txResult.StepLimit,
		Signature:   string(txResult.Signature.Bytes()),
		Timestamp:   txResult.TimeStamp,
		NID:         txResult.NID,
		DataType:    &txResult.DataType,
	}, pack.NewU64(1), nil
}

// SubmitTx ...
func (ct *Client) SubmitTx(ctx context.Context, tx account.Tx) error {
	_, err := ct.sendTransaction(tx)
	return err
}

func (ct Client) sendTransaction(tx account.Tx) (*http.Response, error) {
	rq := &Request{
		Version: jsonprcVersion,
		Method:  "icx_sendTransaction",
		ID:      time.Now().UnixNano() / int64(time.Millisecond),
	}
	b, _ := json.Marshal(tx)
	rq.Params = json.RawMessage(b)
	reqB, err := json.Marshal(rq)
	if err != nil {
		return nil, err
	}
	res, err := http.Post(ct.endpoint, "application/json; charset=utf-8", bytes.NewReader(reqB))
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	return res, nil
}

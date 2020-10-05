package icon

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/icon-project/goloop/common"
	"github.com/icon-project/goloop/common/errors"
	"github.com/icon-project/goloop/server/jsonrpc"
	"github.com/renproject/multichain/chain/icon/intconv"
	"github.com/renproject/multichain/chain/icon/transaction"
	"github.com/renproject/pack"
)

// Client interacts with an instance of ICON network using the REST
// interface exposed by a client node.
type Client struct {
	*JsonRpcClient
	Debug *JsonRpcClient
	conns map[string]*websocket.Conn
}

// NewClient returns a new Client.
func NewClient(endpoint string) *Client {
	client := new(http.Client)
	apiClient := NewJsonRpcClient(client, endpoint)
	var debugClient *JsonRpcClient
	if ep := guessDebugEndpoint(endpoint); len(ep) > 0 {
		debugClient = NewJsonRpcClient(client, ep)
	}

	return &Client{
		JsonRpcClient: apiClient,
		Debug:         debugClient,
		conns:         make(map[string]*websocket.Conn),
	}
}

// Tx query transaction with txHash
func (client Client) Tx(ctx context.Context, hash pack.Bytes) (Tx, pack.U64, error) {
	t := &Transaction{}
	_, err := client.Do("icx_getTransactionByHash", hash, t)
	if err != nil {
		return nil, err
	}
	return t, pack.NewU64(1), nil
}

// SubmitTx to ICON network.
func (client Client) SubmitTx(ctx context.Context, tx Tx) error {
	Tx.Timestamp = jsonrpc.HexInt(intconv.FormatInt(time.Now().UnixNano() / int64(time.Microsecond)))
	js, err := json.Marshal(Tx)
	if err != nil {
		return nil, err
	}

	bs, err := transaction.SerializeJSON(js, nil, txSerializeExcludes)
	if err != nil {
		return nil, err
	}
	bs = append([]byte("icx_sendTransaction."), bs...)

	var result jsonrpc.HexBytes
	if _, err = client.Do("icx_sendTransaction", Tx, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// EstimateStep cost for a transaction
func (client *Client) EstimateStep(param Tx) (*common.HexInt, error) {
	if client.Debug == nil {
		return nil, errors.InvalidStateError.New("UnavailableDebugEndPoint")
	}
	param.Timestamp = jsonrpc.HexInt(intconv.FormatInt(time.Now().UnixNano() / int64(time.Microsecond)))
	var result common.HexInt
	if _, err := client.Debug.Do("debug_estimateStep", param, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

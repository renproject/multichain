package filecoin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api"
	filclient "github.com/filecoin-project/lotus/api/client"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/actors/crypto"
	"github.com/ipfs/go-cid"
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/pack"
)

const (
	// AuthorizationKey is the header key used for authorization
	AuthorizationKey = "Authorization"

	// DefaultClientRpcURL is the RPC websocket URL used by default, to
	// interact with the filecoin lotus node.
	DefaultClientRpcURL = "ws://127.0.0.1:1234/rpc/v0"

	// DefaultClientAuthToken is the auth token used to instantiate the lotus
	// client. A valid lotus auth token is required to write messages to the
	// filecoin storage. To do read-only queries, auth token is not required.
	DefaultClientAuthToken = ""
)

// ClientOptions are used to parameterise the behaviour of the Client.
type ClientOptions struct {
	RpcURL    string
	AuthToken string
}

// DefaultClientOptions returns ClientOptions with the default settings. These
// settings are valid for use with the default local deployment of the
// multichain. In production, the rpc-url and authentication token should be
// changed.
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		RpcURL:    DefaultClientRpcURL,
		AuthToken: DefaultClientAuthToken,
	}
}

// WithRpcURL returns a modified version of the options with the given API
// rpc-url
func (opts ClientOptions) WithRpcURL(rpcURL string) ClientOptions {
	opts.RpcURL = rpcURL
	return opts
}

// WithAuthToken returns a modified version of the options with the given
// authentication token.
func (opts ClientOptions) WithAuthToken(authToken string) ClientOptions {
	opts.AuthToken = authToken
	return opts
}

// Client holds options to connect to a filecoin lotus node, and the underlying
// RPC client instance.
type Client struct {
	opts   ClientOptions
	node   api.FullNode
	closer jsonrpc.ClientCloser
}

// NewClient creates and returns a new JSON-RPC client to the Filecoin node
func NewClient(opts ClientOptions) (*Client, error) {
	requestHeaders := make(http.Header)
	if opts.AuthToken != DefaultClientAuthToken {
		requestHeaders.Add(AuthorizationKey, opts.AuthToken)
	}

	node, closer, err := filclient.NewFullNodeRPC(context.Background(), opts.RpcURL, requestHeaders)
	if err != nil {
		return nil, err
	}

	return &Client{opts, node, closer}, nil
}

// Tx returns the transaction uniquely identified by the given transaction
// hash. It also returns the number of confirmations for the transaction.
func (client *Client) Tx(ctx context.Context, txID pack.Bytes) (account.Tx, pack.U64, error) {
	// parse the transaction ID to a message ID
	msgID, err := cid.Parse([]byte(txID))
	if err != nil {
		return nil, pack.NewU64(0), fmt.Errorf("parsing txid: %v", err)
	}

	// lookup message receipt to get its height
	messageLookup, err := client.node.StateSearchMsg(ctx, msgID)
	if err != nil {
		return nil, pack.NewU64(0), fmt.Errorf("searching state for txid: %v", err)
	}
	if messageLookup == nil {
		return nil, pack.NewU64(0), fmt.Errorf("searching state for txid %v: not found", msgID)
	}

	// get the most recent tipset and its height
	headTipset, err := client.node.ChainHead(ctx)
	if err != nil {
		return nil, pack.NewU64(0), fmt.Errorf("getting head from chain: %v", err)
	}
	confs := headTipset.Height() - messageLookup.Height + 1
	if confs < 0 {
		return nil, pack.NewU64(0), fmt.Errorf("getting head from chain: negative confirmations")
	}

	// get the message
	msg, err := client.node.ChainGetMessage(ctx, msgID)
	if err != nil {
		return nil, pack.NewU64(0), fmt.Errorf("getting txid %v from chain: %v", msgID, err)
	}

	return &Tx{msg: *msg}, pack.NewU64(uint64(confs)), nil
}

// SubmitTx to the underlying blockchain network.
func (client *Client) SubmitTx(ctx context.Context, tx account.Tx) error {
	switch tx := tx.(type) {
	case *Tx:
		// construct crypto.Signature
		signature := crypto.Signature{
			Type: crypto.SigTypeSecp256k1,
			Data: tx.signature.Bytes(),
		}

		// construct types.SignedMessage
		signedMessage := types.SignedMessage{
			Message:   tx.msg,
			Signature: signature,
		}

		// submit transaction to mempool
		msgID, err := client.node.MpoolPush(ctx, &signedMessage)
		if err != nil {
			return fmt.Errorf("pushing txid %v to mpool: %v", msgID, err)
		}
		return nil
	default:
		return fmt.Errorf("expected type %T, got type %T", new(Tx), tx)
	}
}

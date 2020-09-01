package filecoin

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	filaddress "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api"
	filclient "github.com/filecoin-project/lotus/api/client"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/filecoin-project/specs-actors/actors/crypto"
	"github.com/ipfs/go-cid"
	"github.com/minio/blake2b-simd"
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"
)

const (
	// AuthorizationKey is the header key used for authorization
	AuthorizationKey = "Authorization"

	// DefaultClientMultiAddress is the RPC websocket URL used by default, to
	// interact with the filecoin lotus node.
	DefaultClientMultiAddress = "ws://127.0.0.1:1234/rpc/v0"

	// DefaultClientAuthToken is the auth token used to instantiate the lotus
	// client. A valid lotus auth token is required to write messages to the
	// filecoin storage. To do read-only queries, auth token is not required.
	DefaultClientAuthToken = ""
)

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
	return pack.NewBytes(tx.msg.Cid().Hash())
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
	if tx.msg.Method == 0 {
		if len(tx.msg.Params) == 0 {
			return contract.CallData([]byte{})
		}
		return contract.CallData(append([]byte{0}, tx.msg.Params...))
	}
	if len(tx.msg.Params) == 0 {
		return contract.CallData([]byte{byte(tx.msg.Method)})
	}
	return contract.CallData(append([]byte{byte(tx.msg.Method)}, tx.msg.Params...))
}

// Sighashes returns the digests that must be signed before the transaction
// can be submitted by the client.
func (tx Tx) Sighashes() ([]pack.Bytes32, error) {
	return []pack.Bytes32{blake2b.Sum256(tx.Hash())}, nil
}

// Sign the transaction by injecting signatures for the required sighashes.
// The serialized public key used to sign the sighashes must also be
// specified.
func (tx Tx) Sign(signatures []pack.Bytes65, pubkey pack.Bytes) error {
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

// TxBuilder represents a transaction builder that builds transactions to be
// broadcasted to the filecoin network. The TxBuilder is configured using a
// gas price and gas limit.
type TxBuilder struct {
	gasPrice pack.U256
	gasLimit pack.U256
}

// NewTxBuilder creates a new transaction builder.
func NewTxBuilder(gasPrice, gasLimit pack.U256) TxBuilder {
	return TxBuilder{gasPrice: gasPrice, gasLimit: gasLimit}
}

// BuildTx receives transaction fields and constructs a new transaction.
func (txBuilder TxBuilder) BuildTx(from, to address.Address, value, nonce, _, _ pack.U256, payload pack.Bytes) (account.Tx, error) {
	filfrom, err := filaddress.NewFromString(string(from))
	if err != nil {
		return nil, fmt.Errorf("bad from address '%v': %v", from, err)
	}
	filto, err := filaddress.NewFromString(string(to))
	if err != nil {
		return nil, fmt.Errorf("bad to address '%v': %v", to, err)
	}
	methodNum := abi.MethodNum(0)
	if len(payload) > 0 {
		methodNum = abi.MethodNum(payload[0])
		payload = payload[1:]
	}
	return Tx{
		msg: types.Message{
			Version:    types.MessageVersion,
			From:       filfrom,
			To:         filto,
			Value:      big.Int{Int: value.Int()},
			Nonce:      value.Int().Uint64(),
			GasFeeCap:  big.Int{Int: txBuilder.gasPrice.Int()},
			GasPremium: big.Int{Int: pack.NewU256([32]byte{}).Int()},
			GasLimit:   txBuilder.gasLimit.Int().Int64(),
			Method:     methodNum,
			Params:     payload,
		},
		signature: pack.Bytes65{},
	}, nil
}

// ClientOptions are used to parameterise the behaviour of the Client.
type ClientOptions struct {
	MultiAddress string
	AuthToken    string
}

// DefaultClientOptions returns ClientOptions with the default settings. These
// settings are valid for use with the default local deployment of the
// multichain. In production, the multi-address and authentication token should
// be changed.
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		MultiAddress: DefaultClientMultiAddress,
		AuthToken:    DefaultClientAuthToken,
	}
}

// WithAddress returns a modified version of the options with the given API
// multi-address.
func (opts ClientOptions) WithAddress(multiAddr string) ClientOptions {
	opts.MultiAddress = multiAddr
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

	node, closer, err := filclient.NewFullNodeRPC(context.Background(), opts.MultiAddress, requestHeaders)
	if err != nil {
		return nil, err
	}

	return &Client{opts, node, closer}, nil
}

// Tx returns the transaction uniquely identified by the given transaction
// hash. It also returns the number of confirmations for the transaction.
func (client *Client) Tx(ctx context.Context, txID pack.Bytes) (account.Tx, pack.U64, error) {
	// parse the transaction ID to a message ID
	msgID, err := cid.Parse(txID.String())
	if err != nil {
		return nil, pack.NewU64(0), fmt.Errorf("parsing txID: %v", err)
	}

	// lookup message receipt to get its height
	messageLookup, err := client.node.StateSearchMsg(ctx, msgID)
	if err != nil {
		return nil, pack.NewU64(0), fmt.Errorf("searching msg: %v", err)
	}

	// get the most recent tipset and its height
	headTipset, err := client.node.ChainHead(ctx)
	if err != nil {
		return nil, pack.NewU64(0), fmt.Errorf("fetching head: %v", err)
	}
	confs := headTipset.Height() - messageLookup.Height + 1

	// get the message
	msg, err := client.node.ChainGetMessage(ctx, msgID)
	if err != nil {
		return nil, pack.NewU64(0), fmt.Errorf("fetching msg: %v", err)
	}

	return &Tx{msg: *msg}, pack.NewU64(uint64(confs)), nil
}

// SubmitTx to the underlying blockchain network.
// TODO: should also return a transaction hash (pack.Bytes) ?
func (client *Client) SubmitTx(ctx context.Context, tx account.Tx) error {
	switch tx := tx.(type) {
	case Tx:
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
		_, err := client.node.MpoolPush(ctx, &signedMessage)
		if err != nil {
			return fmt.Errorf("pushing msg to mpool: %v", err)
		}
		return nil
	default:
		return fmt.Errorf("invalid tx type: %v", tx)
	}
}

package filecoin

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"

	filaddress "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/cli"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/minio/blake2b-simd"
	"github.com/multiformats/go-multiaddr"
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"
)

const (
	DefaultClientMultiAddress = ""
	DefaultClientAuthToken    = ""
)

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
func (tx Tx) Payload() pack.Bytes {
	if tx.msg.Method == 0 {
		if len(tx.msg.Params) == 0 {
			return pack.NewBytes([]byte{})
		}
		return pack.NewBytes(append([]byte{0}, tx.msg.Params...))
	}
	if len(tx.msg.Params) == 0 {
		return pack.NewBytes([]byte{byte(tx.msg.Method)})
	}
	return pack.NewBytes(append([]byte{byte(tx.msg.Method)}, tx.msg.Params...))
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

type TxBuilder struct {
	gasPrice pack.U256
	gasLimit pack.U256
}

func (txBuilder TxBuilder) BuildTx(from, to address.Address, value, nonce pack.U256, payload pack.Bytes) (account.Tx, error) {
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
			Version:  types.MessageVersion,
			From:     filfrom,
			To:       filto,
			Value:    big.Int{Int: value.Int()},
			Nonce:    value.Int().Uint64(),
			GasPrice: big.Int{Int: txBuilder.gasPrice.Int()},
			GasLimit: txBuilder.gasLimit.Int().Int64(),
			Method:   methodNum,
			Params:   payload,
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

type Client struct {
	opts   ClientOptions
	node   api.FullNode
	closer jsonrpc.ClientCloser
}

func NewClient(opts ClientOptions) (*Client, error) {
	authToken, err := hex.Decode(opts.AuthToken)
	if err != nil {
		return nil, err
	}
	apiInfo := cli.APIInfo{
		Address: multiaddr.NewMultiaddr(opts.Address),
		Token:   authToken,
	}
	apiAddr, err := apiInfo.DialArgs()
	if err != nil {
		return nil, err
	}
	header := apiInfo.AuthHeader()

	node, closer, err := filecoinClient.NewFullNodeRPC(apiAddr, header)
	if err != nil {
		return nil, err
	}

	return &Client{
		opts:   opts,
		node:   node,
		closer: closer,
	}, nil
}

// Tx returns the transaction uniquely identified by the given transaction
// hash. It also returns the number of confirmations for the transaction.
func (client *Client) Tx(context.Context, pack.Bytes) (account.Tx, pack.U64, error) {
	panic("unimplemented")
}

// SubmitTx to the underlying blockchain network.
func (client *Client) SubmitTx(context.Context, account.Tx) error {
	panic("unimplemented")
}

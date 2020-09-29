package iotex

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"sync"

	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"

	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/gas"
	"github.com/renproject/pack"

	"github.com/btcsuite/btcd/btcec"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
)

var (
	_ account.Client = (*Client)(nil)
	_ gas.Estimator  = (*Client)(nil)
)

type ClientOptions struct {
	Endpoint string
	Secure   bool
}

type Client struct {
	sync.RWMutex
	opts     ClientOptions
	grpcConn *grpc.ClientConn
	client   iotexapi.APIServiceClient
}

func NewClient(opts ClientOptions) *Client {
	return &Client{opts: opts}
}

func (c *Client) EstimateGasPrice(ctx context.Context) (pack.U256, error) {
	if err := c.connect(); err != nil {
		return pack.NewU256FromU64(0), err
	}
	response, err := c.client.SuggestGasPrice(ctx, &iotexapi.SuggestGasPriceRequest{})
	return pack.NewU256FromU64(pack.NewU64(response.GetGasPrice())), err
}

func (c *Client) EstimateGasLimit(_ gas.TxType) (pack.U256, error) {
	if err := c.connect(); err != nil {
		return pack.NewU256FromU64(0), err
	}
	// need a valid pubkey to estimate, just use one
	rawPrivKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return pack.NewU256FromU64(0), err
	}
	rawPubKey := rawPrivKey.PubKey()
	act := &iotextypes.Action{
		SenderPubKey: rawPubKey.SerializeUncompressed(),
	}
	act.Core = &iotextypes.ActionCore{
		Action: &iotextypes.ActionCore_Transfer{
			Transfer: &iotextypes.Transfer{},
		},
	}
	response, err := c.client.EstimateGasForAction(context.Background(), &iotexapi.EstimateGasForActionRequest{Action: act})
	return pack.NewU256FromU64(pack.NewU64(response.GetGas())), err
}

func (c *Client) Nonce(ctx context.Context, addr address.Address) (pack.U256, error) {
	if err := c.connect(); err != nil {
		return pack.NewU256FromU64(0), err
	}
	request := &iotexapi.GetAccountRequest{Address: string(addr)}
	res, err := c.client.GetAccount(ctx, request)
	if err != nil {
		return pack.NewU256FromU64(0), err
	}
	return pack.NewU256FromU64(pack.NewU64(res.GetAccountMeta().GetNonce())), nil
}

func (c *Client) Tx(ctx context.Context, h pack.Bytes) (account.Tx, pack.U64, error) {
	if err := c.connect(); err != nil {
		return nil, 0, err
	}
	res, err := c.client.GetActions(ctx, &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByHash{ByHash: &iotexapi.GetActionByHashRequest{ActionHash: hex.EncodeToString(h)}}})
	if err != nil {
		return nil, 0, err
	}
	if len(res.ActionInfo) != 1 {
		return nil, 0, errors.New("action number should be one")
	}
	ser, err := proto.Marshal(res.ActionInfo[0].GetAction())
	if err != nil {
		return nil, 0, err
	}
	tx := Tx{}
	atx, err := tx.Deserialize(ser)
	return atx, 1, err
}

func (c *Client) SubmitTx(ctx context.Context, t account.Tx) error {
	if err := c.connect(); err != nil {
		return err
	}

	iotexTx, ok := t.(*Tx)
	if !ok {
		return errors.New("not iotex tx")
	}

	_, err := c.client.SendAction(ctx, &iotexapi.SendActionRequest{Action: iotexTx.ToIoTeXTransfer()})
	return err
}

func (c *Client) connect() (err error) {
	c.Lock()
	defer c.Unlock()
	// Check if the existing connection is good.
	if c.grpcConn != nil && c.grpcConn.GetState() != connectivity.Shutdown {
		return
	}
	opts := []grpc.DialOption{}
	if c.opts.Secure {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	c.grpcConn, err = grpc.Dial(c.opts.Endpoint, opts...)
	c.client = iotexapi.NewAPIServiceClient(c.grpcConn)
	return err
}

package waves

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/renproject/pack"
	g "github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/node/grpc"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"google.golang.org/grpc"
)

// Provide communication with node through grpc.
type ClientImpl struct {
	addr    string
	chainID proto.Scheme
}

func NewClient(addr string, chainID proto.Scheme) *ClientImpl {
	return &ClientImpl{
		addr:    addr,
		chainID: chainID,
	}
}

// Retrieve tx by id.
func (a ClientImpl) Tx(ctx context.Context, id pack.Bytes) (Tx, pack.U64, error) {
	conn, err := grpc.Dial(a.addr, grpc.WithInsecure())
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = conn.Close()
	}()
	txReq, err := g.NewTransactionsApiClient(conn).GetTransactions(ctx, &g.TransactionsRequest{
		TransactionIds: [][]byte{id},
	})
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to request transaction")
	}
	txResp, err := txReq.Recv()
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get transaction")
	}
	transfer := txResp.Transaction.Transaction.GetTransfer()
	if transfer == nil {
		return nil, 0, errors.Errorf(
			"expected transaction to be '*Transaction_Transfer', got '%T'",
			txResp.Transaction.Transaction)
	}
	heightReq, err := g.NewBlocksApiClient(conn).GetCurrentHeight(ctx, &empty.Empty{})
	if err != nil {
		return nil, 0, err
	}
	height := heightReq.GetValue()
	if height == 0 {
		return nil, 0, errors.New("invalid value '0' for current blockchain height")
	}
	if txResp.Height > int64(height) {
		return nil, 0, errors.New("height changed during requests")
	}
	var c proto.ProtobufConverter
	tx, err := c.SignedTransaction(txResp.Transaction)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to build transaction from protobuf")
	}
	casted, ok := tx.(*proto.TransferWithProofs)
	if !ok {
		return nil, 0, errors.Errorf("expected transaction to be '*proto.TransferWithProofs', got %T", tx)
	}
	out, err := newTx(casted, a.chainID)
	if err != nil {
		return nil, 0, err
	}
	return out, pack.NewU64(uint64(height) - uint64(txResp.Height)), nil
}

// Send transaction to node.
func (a ClientImpl) SubmitTx(ctx context.Context, tx Tx) error {
	conn, err := grpc.Dial(a.addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	original, ok := tx.(OriginalTx)
	if !ok {
		return errors.New("expected `tx` to be instance of OriginalTx")
	}
	transfer := original.OriginalTx()
	pb, err := transfer.ToProtobufSigned(a.chainID)
	if err != nil {
		return err
	}

	_, err = g.NewTransactionsApiClient(conn).Broadcast(ctx, pb)
	if err != nil {
		return errors.Wrap(err, "failed to broadcast transaction")
	}
	return err
}

package icon

import (
	"context"
	"time"

	"github.com/renproject/pack"
)

const (
	// DefaultClientTimeout used by the Client.
	DefaultClientTimeout = time.Minute
	// DefaultClientTimeoutRetry used by the Client.
	DefaultClientTimeoutRetry = time.Second
	// DefaultClientHost used by the Client. This should only be used for local
	// deployments of the multichain.
	DefaultClientHost = "http://0.0.0.0:9000"
	// DefaultBroadcastMode configures the behaviour of a icon client while it
	// interacts with the icon node. Allowed broadcast modes can be async, sync
	// and block. "async" returns immediately after broadcasting, "sync" returns
	// after the transaction has been checked and "block" waits until the
	// transaction is committed to the chain.
	DefaultBroadcastMode = "sync"
)

type Client struct{}

func (client Client) Tx(ctx context.Context, hash pack.Bytes) (Tx, pack.U64, error) {
	var uInt uint64 = 1<<64 - 1
	return Tx{}, pack.U64(uInt), nil
}

func (client Client) SubmitTx(ctx context.Context, tx Tx) error {
	return nil
}

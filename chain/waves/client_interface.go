package waves

import (
	"context"

	"github.com/renproject/pack"
)

type Client interface {
	// Tx returns the transaction uniquely identified by the given transaction
	// hash. It also returns the number of confirmations for the transaction. If
	// the transaction cannot be found before the context is done, or the
	// transaction is invalid, then an error should be returned.
	Tx(ctx context.Context, id pack.Bytes) (Tx, pack.U64, error)

	// SubmitTx to the underlying chain. If the transaction cannot be found
	// before the context is done, or the transaction is invalid, then an error
	// should be returned.
	SubmitTx(context.Context, Tx) error
}

// Package account defines the Account API. All chains that use an account-based
// model should implement this API. The Account API is used to send and confirm
// transactions between addresses.
package cosmos

import (
	"context"

	"github.com/renproject/multichain/api/account"
	"github.com/tendermint/tendermint/libs/bytes"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type (
	// Tx re-exports account.Tx
	Tx = account.Tx

	// TxBuilder re-exports account.TxBuilder
	TxBuilder = account.TxBuilder
)

// The Client interface defines the functionality required to interact with a
// chain over RPC.
type CompositeClient interface {
	// CompositeClient augments account.Client
	account.Client

	// ABCIQuery to the underlying chain (to get an account's number and nonce, for example).
	ABCIQuery(ctx context.Context, path string, data bytes.HexBytes, height int64, prove bool) (*ctypes.ResultABCIQuery, error)
}

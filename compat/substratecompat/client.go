package substratecompat

import (
	"context"

	"github.com/renproject/pack"
	"github.com/renproject/multichain"
)

type Client interface {
	BurnEvent(ctx context.Context, asset multichain.Asset, nonce pack.Bytes32) (amount pack.U256, to pack.String, confs int64, err error)
}

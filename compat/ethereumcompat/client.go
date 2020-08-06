package ethereumcompat

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/renproject/pack"
)

type Client interface {
	ContractCall(ctx context.Context, contract pack.String, input pack.Bytes) (pack.Bytes, error)
}

type client struct {
	ethclient.Client
}

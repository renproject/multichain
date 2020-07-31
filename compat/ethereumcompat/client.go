package ethereumcompat

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/renproject/pack"
)

type Client interface {
	ContractCall(ctx context.Context, contract pack.String, input pack.Value, outputType pack.Type) (out pack.Value, err error)
}

type client struct {
	ethclient.Client
}

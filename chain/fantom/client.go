package fantom

import (
	"github.com/renproject/multichain/chain/evm"
)

const (
	// DefaultClientRPCURL is the RPC URL used by default, to interact with the
	// fantom node.
	DefaultClientRPCURL = "http://127.0.0.1:18545/"
)

// Client re-exports evm.Client.
type Client = evm.Client

// NewClient re-exports evm.NewClient.
var NewClient = evm.NewClient

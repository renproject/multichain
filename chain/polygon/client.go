package polygon

import (
	"github.com/renproject/multichain/chain/evm"
)

const (
	// DefaultClientRPCURL is the RPC URL used by default, to interact with the
	// polygon node.
	DefaultClientRPCURL = "http://127.0.0.1:28545/"
)

// Client re-exports evm.Client.
type Client = evm.Client

// NewClient re-exports evm.NewClient.
var NewClient = evm.NewClient

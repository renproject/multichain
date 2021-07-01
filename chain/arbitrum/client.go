package arbitrum

import (
	"github.com/renproject/multichain/chain/ethereum"
)

const (
	// DefaultClientRPCURL is the RPC URL used by default, to interact with the
	// arbitrum node.
	DefaultClientRPCURL = "http://127.0.0.1:8547"
)

// Client re-exports ethereum.Client.
type Client = ethereum.Client

// NewClient re-exports ethereum.NewClient.
var NewClient = ethereum.NewClient

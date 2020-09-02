package waves

import (
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

type Tx interface {
	Hash() pack.Bytes
	From() PublicKey
	To() address.Address
	Value() pack.U256
	Nonce() pack.U256
	Payload() contract.CallData
	Sign(privateKey pack.Bytes) error
	Serialize() (pack.Bytes, error)
}

// Way to get waves transaction from Tx interface.
type OriginalTx interface {
	OriginalTx() *proto.TransferWithProofs
}

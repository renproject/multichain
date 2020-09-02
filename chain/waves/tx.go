package waves

import (
	"github.com/pkg/errors"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

// Implements Tx interface in tx_interface.go.
type TxImpl struct {
	tx      *proto.TransferWithProofs
	chainID proto.Scheme
}

func (a TxImpl) OriginalTx() *proto.TransferWithProofs {
	return a.tx
}

func newTx(t *proto.TransferWithProofs, chainID proto.Scheme) (*TxImpl, error) {
	if t.Recipient.Address == nil {
		if t.Recipient.Alias != nil {
			return nil, errors.New("unsupported transaction with alias")
		}
		return nil, errors.New("unsupported transaction with empty recipient")
	}
	return &TxImpl{
		tx:      t,
		chainID: chainID,
	}, nil
}

func (a TxImpl) Hash() pack.Bytes {
	id, _ := a.tx.GetID(a.chainID)
	return id
}

func (a TxImpl) From() PublicKey {
	return a.tx.SenderPK.String()
}

func (a TxImpl) To() address.Address {
	return address.Address(a.tx.Recipient.Address.String())
}

func (a TxImpl) Value() pack.U256 {
	return pack.NewU256FromU64(pack.NewU64(a.tx.Amount))
}

func (a TxImpl) Nonce() pack.U256 {
	return pack.NewU256FromU64(pack.NewU64(a.tx.Timestamp))
}

// Looks like this method is useless.
func (a TxImpl) Payload() contract.CallData {
	return contract.CallData{}
}

func (a TxImpl) Sign(privateKey pack.Bytes) error {
	secret, err := crypto.NewSecretKeyFromBytes(privateKey)
	if err != nil {
		return err
	}
	return a.tx.Sign(a.chainID, secret)
}

func (a TxImpl) Serialize() (pack.Bytes, error) {
	return a.tx.MarshalSignedToProtobuf(a.chainID)
}

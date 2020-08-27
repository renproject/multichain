package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"
)

type TxBuilder struct {
	config *params.ChainConfig
}

func NewTxBuilder(config *params.ChainConfig) TxBuilder {
	return TxBuilder{config: config}
}

func (txBuilder TxBuilder) BuildTx(
	from, to address.Address,
	value, nonce pack.U256,
	payload pack.Bytes,
) (account.Tx, error) {
	panic("unimplemented")
}

type Tx struct {
	from Address

	tx *types.Transaction

	signer types.Signer
	signed bool
}

func (tx *Tx) Hash() pack.Bytes {
	return pack.NewBytes(tx.tx.Hash().Bytes())
}

func (tx *Tx) From() address.Address {
	return address.Address(tx.from.String())
}

func (tx *Tx) To() address.Address {
	return address.Address(tx.tx.To().String())
}

func (tx *Tx) Value() pack.U256 {
	return pack.NewU256FromInt(tx.tx.Value())
}

func (tx *Tx) Nonce() pack.U256 {
	return pack.NewU256FromU64(pack.NewU64(tx.tx.Nonce()))
}

func (tx *Tx) Payload() contract.CallData {
	return contract.CallData(pack.NewBytes(tx.tx.Data()))
}

func (tx *Tx) Sighashes() (pack.Bytes32, error) {
	sighash := tx.signer.Hash(tx.tx)
	return pack.NewBytes32(sighash), nil
}

func (tx *Tx) Sign(signature pack.Bytes65, pubKey pack.Bytes) error {
	if tx.signed {
		return fmt.Errorf("already signed")
	}

	signedTx, err := tx.tx.WithSignature(tx.signer, signature.Bytes())
	if err != nil {
		return err
	}

	tx.tx = signedTx
	tx.signed = true
	return nil
}

func (tx *Tx) Serialize() (pack.Bytes, error) {
	serialized, err := tx.tx.MarshalJSON()
	if err != nil {
		return pack.Bytes{}, err
	}

	return pack.NewBytes(serialized), nil
}

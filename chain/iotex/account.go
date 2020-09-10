package iotex

import (
	"errors"
	"math/big"

	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/go-pkgs/hash"
	iotexaddr "github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"

	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"

	"github.com/golang/protobuf/proto"
)

type Tx struct {
	from               address.Address
	to                 address.Address
	value, nonce       pack.U256
	gasLimit, gasPrice pack.U256
	payload            pack.Bytes
	sig                pack.Bytes65
	publicKey          pack.Bytes
}

func (t *Tx) ToIoTeXTransfer() *iotextypes.Action {
	return &iotextypes.Action{
		Core: &iotextypes.ActionCore{
			GasLimit: t.gasLimit.Int().Uint64(),
			GasPrice: t.gasPrice.Int().String(),
			Nonce:    t.nonce.Int().Uint64(),
			Action: &iotextypes.ActionCore_Transfer{
				Transfer: &iotextypes.Transfer{
					Amount:    t.value.String(),
					Recipient: string(t.to),
					Payload:   t.payload,
				},
			},
		},
		SenderPubKey: t.publicKey,
		Signature:    t.sig[:],
	}
}

func (t *Tx) Hash() pack.Bytes {
	sealed, err := t.Serialize()
	if err != nil {
		return nil
	}
	h := hash.Hash256b(sealed)
	return h[:]
}

func (t *Tx) From() address.Address { return t.from }

func (t *Tx) To() address.Address { return t.to }

func (t *Tx) Value() pack.U256 { return t.value }

func (t *Tx) Nonce() pack.U256 { return t.nonce }

func (t *Tx) Payload() contract.CallData { return contract.CallData(t.payload) }

func (t *Tx) PublicKey() pack.Bytes { return t.publicKey }

func (t *Tx) Sighashes() ([]pack.Bytes32, error) {
	act := t.ToIoTeXTransfer()
	core, err := proto.Marshal(act.GetCore())
	if err != nil {
		return nil, err
	}
	h := hash.Hash256b(core)
	return []pack.Bytes32{pack.Bytes32(h)}, nil
}

func (t *Tx) Sign(sig []pack.Bytes65, publicKey pack.Bytes) error {
	copy(t.sig[:], sig[0][:])

	pub, err := crypto.BytesToPublicKey(publicKey)
	if err != nil {
		return err
	}
	pubBytes := pub.Bytes()
	t.publicKey = make([]byte, len(pubBytes))
	copy(t.publicKey[:], pubBytes[:])
	return nil
}

func (t *Tx) Serialize() (pack.Bytes, error) {
	return proto.Marshal(t.ToIoTeXTransfer())
}

func (t *Tx) Deserialize(ser pack.Bytes) (account.Tx, error) {
	act := &iotextypes.Action{}
	if err := proto.Unmarshal(ser, act); err != nil {
		return nil, err
	}
	pub, err := crypto.BytesToPublicKey(act.GetSenderPubKey())
	if err != nil {
		return nil, err
	}
	from, err := iotexaddr.FromBytes(pub.Hash())
	if err != nil {
		return nil, err
	}
	sig := pack.Bytes65{}
	copy(sig[:], act.GetSignature())
	amount, ok := new(big.Int).SetString(act.GetCore().GetTransfer().GetAmount(), 10)
	if !ok {
		return nil, errors.New("amount convert error")
	}
	gasPrice, ok := new(big.Int).SetString(act.GetCore().GetGasPrice(), 10)
	if !ok {
		return nil, errors.New("gas price convert error")
	}
	return &Tx{
		from:      address.Address(from.String()),
		to:        address.Address(act.GetCore().GetTransfer().GetRecipient()),
		value:     pack.NewU256FromInt(amount),
		nonce:     pack.NewU256FromU64(pack.U64(act.GetCore().GetNonce())),
		gasLimit:  pack.NewU256FromU64(pack.U64(act.GetCore().GetGasLimit())),
		gasPrice:  pack.NewU256FromInt(gasPrice),
		payload:   act.GetCore().GetTransfer().GetPayload(),
		sig:       sig,
		publicKey: act.GetSenderPubKey(),
	}, nil
}

type TxBuilder struct{}

func (t *TxBuilder) BuildTx(from, to address.Address, value, nonce pack.U256, gasPrice, gasLimit pack.U256, payload pack.Bytes) (account.Tx, error) {
	return &Tx{
		from:     from,
		to:       to,
		value:    value,
		nonce:    nonce,
		gasLimit: gasLimit,
		gasPrice: gasPrice,
		payload:  payload,
		sig:      pack.NewBytes65([65]byte{0}),
	}, nil
}

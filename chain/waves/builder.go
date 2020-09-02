package waves

import (
	"time"

	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

// Waves requires public key in from field.
type PublicKey = string

// Helper for creating transactions.
type TxBuilder struct {
	chainID byte
}

func NewTxBuilder(chainID byte) *TxBuilder {
	return &TxBuilder{
		chainID: chainID,
	}
}

// BuildTx accepts public key first argument. All other is the same.
func (a TxBuilder) BuildTx(from PublicKey, to address.Address, value, nonce pack.U256, payload pack.Bytes) (Tx, error) {
	// Now is the same as nonce, just unique value to distinguish transactions from each other.
	now := proto.NewTimestampFromTime(time.Now())
	amount := value.Int().Uint64()
	rec, err := proto.NewRecipientFromString(string(to))
	if err != nil {
		return nil, err
	}
	pub, err := crypto.NewPublicKeyFromBase58(from)
	if err != nil {
		return nil, err
	}
	tx := proto.NewUnsignedTransferWithProofs(
		3, // last version
		pub,
		// if empty struct default asset will be used (Waves). Is need some other asset,
		// proto.NewOptionalAssetFromDigest(), proto.NewOptionalAssetFromBytes(), proto.NewOptionalAssetFromString() can be used.
		// It accepts asset unique id (transaction id), https://wavesexplorer.com/tx/4oZLvmLC1jZQXRSw97c2k4VfK6HgAFXXjWvHFTNH4j2t as example.
		proto.OptionalAsset{},
		proto.OptionalAsset{},
		now,
		amount,
		100000, //minimal fee
		rec,
		nil,
	)
	return newTx(tx, a.chainID)
}

package cosmos

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/renproject/multichain/compat/cosmoscompat"
	"github.com/renproject/pack"

	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// NewClient returns returns a new Client with default codec
func NewClient(opts cosmoscompat.ClientOptions) cosmoscompat.Client {
	return cosmoscompat.NewClient(opts, simapp.MakeCodec())
}

type txBuilder struct {
	auth.TxBuilder
	cdc *codec.Codec
}

// NewTxBuilder returns an implementation of the transaction builder interface
// from the Cosmos Compat API, and exposes the functionality to build simple
// Cosmos based transactions.
func NewTxBuilder(options cosmoscompat.TxOptions) cosmoscompat.TxBuilder {
	cdc := simapp.MakeCodec()

	return txBuilder{
		TxBuilder: auth.NewTxBuilder(
			utils.GetTxEncoder(cdc),
			options.AccountNumber.Uint64(),
			options.SequenceNumber.Uint64(),
			options.Gas.Uint64(),
			0,
			false,
			options.ChainID.String(),
			options.Memo.String(),
			options.Fees.Coins(), sdk.DecCoins{},
		),
		cdc: simapp.MakeCodec(),
	}
}

// WithCodec replace codec with custom one
func (builder txBuilder) WithCodec(cdc *codec.Codec) cosmoscompat.TxBuilder {
	builder.WithTxEncoder(utils.GetTxEncoder(cdc))
	builder.cdc = cdc
	return builder
}

func (builder txBuilder) BuildTx(sendMsgs []cosmoscompat.MsgSend) (cosmoscompat.Tx, error) {
	sdkMsgs := []sdk.Msg{}
	for _, sendMsg := range sendMsgs {
		sdkMsgs = append(sdkMsgs, sendMsg.Msg())
	}

	signMsg, err := builder.BuildSignMsg(sdkMsgs)
	if err != nil {
		return nil, err
	}

	return &Tx{cdc: builder.cdc, signMsg: signMsg}, nil
}

// Tx represents a simple Terra transaction that implements the Cosmos Compat
// API.
type Tx struct {
	cdc        *codec.Codec
	signMsg    auth.StdSignMsg
	signatures []auth.StdSignature
}

// Hash return txhash bytes
func (tx *Tx) Hash() (pack.Bytes32, error) {
	if len(tx.signatures) == 0 {
		return pack.Bytes32{}, fmt.Errorf("please do tx.Sign() first to get a hash")
	}

	txBytes, err := tx.Serialize()
	if err != nil {
		return pack.Bytes32{}, err
	}

	hashBytes := pack.Bytes32{}
	hashBytes.Unmarshal(tmhash.Sum(txBytes), 32)
	return hashBytes, nil
}

// SigBytes that need to be signed before this transaction can be
// submitted.
func (tx *Tx) SigBytes() pack.Bytes {
	return tx.signMsg.Bytes()
}

// Sign the transaction by injecting signatures and the serialized pubkey of
// the signer.
func (tx *Tx) Sign(signatures []cosmoscompat.StdSignature) error {
	var stdSignatures []auth.StdSignature
	for _, sig := range signatures {
		var pubKey secp256k1.PubKeySecp256k1
		copy(pubKey[:], sig.PubKey[:secp256k1.PubKeySecp256k1Size])

		stdSignatures = append(stdSignatures, auth.StdSignature{
			Signature: sig.Signature,
			PubKey:    pubKey,
		})
	}

	signers := make(map[string]bool)
	for _, msg := range tx.signMsg.Msgs {
		for _, signer := range msg.GetSigners() {
			fmt.Println("SIBONG", signer.String())
			signers[signer.String()] = true
		}
	}

	for _, sig := range stdSignatures {
		signer := sdk.AccAddress(sig.Address()).String()
		if _, ok := signers[signer]; !ok {
			return fmt.Errorf("wrong signer: %s", signer)
		}
	}

	if len(signers) != len(stdSignatures) {
		return fmt.Errorf("insufficient signers")
	}

	fmt.Println("SIBONG", stdSignatures)
	tx.signatures = stdSignatures
	return nil
}

// Serialize the transaction.
func (tx *Tx) Serialize() (pack.Bytes, error) {
	txBytes, err := tx.cdc.MarshalBinaryLengthPrefixed(auth.NewStdTx(tx.signMsg.Msgs, tx.signMsg.Fee, tx.signatures, tx.signMsg.Memo))
	if err != nil {
		return pack.Bytes{}, err
	}

	return txBytes, nil
}

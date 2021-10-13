package cosmos

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/types"
	txTypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	bankType "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/renproject/id"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"
	"github.com/renproject/surge"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

const (
	// DefaultChainID used by the Client.
	DefaultChainID = pack.String("testnet")
	// DefaultSignMode used in signing the tx
	DefaultSignMode = 1
)

// TxBuilderOptions only contains necessary options to build tx from tx builder
type TxBuilderOptions struct {
	ChainID pack.String
}

// DefaultTxBuilderOptions returns TxBuilderOptions with the default settings.
func DefaultTxBuilderOptions() TxBuilderOptions {
	return TxBuilderOptions{
		ChainID: DefaultChainID,
	}
}

// WithChainID sets the chain ID used by the transaction builder.
func (opts TxBuilderOptions) WithChainID(chainID pack.String) TxBuilderOptions {
	opts.ChainID = chainID
	return opts
}

type txBuilder struct {
	client   *Client
	chainID  pack.String
	signMode int32
}

// NewTxBuilder returns an implementation of the transaction builder interface
// from the Cosmos Compat API, and exposes the functionality to build simple
// Cosmos based transactions.
func NewTxBuilder(options TxBuilderOptions, client *Client) account.TxBuilder {
	return txBuilder{
		signMode: DefaultSignMode,
		client:   client,
		chainID:  options.ChainID,
	}
}

// WithSignMode ad custom sign mode to the txBuilder
func (builder txBuilder) WithSignMode(signMode int32) txBuilder {
	builder.signMode = signMode
	return builder
}

// BuildTx consumes a list of MsgSend to build and return a cosmos transaction.
// This transaction is unsigned, and must be signed before submitting to the
// cosmos chain.
func (builder txBuilder) BuildTx(ctx context.Context, fromPubKey *id.PubKey, to address.Address, value, nonce, gasLimit, gasPrice, gasCap pack.U256, payload pack.Bytes) (account.Tx, error) {
	pubKeyBytes, err := surge.ToBinary(fromPubKey)
	if err != nil {
		return nil, err
	}
	pubKey := secp256k1.PubKey{Key: pubKeyBytes}
	from := multichain.Address(types.AccAddress(pubKey.Address()).String())

	fromAddr, err := types.AccAddressFromBech32(string(from))
	if err != nil {
		return nil, err
	}

	toAddr, err := types.AccAddressFromBech32(string(to))
	if err != nil {
		return nil, err
	}

	sendMsg := MsgSend{
		FromAddress: Address(fromAddr),
		ToAddress:   Address(toAddr),
		Amount: []Coin{
			{
				Denom:  builder.client.opts.CoinDenom,
				Amount: pack.NewU64(value.Int().Uint64()),
			},
		},
	}

	fees := Coins{Coin{
		Denom:  builder.client.opts.CoinDenom,
		Amount: pack.NewU64(gasPrice.Mul(gasLimit).Int().Uint64()),
	}}

	accountNumber, err := builder.client.AccountNumber(ctx, from)
	if err != nil {
		return nil, err
	}

	txBuilder := builder.client.ctx.TxConfig.NewTxBuilder()
	txBuilder.SetFeeAmount(fees.Coins())
	txBuilder.SetGasLimit(gasLimit.Int().Uint64())
	txBuilder.SetMemo(string(payload))

	err = txBuilder.SetMsgs(sendMsg.Msg())
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	sigData := signing.SingleSignatureData{
		SignMode:  signing.SignMode(builder.signMode),
		Signature: nil,
	}

	sig := signing.SignatureV2{
		PubKey:   &pubKey,
		Data:     &sigData,
		Sequence: nonce.Int().Uint64(),
	}
	if err = txBuilder.SetSignatures(sig); err != nil {
		return nil, err
	}
	signerData := authsigning.SignerData{
		AccountNumber: accountNumber.Uint64(),
		ChainID:       string(builder.chainID),
		Sequence:      nonce.Int().Uint64(),
	}
	txConfig := builder.client.ctx.TxConfig
	signMsg, err := txConfig.SignModeHandler().GetSignBytes(signing.SignMode(builder.signMode), signerData, txBuilder.GetTx())
	if err != nil {
		return nil, err
	}
	return &Tx{
		encoder:   builder.client.ctx.TxConfig.TxEncoder(),
		signMsg:   signMsg,
		sigV2:     sig,
		txBuilder: txBuilder,
		sendMsg:   &sendMsg,
		memo:      string(payload),
	}, nil
}

// Coin copy type from types.coin
type Coin struct {
	Denom  pack.String `json:"denom"`
	Amount pack.U64    `json:"amount"`
}

// Coins array of Coin
type Coins []Coin

// Coins parse pack coins to sdk coins
func (coins Coins) Coins() types.Coins {
	sdkCoins := make(types.Coins, 0, len(coins))
	for _, coin := range coins {
		sdkCoins = append(sdkCoins, types.Coin{
			Denom:  coin.Denom.String(),
			Amount: types.NewInt(int64(coin.Amount.Uint64())),
		})
	}

	sdkCoins.Sort()
	return sdkCoins
}

// MsgSend - high level transaction of the coin module
type MsgSend struct {
	FromAddress Address `json:"from_address" yaml:"from_address"`
	ToAddress   Address `json:"to_address" yaml:"to_address"`
	Amount      Coins   `json:"amount" yaml:"amount"`
}

// Msg convert MsgSend to types.Msg
func (msg MsgSend) Msg() types.Msg {
	return bankType.NewMsgSend(
		msg.FromAddress.AccAddress(),
		msg.ToAddress.AccAddress(),
		msg.Amount.Coins(),
	)
}

// Tx is a tx.Tx wrapper
type Tx struct {
	originalTx *txTypes.Tx
	encoder    types.TxEncoder
	sendMsg    *MsgSend
	memo       string
	signMsg    []byte
	sigV2      signing.SignatureV2
	txBuilder  client.TxBuilder
}

// From returns the sender of the transaction
func (t Tx) From() address.Address {
	if t.originalTx != nil {
		return address.Address(t.originalTx.GetBody().Messages[0].GetCachedValue().(*bankType.MsgSend).FromAddress)
	}

	if t.sendMsg != nil {
		return address.Address(t.sendMsg.FromAddress.String())
	}
	return address.Address("")
}

// To returns the recipients of the transaction. For the cosmos chain, there
// can be multiple recipients from a single transaction.
func (t Tx) To() address.Address {
	if t.originalTx != nil {
		return address.Address(t.originalTx.GetBody().Messages[0].GetCachedValue().(*bankType.MsgSend).ToAddress)
	}

	if t.sendMsg != nil {
		return address.Address(t.sendMsg.ToAddress.String())
	}
	return address.Address("")
}

// Value returns the values being transferred in a transaction. For the cosmos
// chain, there can be multiple messages (each with a different value being
// transferred) in a single transaction.
func (t Tx) Value() pack.U256 {
	value := pack.NewU64(0)
	if t.originalTx != nil {
		msgs := t.originalTx.GetBody().Messages
		for _, msg := range msgs {
			value.AddAssign(pack.NewU64(msg.GetCachedValue().(*bankType.MsgSend).Amount[0].Amount.Uint64()))
		}
	} else if t.sendMsg != nil {
		value.AddAssign(pack.NewU64(t.sendMsg.Amount.Coins()[0].Amount.Uint64()))
	}
	return pack.NewU256FromU64(value)
}

// Nonce returns the transaction count of the transaction sender.
func (t Tx) Nonce() pack.U256 {
	if t.originalTx != nil {
		return pack.NewU256FromUint64(t.originalTx.GetAuthInfo().SignerInfos[0].Sequence)
	}

	if t.sendMsg != nil {
		return pack.NewU256FromU64(pack.NewU64(t.sigV2.Sequence))
	}

	return pack.NewU256FromU64(0)
}

// Payload returns the memo attached to the transaction.
func (t Tx) Payload() contract.CallData {
	if t.originalTx != nil {
		return contract.CallData(t.originalTx.GetBody().Memo)
	}

	if t.sendMsg != nil {
		return contract.CallData(t.memo)
	}
	return contract.CallData("")
}

// Hash return txhash bytes.
func (t Tx) Hash() pack.Bytes {
	txBytes, err := t.Serialize()
	if err != nil {
		return pack.Bytes{}
	}

	return pack.NewBytes(tmhash.Sum(txBytes))
}

// Sighashes that need to be signed before this transaction can be submitted.
func (t Tx) Sighashes() ([]pack.Bytes32, error) {
	return []pack.Bytes32{sha256.Sum256(t.signMsg)}, nil
}

// Sign the transaction by injecting signatures and the serialized pubkey of
// the signer.
func (t *Tx) Sign(signatures []pack.Bytes65, pubKey pack.Bytes) error {
	if len(signatures) == 0 {
		return fmt.Errorf("zero signatures found")
	}
	sig := serializeSig(signatureFromBytes(signatures[0].Bytes()))
	singleData := t.sigV2.Data.(*signing.SingleSignatureData)
	singleData.Signature = sig
	t.sigV2.Data = singleData
	err := t.txBuilder.SetSignatures(t.sigV2)
	if err != nil {
		return err
	}
	return nil
}

// Serialize the transaction.
func (t Tx) Serialize() (pack.Bytes, error) {
	var txBytes []byte
	var err error = nil
	if t.originalTx != nil {
		txBytes, err = t.encoder(tx.WrapTx(t.originalTx).GetTx())
	} else if t.sendMsg != nil {
		txBytes, err = t.encoder(t.txBuilder.GetTx())
	}
	if err != nil {
		return pack.Bytes{}, err
	}

	return txBytes, nil
}

func signatureFromBytes(sigStr []byte) *btcec.Signature {
	return &btcec.Signature{
		R: new(big.Int).SetBytes(sigStr[:32]),
		S: new(big.Int).SetBytes(sigStr[32:64]),
	}
}

// Serialize signature to R || S.
// R, S are padded to 32 bytes respectively.
func serializeSig(sig *btcec.Signature) []byte {
	rBytes := sig.R.Bytes()
	sBytes := sig.S.Bytes()
	sigBytes := make([]byte, 64)
	// 0 pad the byte arrays from the left if they aren't big enough.
	copy(sigBytes[32-len(rBytes):32], rBytes)
	copy(sigBytes[64-len(sBytes):64], sBytes)
	return sigBytes
}

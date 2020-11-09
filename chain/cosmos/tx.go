package cosmos

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

const (
	// DefaultChainID used by the Client.
	DefaultChainID = pack.String("testnet")
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
	client  *Client
	chainID pack.String
}

// NewTxBuilder returns an implementation of the transaction builder interface
// from the Cosmos Compat API, and exposes the functionality to build simple
// Cosmos based transactions.
func NewTxBuilder(options TxBuilderOptions, client *Client) account.TxBuilder {
	if client.cdc == nil {
		client.cdc = simapp.MakeCodec()
	}

	return txBuilder{
		client:  client,
		chainID: options.ChainID,
	}
}

// BuildTx consumes a list of MsgSend to build and return a cosmos transaction.
// This transaction is unsigned, and must be signed before submitting to the
// cosmos chain.
func (builder txBuilder) BuildTx(ctx context.Context, from, to address.Address, value, nonce, gasLimit, gasPrice, gasCap pack.U256, payload pack.Bytes) (account.Tx, error) {
	types.GetConfig().SetBech32PrefixForAccount(builder.client.hrp, builder.client.hrp+"pub")

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

	txBuilder := auth.NewTxBuilder(
		utils.GetTxEncoder(builder.client.cdc),
		accountNumber.Uint64(),
		nonce.Int().Uint64(),
		gasLimit.Int().Uint64(),
		0,
		false,
		builder.chainID.String(),
		string(payload),
		fees.Coins(),
		types.DecCoins{},
	)

	sdkMsgs := []types.Msg{sendMsg.Msg()}

	signMsg, err := txBuilder.BuildSignMsg(sdkMsgs)
	if err != nil {
		return nil, err
	}

	return &StdTx{
		msgs:    []MsgSend{sendMsg},
		fee:     parseStdFee(signMsg.Fee),
		memo:    pack.String(payload.String()),
		cdc:     builder.client.cdc,
		signMsg: signMsg,
	}, nil
}

// Coin copy type from types.coin
type Coin struct {
	Denom  pack.String `json:"denom"`
	Amount pack.U64    `json:"amount"`
}

// parseCoin parse types.Coin to Coin
func parseCoin(sdkCoin types.Coin) Coin {
	return Coin{
		Denom:  pack.NewString(sdkCoin.Denom),
		Amount: pack.U64(uint64(sdkCoin.Amount.Int64())),
	}
}

// Coins array of Coin
type Coins []Coin

// parseCoins parse types.Coins to Coins
func parseCoins(sdkCoins types.Coins) Coins {
	coins := make(Coins, 0, len(sdkCoins))
	for _, sdkCoin := range sdkCoins {
		coins = append(coins, parseCoin(sdkCoin))
	}
	return coins
}

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
	return bank.NewMsgSend(
		msg.FromAddress.AccAddress(),
		msg.ToAddress.AccAddress(),
		msg.Amount.Coins(),
	)
}

// NOTE: we only support MsgSend
// parseMsg parse types.Msg to MsgSend
func parseMsg(msg types.Msg) (MsgSend, error) {
	if msg, ok := msg.(bank.MsgSend); ok {
		return MsgSend{
			FromAddress: Address(msg.FromAddress),
			ToAddress:   Address(msg.ToAddress),
			Amount:      parseCoins(msg.Amount),
		}, nil
	}

	return MsgSend{}, fmt.Errorf("Failed to parse %v to MsgSend", msg)
}

// StdFee auth.StdFee wrapper
type StdFee struct {
	Amount Coins    `json:"amount" yaml:"amount"`
	Gas    pack.U64 `json:"gas" yaml:"gas"`
}

// parseStdFee parse auth.StdFee to StdFee
func parseStdFee(stdFee auth.StdFee) StdFee {
	return StdFee{
		Amount: parseCoins(stdFee.Amount),
		Gas:    pack.U64(stdFee.Gas),
	}
}

// StdSignature auth.StdStdSignature wrapper
type StdSignature struct {
	PubKey    pack.Bytes `json:"pub_key" yaml:"pub_key"`
	Signature pack.Bytes `json:"signature" yaml:"signature"`
}

// parseStdSignature parse auth.StdSignature to StdSignature
func parseStdSignature(stdSig auth.StdSignature) StdSignature {
	return StdSignature{
		PubKey:    pack.NewBytes(stdSig.PubKey.Bytes()),
		Signature: pack.NewBytes(stdSig.Signature),
	}
}

// StdTx auth.StStdTx wrapper
type StdTx struct {
	msgs       []MsgSend
	fee        StdFee
	memo       pack.String
	signatures []auth.StdSignature

	cdc     *codec.Codec
	signMsg auth.StdSignMsg
}

// From returns the sender of the transaction
func (tx StdTx) From() address.Address {
	return address.Address(tx.msgs[0].FromAddress.AccAddress().String())
}

// To returns the recipients of the transaction. For the cosmos chain, there
// can be multiple recipients from a single transaction.
func (tx StdTx) To() address.Address {
	return address.Address(tx.msgs[0].ToAddress.AccAddress().String())
}

// Value returns the values being transferred in a transaction. For the cosmos
// chain, there can be multiple messages (each with a different value being
// transferred) in a single transaction.
func (tx StdTx) Value() pack.U256 {
	value := pack.NewU64(0)
	for _, msg := range tx.msgs {
		value.AddAssign(msg.Amount[0].Amount)
	}
	return pack.NewU256FromU64(value)
}

// Nonce returns the transaction count of the transaction sender.
func (tx StdTx) Nonce() pack.U256 {
	return pack.NewU256FromU64(pack.NewU64(tx.signMsg.Sequence))
}

// Payload returns the memo attached to the transaction.
func (tx StdTx) Payload() contract.CallData {
	return contract.CallData(pack.NewBytes([]byte(tx.memo)))
}

// Hash return txhash bytes.
func (tx StdTx) Hash() pack.Bytes {
	if len(tx.signatures) == 0 {
		return pack.Bytes{}
	}

	txBytes, err := tx.Serialize()
	if err != nil {
		return pack.Bytes{}
	}

	return pack.NewBytes(tmhash.Sum(txBytes))
}

// Sighashes that need to be signed before this transaction can be submitted.
func (tx StdTx) Sighashes() ([]pack.Bytes32, error) {
	sighashBytes := crypto.Sha256(tx.signMsg.Bytes())
	if len(sighashBytes) != 32 {
		return nil, fmt.Errorf("expected 32 bytes, got %v bytes", len(tx.signMsg.Bytes()))
	}
	sighash := pack.Bytes32{}
	copy(sighash[:], sighashBytes)
	return []pack.Bytes32{sighash}, nil
}

// Sign the transaction by injecting signatures and the serialized pubkey of
// the signer.
func (tx *StdTx) Sign(signatures []pack.Bytes65, pubKey pack.Bytes) error {
	var stdSignatures []auth.StdSignature
	for _, sig := range signatures {
		var cpPubKey secp256k1.PubKeySecp256k1
		copy(cpPubKey[:], pubKey[:secp256k1.PubKeySecp256k1Size])
		stdSignatures = append(stdSignatures, auth.StdSignature{
			// Cosmos uses 64-bytes signature
			// https://github.com/tendermint/tendermint/blob/v0.33.8/crypto/secp256k1/secp256k1_nocgo.go#L60-L70
			Signature: sig[:64],
			PubKey:    cpPubKey,
		})
	}

	signers := make(map[string]bool)
	for _, msg := range tx.signMsg.Msgs {
		for _, signer := range msg.GetSigners() {
			signers[signer.String()] = true
		}
	}

	for _, sig := range stdSignatures {
		signer := types.AccAddress(sig.Address()).String()
		if _, ok := signers[signer]; !ok {
			return fmt.Errorf("wrong signer: %s", signer)
		}
	}

	if len(signers) != len(stdSignatures) {
		return fmt.Errorf("insufficient signers")
	}

	tx.signatures = stdSignatures
	return nil
}

// Serialize the transaction.
func (tx StdTx) Serialize() (pack.Bytes, error) {
	txBytes, err := tx.cdc.MarshalBinaryLengthPrefixed(auth.NewStdTx(tx.signMsg.Msgs, tx.signMsg.Fee, tx.signatures, tx.signMsg.Memo))
	if err != nil {
		return pack.Bytes{}, err
	}

	return txBytes, nil
}

// parseStdTx parse auth.StdTx to StdTx
func parseStdTx(stdTx auth.StdTx) (StdTx, error) {
	var msgs []MsgSend
	for _, msg := range stdTx.Msgs {
		msg, err := parseMsg(msg)
		if err != nil {
			return StdTx{}, err
		}

		msgs = append(msgs, msg)
	}

	fee := parseStdFee(stdTx.Fee)
	memo := pack.NewString(stdTx.Memo)

	return StdTx{
		msgs:       msgs,
		fee:        fee,
		memo:       memo,
		signatures: stdTx.Signatures,
	}, nil
}

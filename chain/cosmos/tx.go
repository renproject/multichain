package cosmos

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/renproject/pack"
)

// TxBuilder defines an interface that can be used to build simple Bitcoin
// transactions.
type TxBuilder interface {
	// BuildTx returns a simple Bitcoin transaction that consumes a set of
	// Bitcoin outputs and uses the funds to make payments to a set of Bitcoin
	// recipients. The sum value of the inputs must be greater than the sum
	// value of the outputs, and the difference is paid as a fee to the Bitcoin
	// network.
	BuildTx(msgs []MsgSend) (Tx, error)
	WithCodec(cdc *codec.Codec) TxBuilder
}

// Tx defines an interface that must be implemented by all types of Bitcoin
// transactions.
type Tx interface {
	// Hash of the transaction.
	Hash() (pack.Bytes32, error)

	// SigBytes that need to be signed before this transaction can be
	// submitted.
	SigBytes() pack.Bytes

	// Sign the transaction by injecting signatures and the serialized pubkey of
	// the signer.
	Sign([]StdSignature) error

	// Serialize the transaction.
	Serialize() (pack.Bytes, error)
}

// An Address is a public address that can be encoded/decoded to/from strings.
// Addresses are usually formatted different between different network
// configurations.
type Address sdk.AccAddress

// AccAddress convert Address to sdk.AccAddress
func (addr Address) AccAddress() sdk.AccAddress {
	return sdk.AccAddress(addr)
}

// TxBuilderOptions only contains necessary options to build tx from tx builder
type TxBuilderOptions struct {
	AccountNumber  pack.U64    `json:"account_number"`
	SequenceNumber pack.U64    `json:"sequence_number"`
	Gas            pack.U64    `json:"gas"`
	ChainID        pack.String `json:"chain_id"`
	Memo           pack.String `json:"memo"`
	Fees           Coins       `json:"fees"`
}

// Coin copy type from sdk.coin
type Coin struct {
	Denom  pack.String `json:"denom"`
	Amount pack.U64    `json:"amount"`
}

// parseCoin parse sdk.Coin to Coin
func parseCoin(sdkCoin sdk.Coin) Coin {
	return Coin{
		Denom:  pack.NewString(sdkCoin.Denom),
		Amount: pack.U64(uint64(sdkCoin.Amount.Int64())),
	}
}

// Coins array of Coin
type Coins []Coin

// parseCoins parse sdk.Coins to Coins
func parseCoins(sdkCoins sdk.Coins) Coins {
	var coins Coins
	for _, sdkCoin := range sdkCoins {
		coins = append(coins, parseCoin(sdkCoin))
	}
	return coins
}

// Coins parse pack coins to sdk coins
func (coins Coins) Coins() sdk.Coins {
	sdkCoins := sdk.Coins{}
	for _, coin := range coins {
		sdkCoin := sdk.Coin{
			Denom:  coin.Denom.String(),
			Amount: sdk.NewInt(int64(coin.Amount.Uint64())),
		}

		sdkCoins = append(sdkCoins, sdkCoin)
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

// Msg convert MsgSend to sdk.Msg
func (msg MsgSend) Msg() sdk.Msg {
	return bank.NewMsgSend(
		msg.FromAddress.AccAddress(),
		msg.ToAddress.AccAddress(),
		msg.Amount.Coins(),
	)
}

// NOTE: we only support MsgSend
// parseMsg parse sdk.Msg to MsgSend
func parseMsg(msg sdk.Msg) (MsgSend, error) {
	if msg, ok := msg.(bank.MsgSend); ok {
		return MsgSend{
			FromAddress: Address(msg.FromAddress),
			ToAddress:   Address(msg.ToAddress),
			Amount:      parseCoins(msg.Amount),
		}, nil
	} else {
		return MsgSend{}, fmt.Errorf("Failed to parse %v to MsgSend", msg)
	}
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
	Msgs       []MsgSend      `json:"msgs" yaml:"msgs"`
	Fee        StdFee         `json:"fee" yaml:"fee"`
	Signatures []StdSignature `json:"signatures" yaml:"signatures"`
	Memo       pack.String    `json:"memo" yaml:"memo"`
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

	var sigs []StdSignature
	for _, sig := range stdTx.Signatures {
		sigs = append(sigs, parseStdSignature(sig))
	}

	fee := parseStdFee(stdTx.Fee)
	memo := pack.NewString(stdTx.Memo)

	return StdTx{
		Msgs:       msgs,
		Fee:        fee,
		Memo:       memo,
		Signatures: sigs,
	}, nil
}

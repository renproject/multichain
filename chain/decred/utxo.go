package decred

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/chaincfg/v3"
	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/decred/dcrd/txscript/v3"
	"github.com/decred/dcrd/wire"
	"github.com/renproject/multichain/api/utxo"
	"github.com/renproject/pack"
)

// The TxBuilder is an implementation of a UTXO-compatible transaction builder
// for Bitcoin.
type TxBuilder struct {
	params *chaincfg.Params
}

// NewTxBuilder returns a transaction builder that builds UTXO-compatible
// Bitcoin transactions for the given chain configuration (this means that it
// can be used for regnet, testnet, and mainnet, but also for networks that are
// minimally modified forks of the Bitcoin network).
func NewTxBuilder(params *chaincfg.Params) TxBuilder {
	return TxBuilder{params: params}
}

// BuildTx returns a Decred transaction that consumes funds from the given
// inputs, and sends them to the given recipients. The difference in the sum
// value of the inputs and the sum value of the recipients is paid as a fee to
// the Bitcoin network. This fee must be calculated independently of this
// function. Outputs produced for recipients will use P2PKH, P2SH, P2WPKH, or
// P2WSH scripts as the pubkey script, based on the format of the recipient
// address.
func (txBuilder TxBuilder) BuildTx(inputs []utxo.Input, recipients []utxo.Recipient) (utxo.Tx, error) {
	msgTx := wire.NewMsgTx()

	// Inputs
	for _, input := range inputs {
		hash := chainhash.Hash{}
		copy(hash[:], input.Hash)
		index := input.Index.Uint32()
		amt, err := dcrutil.NewAmount(1)
		if err != nil {
			return nil, err
		}
		prevOutV := int64(amt)
		msgTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(&hash, index, wire.TxTreeRegular), prevOutV, []byte{}))
	}

	// Outputs
	for _, recipient := range recipients {
		addr, err := dcrutil.DecodeAddress(string(recipient.To), txBuilder.params)
		if err != nil {
			return nil, err
		}
		// Ensure the address is one of the supported types.
		switch addr.(type) {
		case *dcrutil.AddressPubKeyHash:
			fmt.Printf("Address Type: %+v \n", "PubKeyHash")
		case *dcrutil.AddressScriptHash:
			fmt.Printf("Address Type: %+v \n", "ScriptHash")
		default:
			return nil, errors.New("Invalid address type")
		}

		script, err := txscript.PayToAddrScript(addr)
		if err != nil {
			return nil, err
		}

		value := recipient.Value.Int().Int64()
		if value < 0 {
			return nil, fmt.Errorf("expected value >= 0, got value %v", value)
		}
		msgTx.AddTxOut(wire.NewTxOut(value, script))
	}

	return &Tx{inputs: inputs, recipients: recipients, msgTx: msgTx, signed: false}, nil
}

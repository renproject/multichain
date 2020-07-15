package bitcoincompat

import (
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/renproject/pack"
)

// GatewayScript returns the Bitcoin gateway script that is used to when
// submitting lock-and-mint cross-chain transactions (when working with BTC).
func GatewayScript(gpubkey pack.Bytes, ghash pack.Bytes32) ([]byte, error) {
	pubKeyHash160 := btcutil.Hash160(gpubkey)
	return txscript.NewScriptBuilder().
		AddData(ghash.Bytes()).
		AddOp(txscript.OP_DROP).
		AddOp(txscript.OP_DUP).
		AddOp(txscript.OP_HASH160).
		AddData(pubKeyHash160).
		AddOp(txscript.OP_EQUALVERIFY).
		AddOp(txscript.OP_CHECKSIG).
		Script()
}

// GatewayPubKeyScript returns the pubkey script of a Bitcoin gateway script.
// This is the pubkey script that is expected to be in the underlying
// transaction output for lock-and-mint cross-chain transactions (when working
// with BTC).
func GatewayPubKeyScript(gpubkey pack.Bytes, ghash pack.Bytes32) ([]byte, error) {
	script, err := GatewayScript(gpubkey, ghash)
	if err != nil {
		return nil, fmt.Errorf("invalid script: %v", err)
	}
	pubKeyScript, err := txscript.NewScriptBuilder().
		AddOp(txscript.OP_HASH160).
		AddData(btcutil.Hash160(script)).
		AddOp(txscript.OP_EQUAL).Script()
	if err != nil {
		return nil, fmt.Errorf("invalid pubkeyscript: %v", err)
	}
	return pubKeyScript, nil
}

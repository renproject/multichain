package lbry_test

import (
	"bytes"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/utxo"
	"github.com/renproject/multichain/chain/lbry"
	"github.com/renproject/pack"
)

func publicLBRYFixture(t *testing.T) (*btcec.PrivateKey, []utxo.Input, []utxo.Recipient) {
	t.Helper()
	// Public test scalar 1; never a real wallet or a funded output.
	key, _ := btcec.PrivKeyFromBytes(btcec.S256(), []byte{1})
	owner, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(key.PubKey().SerializeCompressed()), &lbry.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}
	script, err := txscript.PayToAddrScript(owner)
	if err != nil {
		t.Fatal(err)
	}
	inputs := []utxo.Input{{Output: utxo.Output{
		Outpoint: utxo.Outpoint{Hash: pack.NewBytes(bytes.Repeat([]byte{0x42}, 32)), Index: pack.NewU32(3)},
		Value:    pack.NewU256FromUint64(100000), PubKeyScript: pack.NewBytes(script),
	}}}
	recipients := []utxo.Recipient{
		{To: address.Address("bCjHuFv7ggRwqouv9M2Wp5nHbZdhKdPGhd"), Value: pack.NewU256FromUint64(70000)},
		{To: address.Address(owner.EncodeAddress()), Value: pack.NewU256FromUint64(29000)},
	}
	return key, inputs, recipients
}

func TestLBRYTxPreservesOutputsAndFee(t *testing.T) {
	_, inputs, recipients := publicLBRYFixture(t)
	tx, err := lbry.NewTxBuilder(&lbry.MainNetParams).BuildTx(inputs, recipients)
	if err != nil {
		t.Fatal(err)
	}
	serialized, err := tx.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	var decoded wire.MsgTx
	if err := decoded.Deserialize(bytes.NewReader(serialized)); err != nil {
		t.Fatal(err)
	}
	if decoded.Version != 2 || len(decoded.TxIn) != 1 || len(decoded.TxOut) != 2 {
		t.Fatalf("unexpected transaction structure: %+v", decoded)
	}
	if decoded.TxIn[0].PreviousOutPoint.Index != 3 || !bytes.Equal(decoded.TxIn[0].PreviousOutPoint.Hash[:], inputs[0].Hash) {
		t.Fatal("outpoint changed")
	}
	var total int64
	for i, recipient := range recipients {
		addr, err := btcutil.DecodeAddress(string(recipient.To), &lbry.MainNetParams)
		if err != nil {
			t.Fatal(err)
		}
		expectedScript, err := txscript.PayToAddrScript(addr)
		if err != nil {
			t.Fatal(err)
		}
		if decoded.TxOut[i].Value != recipient.Value.Int().Int64() || !bytes.Equal(decoded.TxOut[i].PkScript, expectedScript) {
			t.Fatalf("recipient %d amount or script changed", i)
		}
		total += decoded.TxOut[i].Value
	}
	if fee := inputs[0].Value.Int().Int64() - total; fee != 1000 {
		t.Fatalf("fee: got %d, want 1000", fee)
	}
	outputs, err := tx.Outputs()
	if err != nil || len(outputs) != 2 {
		t.Fatalf("outputs: %v", err)
	}
	hash := decoded.TxHash()
	for i, output := range outputs {
		if !bytes.Equal(output.Hash, hash[:]) || output.Index.Uint32() != uint32(i) || output.Value.Int().Int64() != decoded.TxOut[i].Value {
			t.Fatalf("output %d no longer matches serialized transaction", i)
		}
	}
}

func TestLBRYTxSignaturesValidateAndBindOutputs(t *testing.T) {
	key, inputs, recipients := publicLBRYFixture(t)
	tx, err := lbry.NewTxBuilder(&lbry.MainNetParams).BuildTx(inputs, recipients)
	if err != nil {
		t.Fatal(err)
	}
	if err := tx.Sign(nil, pack.NewBytes(key.PubKey().SerializeCompressed())); err == nil {
		t.Fatal("accepted missing input signature")
	}
	hashes, err := tx.Sighashes()
	if err != nil || len(hashes) != 1 {
		t.Fatalf("sighashes: %v", err)
	}
	signature, err := key.Sign(hashes[0][:])
	if err != nil {
		t.Fatal(err)
	}
	var rsv pack.Bytes65
	r, s := signature.R.Bytes(), signature.S.Bytes()
	copy(rsv[32-len(r):32], r)
	copy(rsv[64-len(s):64], s)
	pubkey := pack.NewBytes(key.PubKey().SerializeCompressed())
	if err := tx.Sign([]pack.Bytes65{rsv}, pubkey); err != nil {
		t.Fatal(err)
	}
	if err := tx.Sign([]pack.Bytes65{rsv}, pubkey); err == nil {
		t.Fatal("signed the same transaction twice")
	}
	serialized, err := tx.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	var decoded wire.MsgTx
	if err := decoded.Deserialize(bytes.NewReader(serialized)); err != nil {
		t.Fatal(err)
	}
	check := func(tx *wire.MsgTx) error {
		engine, err := txscript.NewEngine(inputs[0].PubKeyScript, tx, 0, txscript.StandardVerifyFlags, nil, nil, inputs[0].Value.Int().Int64())
		if err != nil {
			return err
		}
		return engine.Execute()
	}
	if err := check(&decoded); err != nil {
		t.Fatalf("serialized P2PKH signature failed script validation: %v", err)
	}
	decoded.TxOut[0].Value++
	if err := check(&decoded); err == nil {
		t.Fatal("signature accepted an altered payment amount")
	}
}

func TestLBRYTxRejectsWrongNetworkRecipient(t *testing.T) {
	_, inputs, recipients := publicLBRYFixture(t)
	recipients[0].To = "mfWyW5fc9NUj75YAnFgoRLrjxgLDn2MMth"
	if _, err := lbry.NewTxBuilder(&lbry.MainNetParams).BuildTx(inputs, recipients); err == nil {
		t.Fatal("mainnet transaction accepted a testnet recipient")
	}
}

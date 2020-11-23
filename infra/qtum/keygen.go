package main

// DEZU: Straight copy of bitcoin's implementation with Qtums chain configs tacked on

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/renproject/id"
)

// DEZU: values take from from qtumsuit
var RegressionNetParams = chaincfg.Params{
	Name: "regtest",
	DefaultPort: "23888",

	Net: 0xe1c6ddfd,

	// Address encoding magics
	PubKeyHashAddrID: 120, // starts with m or n
	ScriptHashAddrID: 110, // starts with 2
	PrivateKeyID:     239, // starts with 9 (uncompressed) or c (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub

	// Human-readable part for Bech32 encoded segwit addresses, as defined in BIP 173.
	Bech32HRPSegwit: "qcrt",
}

func main() {
	privKey := id.NewPrivKey()
	wif, err := btcutil.NewWIF((*btcec.PrivateKey)(privKey), &RegressionNetParams, true)
	if err != nil {
		panic(err)
	}
	addrPubKeyHash, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &RegressionNetParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("QTUM_PK=%v\n", wif)
	fmt.Printf("QTUM_ADDRESS=%v\n", addrPubKeyHash)
}

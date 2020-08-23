package main

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/renproject/id"
)

func main() {
	privKey := id.NewPrivKey()
	wif, err := btcutil.NewWIF((*btcec.PrivateKey)(privKey), &chaincfg.RegressionNetParams, true)
	if err != nil {
		panic(err)
	}
	addrPubKeyHash, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.RegressionNetParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("BITCOIN_PK=%v\n", wif)
	fmt.Printf("BITCOIN_ADDRESS=%v\n", addrPubKeyHash)
}

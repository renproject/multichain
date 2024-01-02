package main

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

func main() {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		panic(err)
	}
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

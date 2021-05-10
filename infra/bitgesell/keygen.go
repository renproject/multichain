package main

import (
	"fmt"

	"github.com/bitgesellofficial/bgld/btcec"
	"github.com/bitgesellofficial/bglutil"

	"github.com/renproject/id"
	"github.com/renproject/multichain/chain/bitgesell"
)

func main() {
	privKey := id.NewPrivKey()
	wif, err := bglutil.NewWIF((*btcec.PrivateKey)(privKey), &bitgesell.RegressionNetParams, true)
	if err != nil {
		panic(err)
	}
	addrPubKeyHash, err := bglutil.NewAddressPubKeyHash(bglutil.Hash160(wif.SerializePubKey()), &bitgesell.RegressionNetParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("BGL_PK=%v\n", wif)
	fmt.Printf("BGL_ADDRESS=%v\n", addrPubKeyHash)
}

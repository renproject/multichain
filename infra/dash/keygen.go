package main

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"

	"github.com/renproject/id"
	"github.com/renproject/multichain/chain/dash"
)

func main() {
	privKey := id.NewPrivKey()
	wif, err := btcutil.NewWIF((*btcec.PrivateKey)(privKey), &dash.RegressionNetParams, true)
	if err != nil {
		panic(err)
	}
	addrPubKeyHash, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &dash.RegressionNetParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("DASH_PK=%v\n", wif)
	fmt.Printf("DASH_ADDRESS=%v\n", addrPubKeyHash)
}

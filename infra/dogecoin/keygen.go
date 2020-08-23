package main

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"

	"github.com/renproject/id"
	"github.com/renproject/multichain/chain/dogecoin"
)

func main() {
	privKey := id.NewPrivKey()
	wif, err := btcutil.NewWIF((*btcec.PrivateKey)(privKey), &dogecoin.RegressionNetParams, true)
	if err != nil {
		panic(err)
	}
	addrPubKeyHash, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &dogecoin.RegressionNetParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("DOGECOIN_PK=%v\n", wif)
	fmt.Printf("DOGECOIN_ADDRESS=%v\n", addrPubKeyHash)
}

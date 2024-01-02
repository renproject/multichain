package main

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"

	"github.com/renproject/multichain/chain/dogecoin"
)

func main() {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		panic(err)
	}
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

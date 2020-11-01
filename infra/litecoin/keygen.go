package main

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"

	"github.com/renproject/id"
	"github.com/renproject/multichain/chain/litecoin"
)

func main() {
	privKey := id.NewPrivKey()
	wif, err := btcutil.NewWIF((*btcec.PrivateKey)(privKey), &litecoin.RegressionNetParams, true)
	if err != nil {
		panic(err)
	}
	addrPubKeyHash, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &litecoin.RegressionNetParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("LITECOIN_PK=%v\n", wif)
	fmt.Printf("LITECOIN_ADDRESS=%v\n", addrPubKeyHash)
}

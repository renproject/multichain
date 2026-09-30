package main

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/renproject/id"
	"github.com/renproject/multichain/chain/lbry"
)

func main() {
	privKey := id.NewPrivKey()
	wif, err := btcutil.NewWIF((*btcec.PrivateKey)(privKey), &lbry.RegressionNetParams, true)
	if err != nil {
		panic(err)
	}
	addrPubKeyHash, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &lbry.RegressionNetParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("LBRY_PK=%v\n", wif)
	fmt.Printf("LBRY_ADDRESS=%v\n", addrPubKeyHash)
}

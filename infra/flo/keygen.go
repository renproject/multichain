package main

import (
	"fmt"
	"github.com/bitspill/flod/floec"
	"github.com/bitspill/floutil"

	"github.com/renproject/id"
	"github.com/renproject/multichain/chain/flo"
)

func main() {
	privKey := id.NewPrivKey()
	wif, err := floutil.NewWIF((*floec.PrivateKey)(privKey), &flo.RegressionNetParams, true)
	if err != nil {
		panic(err)
	}
	addrPubKeyHash, err := floutil.NewAddressPubKeyHash(floutil.Hash160(wif.SerializePubKey()), &flo.RegressionNetParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("FLO_PK=%v\n", wif)
	fmt.Printf("FLO_ADDRESS=%v\n", addrPubKeyHash)
}

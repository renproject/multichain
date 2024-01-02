package main

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/renproject/multichain/chain/digibyte"
)

func main() {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		panic(err)
	}
	wif, err := btcutil.NewWIF((*btcec.PrivateKey)(privKey), &digibyte.RegressionNetParams, true)
	if err != nil {
		panic(err)
	}
	addrPubKeyHash, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &digibyte.RegressionNetParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("DIGIBYTE_PK=%v\n", wif)
	fmt.Printf("DIGIBYTE_ADDRESS=%v\n", addrPubKeyHash)
}

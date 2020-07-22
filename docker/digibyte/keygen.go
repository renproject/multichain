package main

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"

	"github.com/renproject/multichain/chain/digibyte"
	"github.com/renproject/id"
)

func main() {
	// Use this for main net
	// var network = &digibyte.DigiByteMainNetParams
	var network = digibyte.DigiByteRegtestParams

	privKey := id.NewPrivKey()
	wif, err := btcutil.NewWIF((*btcec.PrivateKey)(privKey), network, true)
	if err != nil {
		panic(err)
	}
	addrPubKeyHash, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), network)
	if err != nil {
		panic(err)
	}
	fmt.Printf("DIGIBYTE_PK=%v\n", wif)
	fmt.Printf("DIGIBYTE_ADDRESS=%v\n", addrPubKeyHash)
}
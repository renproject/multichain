package main

import (
	"fmt"

	"github.com/decred/dcrd/dcrec"
	"github.com/decred/dcrd/dcrec/secp256k1/v3"
	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/decred/dcrd/chaincfg/v3"
)

func main() {
	simNetPrivKeyID := [2]byte{0x23, 0x07}
	privKey, err := secp256k1.GeneratePrivateKey()
	if err != nil {
	    panic(err) 
	}
	wif, err := dcrutil.NewWIF(privKey.Serialize(), simNetPrivKeyID, dcrec.STEcdsaSecp256k1)
	if err != nil {
		panic(err)
	}
	addrPubKeyHash, err := dcrutil.NewAddressPubKeyHash(dcrutil.Hash160(wif.PubKey()), chaincfg.SimNetParams(), dcrec.STEcdsaSecp256k1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("DECRED_PK=%v\n", wif)
	fmt.Printf("DECRED_ADDRESS=%v\n", addrPubKeyHash)
}

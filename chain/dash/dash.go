package dash

import (
	"github.com/btcsuite/btcd/chaincfg"
)

func init() {
	if err := chaincfg.Register(&MainNetParams); err != nil {
		panic(err)
	}
	if err := chaincfg.Register(&TestNetParams); err != nil {
		panic(err)
	}
	if err := chaincfg.Register(&RegressionNetParams); err != nil {
		panic(err)
	}
}

var MainNetParams = chaincfg.Params{
	Name: "mainnet",
	Net:  0xd9b4bef9,


    // Address encoding magics
    PubKeyHashAddrID: 0x00, // starts with 1
	ScriptHashAddrID: 0x05, // starts with 3
	PrivateKeyID:     0x80, // starts with 5 (uncompressed) or K (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173. Dash does not actually support this, but we do not want to
	// collide with real addresses, so we specify it.
	Bech32HRPSegwit: "dash",
}

var TestNetParams = chaincfg.Params{
	Name: "testnet",
	Net:  0x0709110b,


	// Address encoding magics
	PubKeyHashAddrID: 0x6f, // starts with m or n
	ScriptHashAddrID: 0xc4, // starts with 2
	PrivateKeyID:     0xef, // starts with 9 (uncompressed) or c (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173. Dash does not actually support this, but we do not want to
	// collide with real addresses, so we specify it.
	Bech32HRPSegwit: "dasht",
}

var RegressionNetParams = chaincfg.Params{
	Name: "regtest",

	// Dash has 0xdab5bffa as RegTest (same as Bitcoin's RegTest).
	// Setting it to an arbitrary value, so that we can
	// register the regtest network.
	Net: 0xdab5bffb,

	// Address encoding magics
	PubKeyHashAddrID: 0x6f, // starts with m or n
	ScriptHashAddrID: 0xc4, // starts with 2
	PrivateKeyID:     0xef, // starts with 9 (uncompressed) or c (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173. Dash does not actually support this, but we do not want to
	// collide with real addresses, so we specify it.
	Bech32HRPSegwit: "dashrt",
}

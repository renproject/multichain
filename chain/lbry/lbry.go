package lbry

import (
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
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

// genesisCoinbaseTx is the coinbase transaction for the genesis blocks for
// the main network, regression test network, and test network (version 3).
var genesisCoinbaseTx = wire.MsgTx{
	Version: 1,
	TxIn: []*wire.TxIn{
		{
			PreviousOutPoint: wire.OutPoint{
				Hash:  chainhash.Hash{},
				Index: 0xffffffff,
			},
			SignatureScript: []byte{
				0x04, 0xff, 0xff, 0x00, 0x1d, 0x01, 0x04, 0x17,
				0x69, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x20, 0x74,
				0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
				0x20, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67,
			},
			Sequence: 0xffffffff,
		},
	},
	TxOut: []*wire.TxOut{
		{
			Value: 0x12a05f200,
			PkScript: []byte{
				0x76, 0xa9, 0x14, 0x34, 0x59, 0x91, 0xdb, 0xf5,
				0x7b, 0xfb, 0x01, 0x4b, 0x87, 0x00, 0x6a, 0xcd,
				0xfa, 0xfb, 0xfc, 0x5f, 0xe8, 0x29, 0x2f, 0x88,
				0xac,
			},
		},
	},
	LockTime: 0,
}

// https://github.com/lbryio/lbrycrd/blob/master/src/chainparams.cpp#L19
var genesisMerkleRoot = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0xcc, 0x59, 0xe5, 0x9f, 0xf9, 0x7a, 0xc0, 0x92,
	0xb5, 0x5e, 0x42, 0x3a, 0xa5, 0x49, 0x51, 0x51,
	0xed, 0x6f, 0xb8, 0x05, 0x70, 0xa5, 0xbb, 0x78,
	0xcd, 0x5b, 0xd1, 0xc3, 0x82, 0x1c, 0x21, 0xb8,
})

var genesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: genesisMerkleRoot,        // b8211c82c3d15bcd78bba57005b86fed515149a53a425eb592c07af99fe559cc
		Timestamp:  time.Unix(1446058291, 0), // Wednesday, October 28, 2015 6:51:31 PM GMT
		Bits:       0x1f00ffff,
		Nonce:      1287,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

var genesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x63, 0xf4, 0x34, 0x6a, 0x4d, 0xb3, 0x4f, 0xdf,
	0xce, 0x29, 0xa7, 0x0f, 0x5e, 0x8d, 0x11, 0xf0,
	0x65, 0xf6, 0xb9, 0x16, 0x02, 0xb7, 0x03, 0x6c,
	0x7f, 0x22, 0xf3, 0xa0, 0x3b, 0x28, 0x89, 0x9c,
})

func newHashFromStr(hexStr string) *chainhash.Hash {
	hash, err := chainhash.NewHashFromStr(hexStr)
	if err != nil {
		panic(err)
	}
	return hash
}

// MainNetParams returns the chain configuration for mainnet.
var MainNetParams = chaincfg.Params{
	Name:        "mainnet",
	Net:         0xfae4aaf1,
	DefaultPort: "9246",

	// Chain parameters
	GenesisBlock: &genesisBlock,
	GenesisHash:  &genesisHash,

	// Address encoding magics
	PubKeyHashAddrID:        85,
	ScriptHashAddrID:        122,
	PrivateKeyID:            28,
	WitnessPubKeyHashAddrID: 0x06, // starts with p2
	WitnessScriptHashAddrID: 0x0A, // starts with 7Xh
	BIP0034Height:           1,
	BIP0065Height:           200000,
	BIP0066Height:           200000,

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "lbc",

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 0x8c,
}

// TestNetParams returns the chain configuration for testnet.
var TestNetParams = chaincfg.Params{
	Name:        "testnet",
	Net:         0xfae4aae1,
	DefaultPort: "19246",

	// Chain parameters
	GenesisBlock: &genesisBlock,
	GenesisHash:  &genesisHash,

	// Address encoding magics
	PubKeyHashAddrID: 111,
	ScriptHashAddrID: 196,
	PrivateKeyID:     239,

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with xpub

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "tlbc",

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 0x8c,
}

// RegressionNetParams returns the chain configuration for regression net.
var RegressionNetParams = chaincfg.Params{
	Name: "regtest",

	// LBRY has 0xfae4aad1 as RegTest (same as Bitcoin's RegTest).
	// Setting it to an arbitrary value (leet_hex(LBRY)), so that we can
	// register the regtest network.
	Net: 0xfae4aad1,

	// Address encoding magics
	PubKeyHashAddrID: 111,
	ScriptHashAddrID: 196,
	PrivateKeyID:     239,

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with xpub

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "rlbc",

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 0x8c,
}

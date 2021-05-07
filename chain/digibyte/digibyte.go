package digibyte

import (
	"math/big"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

func init() {
	if err := chaincfg.Register(&MainNetParams); err != nil {
		panic(err)
	}
	if err := chaincfg.Register(&TestnetParams); err != nil {
		panic(err)
	}
	if err := chaincfg.Register(&RegressionNetParams); err != nil {
		panic(err)
	}
}

var (
	bigOne       = big.NewInt(1)
	mainPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne)
)

const (
	// DeploymentTestDummy ...
	DeploymentTestDummy = iota

	// DeploymentCSV ...
	DeploymentCSV

	// DeploymentSegwit ...
	DeploymentSegwit

	// DefinedDeployments ...
	DefinedDeployments
)

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
				0x04, 0xff, 0xff, 0x00, 0x1d, 0x01, 0x04, 0x45, 0x55, 0x53, 0x41, 0x20, 0x54, 0x6f, 0x64, 0x61, /* |.......EUSA Toda| */
				0x79, 0x3a, 0x20, 0x31, 0x30, 0x2f, 0x4a, 0x61, 0x6e, 0x2f, 0x32, 0x30, 0x31, 0x34, 0x2c, 0x20, /* |y: 10/Jan/2014, | */
				0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x3a, 0x20, 0x44, 0x61, 0x74, 0x61, 0x20, 0x73, 0x74, 0x6f, /* |Target: Data sto| */
				0x6c, 0x65, 0x6e, 0x20, 0x66, 0x72, 0x6f, 0x6d, 0x20, 0x75, 0x70, 0x20, 0x74, 0x6f, 0x20, 0x31, /* |len from up to 1| */
				0x31, 0x30, 0x4d, 0x20, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x73, 0x61, 0x32, 0x30, /* |10M customers|    */
			},
			Sequence: 0xffffffff,
		},
	},
	TxOut: []*wire.TxOut{
		{
			Value: 0x12a05f200,
			PkScript: []byte{ // ToDo
				0x41, 0x04, 0x67, 0x8a, 0xfd, 0xb0, 0xfe, 0x55, /* |A.g....U| */
				0x48, 0x27, 0x19, 0x67, 0xf1, 0xa6, 0x71, 0x30, /* |H'.g..q0| */
				0xb7, 0x10, 0x5c, 0xd6, 0xa8, 0x28, 0xe0, 0x39, /* |..\..(.9| */
				0x09, 0xa6, 0x79, 0x62, 0xe0, 0xea, 0x1f, 0x61, /* |..yb...a| */
				0xde, 0xb6, 0x49, 0xf6, 0xbc, 0x3f, 0x4c, 0xef, /* |..I..?L.| */
				0x38, 0xc4, 0xf3, 0x55, 0x04, 0xe5, 0x1e, 0xc1, /* |8..U....| */
				0x12, 0xde, 0x5c, 0x38, 0x4d, 0xf7, 0xba, 0x0b, /* |..\8M...| */
				0x8d, 0x57, 0x8a, 0x4c, 0x70, 0x2b, 0x6b, 0xf1, /* |.W.Lp+k.| */
				0x1d, 0x5f, 0xac, /* |._.| */
			},
		},
	},
	LockTime: 0,
}

// USA Today: 10/Jan/2014, Target: Data stolen from up to 110M customers
var genesisMerkleRoot = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x96, 0x84, 0x1e, 0x6e, 0xcc, 0x8d,
	0xc9, 0x64, 0x3a, 0xad, 0xdf, 0xb6,
	0xfc, 0xd6, 0x16, 0xe0, 0x8f, 0x07,
	0x77, 0xc8, 0x7b, 0x50, 0x8f, 0x1c,
	0x9f, 0xb3, 0x5e, 0x46, 0x1b, 0xea,
	0x97, 0x74,
})

var genesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: genesisMerkleRoot,        // 7497ea1b465eb39f1c8f507bc877078fe016d6fcb6dfad3a64c98dcc6e1e8496
		Timestamp:  time.Unix(1389388394, 0), // 2014-01-10T21:13:14.000Z
		Bits:       0x1e0ffff0,               // 486604799 [00000000ffff0000000000000000000000000000000000000000000000000000]
		Nonce:      2447652,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

var genesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x96, 0x84, 0x1e, 0x6e, 0xcc, 0x8d, 0xc9, 0x64,
	0x3a, 0xad, 0xdf, 0xb6, 0xfc, 0xd6, 0x16, 0xe0,
	0x8f, 0x07, 0x77, 0xc8, 0x7b, 0x50, 0x8f, 0x1c,
	0x9f, 0xb3, 0x5e, 0x46, 0x1b, 0xea, 0x97, 0x74,
})

func newHashFromStr(hexStr string) *chainhash.Hash {
	hash, err := chainhash.NewHashFromStr(hexStr)
	if err != nil {
		panic(err)
	}
	return hash
}

// MainNetParams returns the chain configuration for mainnet
var MainNetParams = chaincfg.Params{
	Name:        "mainnet",
	Net:         0xdab6c3fa,
	DefaultPort: "12024",

	// Chain parameters
	GenesisBlock: &genesisBlock,
	GenesisHash:  &genesisHash,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "dgb", // always bc for main net

	// Address encoding magics
	PubKeyHashAddrID:        0x1e, // starts with 1
	ScriptHashAddrID:        0x3f, // starts with 3
	PrivateKeyID:            0x80, // starts with 5 (uncompressed) or K (compressed)
	WitnessPubKeyHashAddrID: 0x06, // starts with p2
	WitnessScriptHashAddrID: 0x0A, // starts with 7Xh

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 0x14,
}

// TestnetParams returns the chain configuration for testnet
var TestnetParams = chaincfg.Params{
	Name: "testnet",

	// DigiByte has 0xdab5bffa as RegTest (same as Bitcoin's RegTest).
	// Setting it to an arbitrary value (leet_hex(digibyte)), so that we can
	// register the regtest network.
	// DigiByte Core Developers will change this soon.
	Net:         0xddbdc8fd,
	DefaultPort: "12026",

	// Chain parameters
	GenesisBlock: &genesisBlock,
	GenesisHash:  &genesisHash,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "dgbt", // always bc for main net

	// Address encoding magics
	PubKeyHashAddrID:        0x7e, // starts with 1
	ScriptHashAddrID:        0x8c, // starts with 3
	PrivateKeyID:            0xfe, // starts with 5 (uncompressed) or K (compressed)
	WitnessPubKeyHashAddrID: 0x06, // starts with p2
	WitnessScriptHashAddrID: 0x0A, // starts with 7Xh

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with xpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 0x14,
}

// RegressionNetParams returns the chain configuration for regression net
var RegressionNetParams = chaincfg.Params{
	Name: "regtest",

	// DigiByte has 0xdab5bffa as RegTest (same as Bitcoin's RegTest).
	// Setting it to an arbitrary value (leet_hex(digibyte)), so that we can
	// register the regtest network.
	// DigiByte Core Developers will change this soon.
	Net:         0xd191841e,
	DefaultPort: "18444",

	// Chain parameters
	GenesisBlock: &genesisBlock,
	GenesisHash:  &genesisHash,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "dgbrt", // always bc for main net

	// Address encoding magics
	PubKeyHashAddrID:        0x7e, // starts with 1
	ScriptHashAddrID:        0x8c, // starts with 3
	PrivateKeyID:            0xfe, // starts with 5 (uncompressed) or K (compressed)
	WitnessPubKeyHashAddrID: 0x06, // starts with p2
	WitnessScriptHashAddrID: 0x0A, // starts with 7Xh

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with xpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 0x14,
}

package zcash

import (
	"github.com/btcsuite/btcd/chaincfg"
)

const (
	sighashMask                 = 0x1f
	blake2BSighash              = "ZcashSigHash"
	prevoutsHashPersonalization = "ZcashPrevoutHash"
	sequenceHashPersonalization = "ZcashSequencHash"
	outputsHashPersonalization  = "ZcashOutputsHash"

	versionOverwinter        int32  = 3
	versionOverwinterGroupID uint32 = 0x3C48270
	versionSapling                  = 4
	versionSaplingGroupID           = 0x892f2085
)

type Params struct {
	// TODO: We do not actually need to embed the entire chaincfg params object.
	*chaincfg.Params

	P2SHPrefix  []byte
	P2PKHPrefix []byte
	Upgrades    []ParamsUpgrade
}

type ParamsUpgrade struct {
	ActivationHeight uint32
	BranchID         []byte
}

var (
	witnessMarkerBytes = []byte{0x00, 0x01}

	MainNetParams = Params{
		Params: &chaincfg.MainNetParams,

		P2PKHPrefix: []byte{0x1C, 0xB8},
		P2SHPrefix:  []byte{0x1C, 0xBD},
		Upgrades: []ParamsUpgrade{
			{0, []byte{0x00, 0x00, 0x00, 0x00}},
			{347500, []byte{0x19, 0x1B, 0xA8, 0x5B}},
			{419200, []byte{0xBB, 0x09, 0xB8, 0x76}},
			{653600, []byte{0x60, 0x0E, 0xB4, 0x2B}},
		},
	}
	TestNet3Params = Params{
		Params: &chaincfg.TestNet3Params,

		P2PKHPrefix: []byte{0x1D, 0x25},
		P2SHPrefix:  []byte{0x1C, 0xBA},
		Upgrades: []ParamsUpgrade{
			{0, []byte{0x00, 0x00, 0x00, 0x00}},
			{207500, []byte{0x19, 0x1B, 0xA8, 0x5B}},
			{280000, []byte{0xBB, 0x09, 0xB8, 0x76}},
			{584000, []byte{0x60, 0x0E, 0xB4, 0x2B}},
		},
	}
	RegressionNetParams = Params{
		Params: &chaincfg.RegressionNetParams,

		P2PKHPrefix: []byte{0x1D, 0x25},
		P2SHPrefix:  []byte{0x1C, 0xBA},
		Upgrades: []ParamsUpgrade{
			{0, []byte{0x00, 0x00, 0x00, 0x00}},
			{60, []byte{0x19, 0x1B, 0xA8, 0x5B}},
			{80, []byte{0xBB, 0x09, 0xB8, 0x76}},
			{100, []byte{0x60, 0x0E, 0xB4, 0x2B}},
		},
	}
)

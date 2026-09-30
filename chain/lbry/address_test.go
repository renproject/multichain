package lbry_test

import (
	"bytes"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/chain/bitcoin"
	"github.com/renproject/multichain/chain/lbry"
)

// These Base58Check fixtures encode the public hash 000102...1213. Testnet and
// regtest intentionally share address versions; they are not distinguishable by
// their legacy addresses. No key, node or funded account is used.
func TestLBRYAddressRoundTrip(t *testing.T) {
	cases := []struct {
		name    string
		params  *chaincfg.Params
		encoded string
		version byte
	}{
		{"mainnet-p2pkh", &lbry.MainNetParams, "bCjHuFv7ggRwqouv9M2Wp5nHbZdhKdPGhd", 85},
		{"mainnet-p2sh", &lbry.MainNetParams, "r6EcLGwmxMaN6r583sMJkirPtEAbfrdFcU", 122},
		{"testnet-p2pkh", &lbry.TestNetParams, "mfWyW5fc9NUj75YAnFgoRLrjxgLDn2MMth", 111},
		{"testnet-p2sh", &lbry.TestNetParams, "2MsFFCK16VhsCcvPXruztdzzcTZEQCbNKjJ", 196},
		{"regtest-p2pkh", &lbry.RegressionNetParams, "mfWyW5fc9NUj75YAnFgoRLrjxgLDn2MMth", 111},
		{"regtest-p2sh", &lbry.RegressionNetParams, "2MsFFCK16VhsCcvPXruztdzzcTZEQCbNKjJ", 196},
	}
	payload := make([]byte, 20)
	for i := range payload {
		payload[i] = byte(i)
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var codec lbry.AddressEncodeDecoder = bitcoin.NewAddressEncodeDecoder(tc.params)
			raw, err := codec.DecodeAddress(address.Address(tc.encoded))
			if err != nil {
				t.Fatal(err)
			}
			if len(raw) != 25 || raw[0] != tc.version || !bytes.Equal(raw[1:21], payload) {
				t.Fatalf("incorrect network or payload: %x", raw)
			}
			encoded, err := codec.EncodeAddress(raw)
			if err != nil || encoded != address.Address(tc.encoded) {
				t.Fatalf("round trip: got %q, %v", encoded, err)
			}
		})
	}
}

func TestLBRYAddressRejectsMalformedAndWrongNetwork(t *testing.T) {
	decoder := bitcoin.NewAddressDecoder(&lbry.MainNetParams)
	for _, invalid := range []string{
		"", "bCjHuFv7ggRwqouv9M2Wp5nHbZdhKdPGh1", "not-an-address", "\u0100",
		"mfWyW5fc9NUj75YAnFgoRLrjxgLDn2MMth", // LBRY testnet/regtest
		"1BoatSLRHtKNngkdXEeobR76b53LETtpyT", // Bitcoin mainnet
	} {
		t.Run(invalid, func(t *testing.T) {
			if _, err := decoder.DecodeAddress(address.Address(invalid)); err == nil {
				t.Fatal("accepted invalid or wrong-network address")
			}
		})
	}
}

func TestLBRYAddressChainMappings(t *testing.T) {
	if multichain.LBC.OriginChain() != multichain.LBRY || multichain.LBRY.NativeAsset() != multichain.LBC {
		t.Fatal("LBC and LBRY must map to each other")
	}
	if multichain.LBC.ChainType() != multichain.ChainTypeUTXOBased || !multichain.LBRY.IsUTXOBased() {
		t.Fatal("LBRY must use UTXO semantics")
	}
	// The merge must preserve current upstream mappings alongside the addition.
	for chain, asset := range map[multichain.Chain]multichain.Asset{
		multichain.Polygon: multichain.MATIC, multichain.Moonbeam: multichain.GLMR,
		multichain.Kava: multichain.KAVA, multichain.Optimism: multichain.Asset("oETH"),
	} {
		if chain.NativeAsset() != asset || asset.OriginChain() != chain || !chain.IsAccountBased() {
			t.Fatalf("lost upstream mapping for %s", chain)
		}
	}
}

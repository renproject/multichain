package acala_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/acala"
	"github.com/renproject/multichain/chain/bitcoin"
	"github.com/renproject/pack"
)

var _ = Describe("Mint Burn", func() {
	client, err := acala.NewClient(acala.DefaultClientOptions())
	Expect(err).NotTo(HaveOccurred())

	opts := types.SerDeOptions{NoPalletIndices: true}
	types.SetSerDeOptions(opts)

	Context("when minting over renbridge", func() {
		It("should succeed", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Ignore recipient
			pHash, nHash, sig, amount, _ := getMintParams()

			txhash, err := client.Mint(ctx, pHash, nHash, sig, amount)
			Expect(err).NotTo(HaveOccurred())

			fmt.Printf("txhash = %v\n", txhash)
		})
	})

	Context("when burning over renbridge", func() {
		It("should succeed", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Ignore phash, nhash, sig
			_, _, _, amount, recipient := getMintParams()

			txhash, err := client.Burn(ctx, recipient, amount)
			Expect(err).NotTo(HaveOccurred())

			fmt.Printf("txhash = %v\n", txhash)
		})
	})

	Context("when reading burn info", func() {
		It("should succeed", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// FIXME: set this appropriately
			blockheight := pack.U64(uint64(1))
			amount, recipient, confs, err := client.BurnEvent(ctx, blockheight)
			Expect(err).NotTo(HaveOccurred())

			fmt.Printf("amount = %v\n", amount)
			fmt.Printf("recipient = %v\n", recipient)
			fmt.Printf("confs = %v\n", confs)
		})
	})
})

func getMintParams() (pack.Bytes32, pack.Bytes32, pack.Bytes65, uint64, pack.Bytes) {
	pHashHex := "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
	nHashHex := "0x1f1a01537e418859cd99eb15099dcdcb98483ad723cd20ccaa5a2677b755572b"
	sigHex := "0x60930a2c1c933c30bb7f88d6183e81a71394d81ead26d68a3b12d6b4efdc3ef563f91a945375b71d78accb5860e8154bc01681577db544e2e53611aa14613a9c1b"
	amount := uint64(95000)
	recipient := multichain.Address(pack.String("miMi2VET41YV1j6SDNTeZoPBbmH8B4nEx6"))

	pHashBytes, err := types.HexDecodeString(pHashHex)
	Expect(err).NotTo(HaveOccurred())
	nHashBytes, err := types.HexDecodeString(nHashHex)
	Expect(err).NotTo(HaveOccurred())
	sigBytes, err := types.HexDecodeString(sigHex)
	Expect(err).NotTo(HaveOccurred())

	btcEncodeDecoder := bitcoin.NewAddressEncodeDecoder(&chaincfg.RegressionNetParams)
	rawAddr, err := btcEncodeDecoder.DecodeAddress(recipient)
	Expect(err).NotTo(HaveOccurred())

	var pHash [32]byte
	var nHash [32]byte
	var sig [65]byte
	copy(pHash[:], pHashBytes)
	copy(nHash[:], nHashBytes)
	copy(sig[:], sigBytes)

	return pack.Bytes32(pHash), pack.Bytes32(nHash), pack.Bytes65(sig), amount, pack.Bytes(rawAddr)
}

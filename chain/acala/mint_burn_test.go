package acala_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/centrifuge/go-substrate-rpc-client/signature"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/renproject/id"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/acala"
	"github.com/renproject/multichain/chain/bitcoin"
	"github.com/renproject/multichain/chain/ethereum"
	"github.com/renproject/pack"
	"github.com/renproject/surge"
)

var _ = Describe("Mint Burn", func() {
	client, err := acala.NewClient(acala.DefaultClientOptions())
	Expect(err).NotTo(HaveOccurred())

	Context("when minting over renbridge", func() {
		It("should succeed", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Ignore recipient
			alice, phash, nhash, sig, amount, _ := constructMintParams()

			txhash, err := client.Mint(ctx, alice, phash, nhash, sig, amount)
			Expect(err).NotTo(HaveOccurred())

			fmt.Printf("mint tx = %v\n", hex.EncodeToString(txhash))
		})
	})

	Context("when burning over renbridge", func() {
		It("should succeed", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Ignore phash, nhash, sig
			alice, _, _, _, amount, recipient := constructMintParams()

			txhash, err := client.Burn(ctx, alice, recipient, amount)
			Expect(err).NotTo(HaveOccurred())

			fmt.Printf("burn tx = %v\n", hex.EncodeToString(txhash))
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

func constructMintParams() (signature.KeyringPair, pack.Bytes32, pack.Bytes32, pack.Bytes65, uint64, [20]byte) {
	// Get RenVM priv key.
	renVmPrivKeyBytes, err := hex.DecodeString("c44700049a72c02bbacbec25551190427315f046c1f656f23884949da3fbdc3a")
	Expect(err).NotTo(HaveOccurred())
	renVmPrivKey := id.PrivKey{}
	err = surge.FromBinary(&renVmPrivKey, renVmPrivKeyBytes)
	Expect(err).NotTo(HaveOccurred())

	// Get random pHash and nHash.
	phashBytes := make([]byte, 32)
	nhashBytes := make([]byte, 32)
	_, err = rand.Read(phashBytes)
	Expect(err).NotTo(HaveOccurred())
	_, err = rand.Read(nhashBytes)
	Expect(err).NotTo(HaveOccurred())

	// Amount to be minted.
	amount := uint64(25000)

	// Selector for this cross-chain mint.
	token, err := hex.DecodeString("0000000000000000000000000a9add98c076448cbcfacf5e457da12ddbef4a8f")
	Expect(err).NotTo(HaveOccurred())
	token32 := [32]byte{}
	copy(token32[:], token[:])

	// Initialise message args
	sighash32 := [32]byte{}
	phash32 := [32]byte{}
	nhash32 := [32]byte{}
	to := [32]byte{}
	rawAddr, err := hex.DecodeString("d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d") // Alice.
	Expect(err).NotTo(HaveOccurred())

	// Get message sighash.
	copy(to[:], rawAddr)
	copy(phash32[:], phashBytes)
	copy(nhash32[:], nhashBytes)
	copy(sighash32[:], crypto.Keccak256(ethereum.Encode(
		pack.Bytes32(phash32),
		pack.NewU256FromUint64(amount),
		pack.Bytes32(token32),
		pack.Bytes32(to),
		pack.Bytes32(nhash32),
	)))

	// Sign the sighash.
	hash := id.Hash(sighash32)
	sig65, err := renVmPrivKey.Sign(&hash)
	Expect(err).NotTo(HaveOccurred())
	sig65[64] = sig65[64] + 27

	// Get the address of the burn recipient.
	recipientAddr := multichain.Address(pack.String("miMi2VET41YV1j6SDNTeZoPBbmH8B4nEx6"))
	btcEncodeDecoder := bitcoin.NewAddressEncodeDecoder(&chaincfg.RegressionNetParams)
	rawRecipientAddr, err := btcEncodeDecoder.DecodeAddress(recipientAddr)
	Expect(err).NotTo(HaveOccurred())
	recipient := [20]byte{}
	copy(recipient[:], rawRecipientAddr)

	return signature.TestKeyringPairAlice, pack.Bytes32(phash32), pack.Bytes32(nhash32), pack.Bytes65(sig65), amount, recipient
}

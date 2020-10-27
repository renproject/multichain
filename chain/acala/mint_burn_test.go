package acala_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/centrifuge/go-substrate-rpc-client/signature"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/renproject/id"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/multichain/chain/acala"
	"github.com/renproject/multichain/chain/bitcoin"
	"github.com/renproject/multichain/chain/ethereum"
	"github.com/renproject/pack"
	"github.com/renproject/surge"
)

var _ = FDescribe("Mint Burn", func() {
	r := rand.New(rand.NewSource(GinkgoRandomSeed()))

	blockhash, extSign := pack.Bytes32{}, pack.Bytes{}
	balanceBefore, balanceAfter := pack.U256{}, pack.U256{}

	client, err := acala.NewClient(acala.DefaultClientOptions())
	Expect(err).NotTo(HaveOccurred())

	alice, phash, nhash, sig, mintAmount, burnAmount, recipient := constructMintParams(r)

	Context("when minting over renbridge", func() {
		It("should succeed", func() {
			balanceBefore, err = client.Balance(alice)
			if err != nil {
				// This means there are no tokens allocated for that address.
				Expect(err).To(Equal(fmt.Errorf("get storage: <nil>")))
				balanceBefore = pack.NewU256FromUint64(uint64(0))
			}

			_, err = client.Mint(alice, phash, nhash, sig, mintAmount)
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(5 * time.Second)

			balanceAfter, err = client.Balance(alice)
			Expect(err).NotTo(HaveOccurred())

			Expect(balanceBefore.Add(pack.NewU256FromUint64(mintAmount))).To(Equal(balanceAfter))
		})
	})

	Context("when burning over renbridge", func() {
		It("should succeed", func() {
			balanceBefore, err = client.Balance(alice)
			Expect(err).NotTo(HaveOccurred())

			blockhash, extSign, err = client.Burn(alice, recipient, burnAmount)
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(5 * time.Second)

			balanceAfter, err = client.Balance(alice)
			Expect(err).NotTo(HaveOccurred())

			Expect(balanceBefore.Sub(pack.NewU256FromUint64(burnAmount))).To(Equal(balanceAfter))
		})
	})

	Context("when reading burn log", func() {
		It("should succeed", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			input := acala.BurnLogInput{
				Blockhash: blockhash,
				ExtSign:   extSign,
			}
			calldata, err := surge.ToBinary(input)
			Expect(err).NotTo(HaveOccurred())
			outputBytes, err := client.ContractCall(ctx, multichain.Address(""), contract.CallData(calldata))
			Expect(err).NotTo(HaveOccurred())

			output := acala.BurnLogOutput{}
			Expect(surge.FromBinary(&output, outputBytes)).To(Succeed())

			Expect(output.Amount).To(Equal(pack.NewU256FromUint64(burnAmount)))
			Expect(output.Recipient).To(Equal(multichain.RawAddress(recipient)))
			Expect(output.Confs).To(BeNumerically(">", 0))
		})
	})
})

func constructMintParams(r *rand.Rand) (signature.KeyringPair, pack.Bytes32, pack.Bytes32, pack.Bytes65, uint64, uint64, pack.Bytes) {
	// Get RenVM priv key.
	renVmPrivKeyBytes, err := hex.DecodeString("c44700049a72c02bbacbec25551190427315f046c1f656f23884949da3fbdc3a")
	Expect(err).NotTo(HaveOccurred())
	renVmPrivKey := id.PrivKey{}
	err = surge.FromBinary(&renVmPrivKey, renVmPrivKeyBytes)
	Expect(err).NotTo(HaveOccurred())

	// Get random pHash and nHash.
	phashBytes := make([]byte, 32)
	nhashBytes := make([]byte, 32)
	_, err = r.Read(phashBytes)
	Expect(err).NotTo(HaveOccurred())
	_, err = r.Read(nhashBytes)
	Expect(err).NotTo(HaveOccurred())

	// Amount to be minted/burnt.
	mintAmount := uint64(100000)
	burnAmount := uint64(25000)

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
		pack.NewU256FromUint64(mintAmount),
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
	recipientRawAddr, err := btcEncodeDecoder.DecodeAddress(recipientAddr)
	Expect(err).NotTo(HaveOccurred())

	return signature.TestKeyringPairAlice, pack.Bytes32(phash32), pack.Bytes32(nhash32), pack.Bytes65(sig65), mintAmount, burnAmount, pack.Bytes(recipientRawAddr)
}

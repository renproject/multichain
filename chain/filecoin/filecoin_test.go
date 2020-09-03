package filecoin_test

import (
	"context"
	"fmt"
	"os"
	"time"

	filaddress "github.com/filecoin-project/go-address"
	"github.com/renproject/id"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/filecoin"
	"github.com/renproject/pack"
	"github.com/renproject/surge"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filecoin", func() {
	Context("when broadcasting a tx", func() {
		It("should work", func() {
			// create context for the test
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// instantiate the client
			client, err := filecoin.NewClient(
				filecoin.DefaultClientOptions().
					WithAddress("ws://127.0.0.1:1234/rpc/v0").
					WithAuthToken("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.673MLa4AmbhNeC1Hj2Bn6c4t_ci68I0amkqAEHea8ik"),
			)
			Expect(err).ToNot(HaveOccurred())

			// read the private key that we will send transactions from
			filPrivKeyStr := os.Getenv("FILECOIN_PK")
			if filPrivKeyStr == "" {
				panic("FILECOIN_PK is undefined")
			}
			filAddress := os.Getenv("FILECOIN_ADDRESS")
			if filAddress == "" {
				panic("FILECOIN_ADDRESS is undefined")
			}
			filPrivKey := id.PrivKey{}
			err = surge.FromBinary(&filPrivKey, []byte(filPrivKeyStr))
			Expect(err).ToNot(HaveOccurred())
			filPubKey := filPrivKey.PubKey()
			filPubKeyCompressed, err := surge.ToBinary(filPubKey)
			Expect(err).NotTo(HaveOccurred())

			// random recipient
			recipientPK := id.NewPrivKey()
			recipientPubKey := recipientPK.PubKey()
			recipientPubKeyCompressed, err := surge.ToBinary(recipientPubKey)
			Expect(err).NotTo(HaveOccurred())
			recipientAddr, err := filaddress.NewSecp256k1Address(recipientPubKeyCompressed)

			// construct the transaction builder
			gasPrice := pack.NewU256FromU64(pack.NewU64(100))
			gasLimit := pack.NewU256FromU64(pack.NewU64(100000))
			amount := pack.NewU256FromU64(pack.NewU64(100))
			nonce := pack.NewU256FromU64(pack.NewU64(0))
			payload := pack.Bytes(nil)
			filTxBuilder := filecoin.NewTxBuilder(gasPrice, gasLimit)

			// build the transaction
			tx, err := filTxBuilder.BuildTx(
				multichain.Address(pack.String(filAddress)),
				multichain.Address(pack.String(recipientAddr.String())),
				amount, nonce,
				pack.U256{}, pack.U256{},
				payload,
			)
			Expect(err).ToNot(HaveOccurred())

			// Sign the filecoin-side lock transaction
			txSighashes, err := tx.Sighashes()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(txSighashes)).To(Equal(1))
			Expect(len(txSighashes[0])).To(Equal(32))
			sighash32 := [32]byte{}
			for i, b := range []byte(txSighashes[0]) {
				sighash32[i] = b
			}
			hash := id.Hash(sighash32)
			sig, err := filPrivKey.Sign(&hash)
			Expect(err).NotTo(HaveOccurred())
			sigBytes, err := surge.ToBinary(sig)
			Expect(err).NotTo(HaveOccurred())
			txSignature := pack.NewBytes(sigBytes)
			Expect(tx.Sign([]pack.Bytes{txSignature}, pack.NewBytes(filPubKeyCompressed))).To(Succeed())

			// submit the transaction
			txHash := tx.Hash()
			fmt.Printf("tx     = %v\n", tx)
			fmt.Printf("txhash = %v\n", txHash)
			err = client.SubmitTx(ctx, tx)
			Expect(err).ToNot(HaveOccurred())

			// wait for the transaction to be included in a block
			for {
				time.Sleep(time.Second)
				fetchedTx, confs, err := client.Tx(ctx, txHash)
				Expect(err).ToNot(HaveOccurred())
				Expect(fetchedTx.Hash().Equal(txHash)).To(BeTrue())
				Expect(confs).To(BeNumerically(">=", 0))
				if confs >= 1 {
					break
				}
			}
		})
	})
})

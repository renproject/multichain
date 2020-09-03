package filecoin_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	filaddress "github.com/filecoin-project/go-address"
	filtypes "github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
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
			senderPrivKeyStr := os.Getenv("FILECOIN_PK")
			if senderPrivKeyStr == "" {
				panic("FILECOIN_PK is undefined")
			}
			var ki filtypes.KeyInfo
			data, err := hex.DecodeString(senderPrivKeyStr)
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal(data, &ki)
			Expect(err).NotTo(HaveOccurred())
			senderPrivKey := id.PrivKey{}
			err = surge.FromBinary(&senderPrivKey, ki.PrivateKey)
			Expect(err).NotTo(HaveOccurred())

			// read sender's address into the filecoin-compatible format
			senderAddr := os.Getenv("FILECOIN_ADDRESS")
			if senderAddr == "" {
				panic("FILECOIN_ADDRESS is undefined")
			}
			senderFilAddr, err := filaddress.NewFromString(string(senderAddr))
			Expect(err).NotTo(HaveOccurred())

			// FIXME
			// random recipient
			recipientPK := id.NewPrivKey()
			recipientPubKey := recipientPK.PubKey()
			recipientPubKeyCompressed, err := surge.ToBinary(recipientPubKey)
			Expect(err).NotTo(HaveOccurred())
			recipientFilAddr, err := filaddress.NewSecp256k1Address(recipientPubKeyCompressed)
			Expect(err).NotTo(HaveOccurred())

			// construct the transaction builder
			gasFeeCap := pack.NewU256FromU64(pack.NewU64(149838))
			gasLimit := pack.NewU256FromU64(pack.NewU64(495335))
			gasPremium := pack.NewU256FromU64(pack.NewU64(149514))
			filTxBuilder := filecoin.NewTxBuilder(gasFeeCap, gasLimit, gasPremium)

			// build the transaction
			amount := pack.NewU256FromU64(pack.NewU64(100000000))
			nonce := pack.NewU256FromU64(pack.NewU64(0))
			payload := pack.Bytes(nil)
			tx, err := filTxBuilder.BuildTx(
				multichain.Address(pack.String(senderFilAddr.String())),
				multichain.Address(pack.String(recipientFilAddr.String())),
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
			sig, err := senderPrivKey.Sign(&hash)
			Expect(err).NotTo(HaveOccurred())
			sigBytes, err := surge.ToBinary(sig)
			Expect(err).NotTo(HaveOccurred())
			txSignature := pack.NewBytes(sigBytes)
			Expect(tx.Sign([]pack.Bytes{txSignature}, []byte{})).To(Succeed())

			// submit the transaction
			txHash := tx.Hash()
			txID, err := cid.Parse([]byte(txHash))
			Expect(err).NotTo(HaveOccurred())
			fmt.Printf("msgID = %v\n", txID)
			err = client.SubmitTx(ctx, tx)
			Expect(err).ToNot(HaveOccurred())

			// wait for the transaction to be included in a block
			for {
				time.Sleep(2 * time.Second)
				fetchedTx, confs, err := client.Tx(ctx, txHash)
				Expect(err).ToNot(HaveOccurred())
				if fetchedTx != nil {
					Expect(confs).To(BeNumerically(">=", 0))
					break
				}
			}
		})
	})
})

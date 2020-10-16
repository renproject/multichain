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
					WithAuthToken(fetchAuthToken()),
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

			// random recipient
			recipientPK := id.NewPrivKey()
			recipientPubKey := recipientPK.PubKey()
			recipientPubKeyCompressed, err := surge.ToBinary(recipientPubKey)
			Expect(err).NotTo(HaveOccurred())
			recipientFilAddr, err := filaddress.NewSecp256k1Address(recipientPubKeyCompressed)
			Expect(err).NotTo(HaveOccurred())

			// get good gas estimates
			gasLimit := uint64(2200000)
			gasEstimator := filecoin.NewGasEstimator(client, int64(gasLimit))
			gasPremium, gasFeeCap, err := gasEstimator.EstimateGas(ctx)
			Expect(err).ToNot(HaveOccurred())

			// construct the transaction builder
			filTxBuilder := filecoin.NewTxBuilder()

			// build the transaction
			sender := multichain.Address(pack.String(senderFilAddr.String()))
			amount := pack.NewU256FromU64(pack.NewU64(100000000))
			nonce, err := client.AccountNonce(ctx, sender)
			Expect(err).ToNot(HaveOccurred())

			tx, err := filTxBuilder.BuildTx(
				ctx,
				sender,
				multichain.Address(pack.String(recipientFilAddr.String())),
				amount, // amount
				nonce,  // nonce
				pack.NewU256FromU64(pack.NewU64(gasLimit)), // gasLimit
				gasPremium,
				gasFeeCap,
				pack.Bytes(nil), // payload
			)
			Expect(err).ToNot(HaveOccurred())

			// Sign the filecoin-side lock transaction
			txSighashes, err := tx.Sighashes()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(txSighashes)).To(Equal(1))
			Expect(len(txSighashes[0])).To(Equal(32))
			sighash32 := txSighashes[0]
			hash := id.Hash(sighash32)
			sig, err := senderPrivKey.Sign(&hash)
			Expect(err).NotTo(HaveOccurred())
			sigBytes, err := surge.ToBinary(sig)
			Expect(err).NotTo(HaveOccurred())
			txSignature := pack.Bytes65{}
			copy(txSignature[:], sigBytes)
			Expect(tx.Sign([]pack.Bytes65{txSignature}, []byte{})).To(Succeed())

			// submit the transaction
			txHash := tx.Hash()
			txID, err := cid.Parse([]byte(txHash))
			Expect(err).NotTo(HaveOccurred())
			fmt.Printf("msgID = %v\n", txID)
			err = client.SubmitTx(ctx, tx)
			Expect(err).ToNot(HaveOccurred())

			// Wait slightly before we query the chain's node.
			time.Sleep(time.Second)

			// wait for the transaction to be included in a block
			for {
				// Loop until the transaction has at least a few confirmations.
				tx, confs, err := client.Tx(ctx, txHash)
				if err == nil {
					Expect(confs.Uint64()).To(BeNumerically(">", 0))
					Expect(tx.From()).To(Equal(sender))
					Expect(tx.Value()).To(Equal(amount))
					break
				}

				// wait and retry querying for the transaction
				time.Sleep(5 * time.Second)
			}
		})
	})
})

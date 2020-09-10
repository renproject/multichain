package filecoin_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
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

			// fetch the auth token from filecoin's running docker container
			cmd := exec.Command("docker", "exec", "infra_filecoin_1", "/bin/bash", "-c", "/app/lotus auth api-info --perm admin")
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
				panic(fmt.Sprintf("could not run command: %v", err))
			}
			tokenWithSuffix := strings.TrimPrefix(out.String(), "FULLNODE_API_INFO=")
			authToken := strings.Split(tokenWithSuffix, ":/")

			// instantiate the client
			client, err := filecoin.NewClient(
				filecoin.DefaultClientOptions().
					WithRPCURL("ws://127.0.0.1:1234/rpc/v0").
					WithAuthToken(fmt.Sprintf("Bearer %s", authToken[0])),
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

			// construct the transaction builder
			gasPremium := pack.NewU256FromU64(pack.NewU64(149514))
			filTxBuilder := filecoin.NewTxBuilder(gasPremium)

			// build the transaction
			tx, err := filTxBuilder.BuildTx(
				multichain.Address(pack.String(senderFilAddr.String())),
				multichain.Address(pack.String(recipientFilAddr.String())),
				pack.NewU256FromU64(pack.NewU64(100000000)), // amount
				pack.NewU256FromU64(pack.NewU64(0)),         // nonce
				pack.NewU256FromU64(pack.NewU64(495335)),    // gasFeeCap
				pack.NewU256FromU64(pack.NewU64(149838)),    // gasPrice
				pack.Bytes(nil),                             // payload
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

			// wait for the transaction to be included in a block
			for {
				time.Sleep(2 * time.Second)
				fetchedTx, confs, err := client.Tx(ctx, txHash)
				if fetchedTx != nil {
					Expect(err).ToNot(HaveOccurred())
					Expect(confs).To(BeNumerically(">=", 0))
					break
				}
			}
		})
	})
})

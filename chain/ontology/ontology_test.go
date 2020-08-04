package ontology_test

import (
	"log"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-crypto/signature"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/service/native/ont"
	"github.com/renproject/multichain/chain/ontology"
)

var _ = Describe("Ontology", func() {
	Context("when submitting transactions", func() {
		Context("when sending ONT to multiple addresses", func() {
			It("should work", func() {
				pkEnv := os.Getenv("ONTOLOGY_PK")
				if pkEnv == "" {
					panic("ONTOLOGY_PK is undefined")
				}
				pri, err := keypair.WIF2Key([]byte(pkEnv))
				Expect(err).ToNot(HaveOccurred())
				acc := &sdk.Account{
					PrivateKey: pri,
					PublicKey:  pri.Public(),
					Address:    types.AddressFromPubKey(pri.Public()),
					SigScheme:  signature.SHA256withECDSA,
				}

				// Setup client
				client := ontology.NewClient("http://127.0.0.1:20336")

				// Gen addresses
				address1 := client.GenAddress()
				address2 := client.GenAddress()

				// Build transfer states
				transfer1 := &ont.State{
					From:  acc.Address,
					To:    address1,
					Value: 1,
				}
				transfer2 := &ont.State{
					From:  acc.Address,
					To:    address2,
					Value: 1,
				}
				states := []*ont.State{transfer1, transfer2}

				// Send the transfer ont to multiple addresses transaction
				txHash, err := client.MultiTransferOnt(0, 20000, acc, states, acc)
				Expect(err).ToNot(HaveOccurred())
				log.Printf("TxHash               %v", txHash.ToHexString())

				// We wait for 1000 ms before beginning to check transaction.
				time.Sleep(1 * time.Second)

				for {
					evt, err := client.GetEvent(txHash.ToHexString())
					Expect(err).ToNot(HaveOccurred())
					if evt != nil {
						Expect(evt.State).Should(BeEquivalentTo(1))
						log.Printf("Tx success       %v", txHash.ToHexString())
						break
					}
					time.Sleep(1 * time.Second)
				}
			})
		})
	})
})

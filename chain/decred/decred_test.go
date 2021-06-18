package decred_test

import (
	"context"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/decred/dcrd/chaincfg/v3"
	"github.com/decred/dcrd/dcrec"
	"github.com/decred/dcrd/dcrec/secp256k1/v3"
	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/renproject/id"
	"github.com/renproject/pack"

	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/utxo"
	"github.com/renproject/multichain/chain/decred"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decred", func() {
	Context("when submitting transactions", func() {
		Context("when sending DCR to multiple addresses", func() {
			It("should work", func() {
				// Load private key, and assume that the associated address has
				// funds to spend. You can do this by setting DECRED_PK to the
				// value specified in the `./multichaindeploy/.env` file.
				pkEnv := os.Getenv("DECRED_PK")
				if pkEnv == "" {
					panic("DECRED_PK is undefined")
				}
				simNetPrivKeyID := [2]byte{0x23, 0x07}
				wif, err := dcrutil.DecodeWIF(pkEnv, simNetPrivKeyID)
				Expect(err).ToNot(HaveOccurred())

				// PKH
				addrPubKeyHash, err := dcrutil.NewAddressPubKeyHash(dcrutil.Hash160(wif.PubKey()), chaincfg.SimNetParams(), dcrec.STEcdsaSecp256k1)
				Expect(err).ToNot(HaveOccurred())

				log.Printf("PKH                %v", addrPubKeyHash.String())

				// Setup the client and load the unspent transaction outputs.
				client := decred.NewClient(decred.DefaultClientOptions())
				outputs, err := client.UnspentOutputs(context.Background(), 0, 999999999, address.Address(addrPubKeyHash.String()))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(outputs)).To(BeNumerically(">", 0))
				output := outputs[0]

				// Check that we can load the output and that it is equal.
				// Otherwise, something strange is happening with the RPC
				// client.
				output2, _, err := client.Output(context.Background(), output.Outpoint)
				Expect(err).ToNot(HaveOccurred())
				Expect(reflect.DeepEqual(output, output2)).To(BeTrue())
				output2, _, err = client.UnspentOutput(context.Background(), output.Outpoint)
				Expect(err).ToNot(HaveOccurred())
				Expect(reflect.DeepEqual(output, output2)).To(BeTrue())

				// Build the transaction by consuming the outputs and spending
				// them to a set of recipients.
				inputs := []utxo.Input{
					{Output: utxo.Output{
						Outpoint: utxo.Outpoint{
							Hash:  output.Outpoint.Hash[:],
							Index: output.Outpoint.Index,
						},
						PubKeyScript: output.PubKeyScript,
						Value:        output.Value,
					}},
				}
				recipients := []utxo.Recipient{
					{
						To:    address.Address(addrPubKeyHash.String()),
						Value: pack.NewU256FromU64(pack.NewU64((output.Value.Int().Uint64() - 1000) / 3)),
					},
				}
				tx, err := decred.NewTxBuilder(chaincfg.SimNetParams()).BuildTx(inputs, recipients)
				Expect(err).ToNot(HaveOccurred())

				// Get the digests that need signing from the transaction, and
				// sign them. In production, this would be done using the RZL
				// MPC algorithm, but for the purposes of this test, using an
				// explicit privkey is ok.
				sighashes, err := tx.Sighashes()
				signatures := make([]pack.Bytes65, len(sighashes))
				Expect(err).ToNot(HaveOccurred())
				for i := range sighashes {
					hash := id.Hash(sighashes[i])
					priv := secp256k1.PrivKeyFromBytes(wif.PrivKey())
					privk := priv.ToECDSA()
					privKey := (*id.PrivKey)(privk)
					signature, err := privKey.Sign(&hash)
					Expect(err).ToNot(HaveOccurred())
					signatures[i] = pack.NewBytes65(signature)
				}
				Expect(tx.Sign(signatures, pack.NewBytes(wif.PubKey()))).To(Succeed())

				// Submit the transaction to the dcrd node. Again, this
				// should be running a la `./multichaindeploy`.
				txHash, err := tx.Hash()
				Expect(err).ToNot(HaveOccurred())
				err = client.SubmitTx(context.Background(), tx)
				Expect(err).ToNot(HaveOccurred())
				log.Printf("TXID               %v", txHash)

				for {
					// Loop until the transaction has at least a few
					// confirmations. This implies that the transaction is
					// definitely valid, and the test has passed. We were
					// successfully able to use the multichain to construct and
					// submit a Bitcoin transaction!
					confs, err := client.Confirmations(context.Background(), txHash)
					Expect(err).ToNot(HaveOccurred())
					log.Printf("                   %v/3 confirmations", confs)
					if confs >= 1 {
						break
					}
					time.Sleep(10 * time.Second)
				}

				// Check that we can load the output and that it is equal.
				// Otherwise, something strange is happening with the RPC
				// client.
				output2, _, err = client.Output(context.Background(), output.Outpoint)
				Expect(err).ToNot(HaveOccurred())
				Expect(reflect.DeepEqual(output, output2)).To(BeTrue())
			})
		})
	})
})

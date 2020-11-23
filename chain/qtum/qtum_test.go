package qtum_test

// DEZU: This is a straight rip of dogecoin's implementation

import (
	"context"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/renproject/id"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/utxo"
	"github.com/renproject/multichain/chain/qtum"
	"github.com/renproject/pack"
	//"github.com/qtumproject/qtumsuite/chaincfg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

	// DEZU: values take from from qtumsuit
var RegressionNetParams = chaincfg.Params{
	Name: "regtest",
	DefaultPort: "23888",

	Net: 0xe1c6ddfd,

	// Address encoding magics
	PubKeyHashAddrID: 120, // starts with m or n
	ScriptHashAddrID: 110, // starts with 2
	PrivateKeyID:     239, // starts with 9 (uncompressed) or c (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub

	// Human-readable part for Bech32 encoded segwit addresses, as defined in BIP 173.
	Bech32HRPSegwit: "qcrt",
}

var _ = Describe("Qtum", func() {
	Context("when submitting transactions", func() {
		Context("when sending QTUM to multiple addresses", func() {
			It("should work", func() {
				// Load private key, and assume that the associated address has
				// funds to spend. You can do this by setting QTUM_PK to the
				// value specified in the `./multichaindeploy/.env` file.
				pkEnv := os.Getenv("QTUM_PK")
				if pkEnv == "" {
					panic("QTUM_PK is undefined")
				}
				wif, err := btcutil.DecodeWIF(pkEnv)
				Expect(err).ToNot(HaveOccurred())

				// PKH
				pkhAddr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.PrivKey.PubKey().SerializeCompressed()), &RegressionNetParams)
				Expect(err).ToNot(HaveOccurred())
				pkhAddrUncompressed, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.PrivKey.PubKey().SerializeUncompressed()), &RegressionNetParams)
				Expect(err).ToNot(HaveOccurred())
				log.Printf("PKH                %v", pkhAddr.EncodeAddress())
				log.Printf("PKH (uncompressed) %v", pkhAddrUncompressed.EncodeAddress())

				// Setup the client and load the unspent transaction outputs.
				client := qtum.NewClient(qtum.DefaultClientOptions().WithHost("http://127.0.0.1:13889")) // DEZU: TODO: This is QTUM's testnet address, is that the one?
				outputs, err := client.UnspentOutputs(context.Background(), 0, 999999999, address.Address(pkhAddr.EncodeAddress()))
				Expect(err).ToNot(HaveOccurred())
				Expect(len(outputs)).To(BeNumerically(">", 0))
				output := outputs[0]

				// Check that we can load the output and that it is equal.
				// Otherwise, something strange is happening with the RPC
				// client.
				output2, _, err := client.Output(context.Background(), output.Outpoint)
				Expect(err).ToNot(HaveOccurred())
				Expect(reflect.DeepEqual(output, output2)).To(BeTrue())

				// Build the transaction by consuming the outputs and spending
				// them to a set of recipients.
				inputs := []utxo.Input{
					{Output: output},
				}
				recipients := []utxo.Recipient{
					{
						To:    address.Address(pkhAddr.EncodeAddress()),
						Value: pack.NewU256FromU64(pack.NewU64((output.Value.Int().Uint64() - 100000) / 2)), // DEZU: the constant fee here will cause error if too low, demand of > 90000 seems common
					},
					{
						To:    address.Address(pkhAddrUncompressed.EncodeAddress()),
						Value: pack.NewU256FromU64(pack.NewU64((output.Value.Int().Uint64() - 100000) / 2)), // DEZU: Ditto above
					},
				}
				tx, err := qtum.NewTxBuilder(&RegressionNetParams).BuildTx(inputs, recipients)
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
					privKey := (*id.PrivKey)(wif.PrivKey)
					signature, err := privKey.Sign(&hash)
					Expect(err).ToNot(HaveOccurred())
					signatures[i] = pack.NewBytes65(signature)
				}
				Expect(tx.Sign(signatures, pack.NewBytes(wif.SerializePubKey()))).To(Succeed())

				// Submit the transaction to the Qtum node. Again, this
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
					// submit a Qtum transaction!
					confs, err := client.Confirmations(context.Background(), txHash)
					Expect(err).ToNot(HaveOccurred())
					log.Printf("                   %v/3 confirmations", confs)
					if confs >= 3 {
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

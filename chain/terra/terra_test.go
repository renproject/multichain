package terra_test

import (
	"context"
	"encoding/hex"
	"os"
	"time"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/app"

	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/cosmos"
	"github.com/renproject/multichain/chain/terra"
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Terra", func() {
	Context("when submitting transactions", func() {
		Context("when sending LUNA", func() {
			It("should work", func() {
				// create context for the test
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				// Load private key, and assume that the associated address has
				// funds to spend. You can do this by setting TERRA_PK to the
				// value specified in the `./multichaindeploy/.env` file.
				pkEnv := os.Getenv("TERRA_PK")
				if pkEnv == "" {
					panic("TERRA_PK is undefined")
				}

				addrEnv := os.Getenv("TERRA_ADDRESS")
				if addrEnv == "" {
					panic("TERRA_ADDRESS is undefined")
				}

				pkBz, err := hex.DecodeString(pkEnv)
				Expect(err).ToNot(HaveOccurred())

				var pk secp256k1.PrivKeySecp256k1
				copy(pk[:], pkBz)

				addr := terra.Address(pk.PubKey().Address())

				decoder := terra.NewAddressDecoder("terra")
				_, err = decoder.DecodeAddress(multichain.Address(pack.NewString(addrEnv)))
				Expect(err).ToNot(HaveOccurred())

				// random recipient
				pkRecipient := secp256k1.GenPrivKey()
				recipient := sdk.AccAddress(pkRecipient.PubKey().Address())

				// instantiate a new client
				client := terra.NewClient(cosmos.DefaultClientOptions())

				// create a new cosmos-compatible transaction builder
				txBuilder := terra.NewTxBuilder(terra.TxBuilderOptions{
					AccountNumber: pack.NewU64(1),
					ChainID:       "testnet",
					CoinDenom:     "uluna",
					Cdc:           app.MakeCodec(),
				})

				// build the transaction
				payload := pack.NewBytes([]byte("multichain"))
				tx, err := txBuilder.BuildTx(
					multichain.Address(addr.String()),      // from
					multichain.Address(recipient.String()), // to
					pack.NewU256FromU64(pack.U64(2000000)), // amount
					pack.NewU256FromU64(0),                 // nonce
					pack.NewU256FromU64(pack.U64(300000)),  // gas
					pack.NewU256FromU64(pack.U64(300)),     // fee
					payload,                                // memo
				)
				Expect(err).NotTo(HaveOccurred())

				// get the transaction bytes and sign it
				sighashes, err := tx.Sighashes()
				Expect(err).NotTo(HaveOccurred())
				sigBytes, err := pk.Sign([]byte(sighashes[0]))
				Expect(err).NotTo(HaveOccurred())

				// attach the signature to the transaction
				pubKey := pk.PubKey().(secp256k1.PubKeySecp256k1)
				err = tx.Sign(
					[]pack.Bytes{pack.NewBytes(sigBytes)},
					pack.NewBytes(pubKey[:]),
				)
				Expect(err).NotTo(HaveOccurred())

				// submit the transaction to the chain
				txHash := tx.Hash()
				err = client.SubmitTx(ctx, tx)
				Expect(err).NotTo(HaveOccurred())

				for {
					// Loop until the transaction has at least a few
					// confirmations. This implies that the transaction is
					// definitely valid, and the test has passed. We were
					// successfully able to use the multichain to construct and
					// submit a Bitcoin transaction!
					foundTx, confs, err := client.Tx(ctx, txHash)
					if err == nil {
						Expect(confs.Uint64()).To(Equal(uint64(1)))
						Expect(foundTx.Payload()).To(Equal(multichain.ContractCallData([]byte(payload.String()))))
						break
					}

					// wait and retry querying for the transaction
					time.Sleep(2 * time.Second)
				}
			})
		})
	})
})

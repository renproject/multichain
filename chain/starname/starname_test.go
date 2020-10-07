package starname_test

import (
	"context"
	"encoding/hex"
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/app"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/cosmos"
	"github.com/renproject/multichain/chain/starname"
	"github.com/renproject/pack"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Starname (IOV)", func() {
	Context("when submitting transactions", func() {
		Context("when sending IOV", func() {
			It("should work", func() {
				// create context for the test
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				// Load private key, and assume that the associated address has
				// funds to spend. You can do this by setting IOV_PK to the
				// value specified in the `./multichaindeploy/.env` file.
				pkEnv := os.Getenv("IOV_PK")
				if pkEnv == "" {
					panic("IOV_PK is undefined")
				}

				addrEnv := os.Getenv("IOV_ADDRESS")
				if addrEnv == "" {
					panic("IOV_ADDRESS is undefined")
				}

				pkBz, err := hex.DecodeString(pkEnv)
				Expect(err).ToNot(HaveOccurred())

				var pk secp256k1.PrivKeySecp256k1
				copy(pk[:], pkBz)

				addr := starname.Address(pk.PubKey().Address())

				decoder := starname.NewAddressDecoder("star")
				_, err = decoder.DecodeAddress(multichain.Address(pack.NewString(addrEnv)))
				Expect(err).ToNot(HaveOccurred())

				// random recipient
				pkRecipient := secp256k1.GenPrivKey()
				recipient := types.AccAddress(pkRecipient.PubKey().Address())

				// instantiate a new client
				client := starname.NewClient(cosmos.DefaultClientOptions())

				// create a new cosmos-compatible transaction builder
				txBuilder := starname.NewTxBuilder(starname.TxBuilderOptions{
					AccountNumber: pack.NewU64(1),
					ChainID:       "testnet",
					CoinDenom:     "tiov",
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
				sigBytes, err := pk.Sign(sighashes[0][:])
				Expect(err).NotTo(HaveOccurred())
				sig65 := pack.Bytes65{}
				copy(sig65[:], sigBytes)

				// attach the signature to the transaction
				pubKey := pk.PubKey().(secp256k1.PubKeySecp256k1)
				err = tx.Sign(
					[]pack.Bytes65{sig65},
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
					// submit a Starname (IOV) transaction!
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

package starname_test

import (
	"context"
	"encoding/hex"
	"os"

	"github.com/renproject/id"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/chain/starname"
	"github.com/renproject/pack"
	"github.com/renproject/surge"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Starname", func() {
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

				pkBz, err := hex.DecodeString(pkEnv)
				Expect(err).ToNot(HaveOccurred())

				var pk secp256k1.PrivKeySecp256k1
				copy(pk[:], pkBz)

				var privKey id.PrivKey
				err = surge.FromBinary(&privKey, pkBz)
				Expect(err).NotTo(HaveOccurred())

				addr := starname.Address(pk.PubKey().Address())

				// random recipient
				pkRecipient := secp256k1.GenPrivKey()
				addrEncoder := starname.NewAddressEncoder("star")
				recipient, err := addrEncoder.EncodeAddress(address.RawAddress(pack.Bytes(pkRecipient.PubKey().Address())))
				Expect(err).NotTo(HaveOccurred())

				// avoid a port collision with terra and set BroadcastMode to block to avoid having to loop on client.Tx()
				options := starname.DefaultClientOptions()
				options.Host = "http://0.0.0.0:46657"
				options.BroadcastMode = "block"

				// instantiate a new client
				client := starname.NewClient(options)
				nonce, err := client.AccountNonce(ctx, multichain.Address(addr.String()))
				Expect(err).NotTo(HaveOccurred())

				// create a new cosmos-compatible transaction builder
				txBuilder := starname.NewTxBuilder(starname.TxBuilderOptions{
					ChainID:   "testnet",
					CoinDenom: "tiov",
				}, client)

				// build the transaction
				payload := pack.NewBytes([]byte("multichain"))
				amount := pack.NewU256FromU64(pack.U64(2000000))
				tx, err := txBuilder.BuildTx(
					ctx,
					multichain.Address(addr.String()),     // from
					recipient,                             // to
					amount,                                // amount
					nonce,                                 // nonce
					pack.NewU256FromU64(pack.U64(200000)), // gas limit
					pack.NewU256FromU64(pack.U64(1)),      // gas price
					pack.NewU256FromU64(pack.U64(1)),      // gas cap
					payload,                               // memo
				)
				Expect(err).NotTo(HaveOccurred())

				// get the transaction bytes and sign it
				sighashes, err := tx.Sighashes()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(sighashes)).To(Equal(1))
				hash := id.Hash(sighashes[0])
				sig, err := privKey.Sign(&hash)
				Expect(err).NotTo(HaveOccurred())
				sigBytes, err := surge.ToBinary(sig)
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

				// We don't need to loop due to instant finality and broadcast
				// mode "block".  If err == nil then we were successfully able to
				// use the multichain to construct and submit a Starname (IOV)
				// transaction!
				foundTx, confs, err := client.Tx(ctx, txHash)
				Expect(err).NotTo(HaveOccurred())
				Expect(confs.Uint64()).To(Equal(uint64(1)))
				Expect(foundTx.Payload()).To(Equal(multichain.ContractCallData([]byte(payload.String()))))
				Expect(foundTx.From()).To(Equal(multichain.Address(addr.String())))
				Expect(foundTx.To()).To(Equal(recipient))
				Expect(foundTx.Value()).To(Equal(amount))
			})
		})
	})
})

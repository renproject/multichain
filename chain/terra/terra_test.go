package terra_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
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
		Context("when sending LUNA to multiple addresses", func() {
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

				// pkEnv := "a96e62ed3955e65be32703f12d87b6b5cf26039ecfa948dc5107a495418e5330"
				// addrEnv := "terra10s4mg25tu6termrk8egltfyme4q7sg3hl8s38u"

				pkBz, err := hex.DecodeString(pkEnv)
				Expect(err).ToNot(HaveOccurred())

				var pk secp256k1.PrivKeySecp256k1
				copy(pk[:], pkBz)

				addr := cosmos.Address(pk.PubKey().Address())

				decoder := terra.NewAddressDecoder("terra")
				expectedAddr, err := decoder.DecodeAddress(multichain.Address(pack.NewString(addrEnv)))
				Expect(err).ToNot(HaveOccurred())
				Expect(addr).Should(Equal(expectedAddr))

				pk1 := secp256k1.GenPrivKey()
				// pk2 := secp256k1.GenPrivKey()

				recipient1 := sdk.AccAddress(pk1.PubKey().Address())
				// recipient2 := sdk.AccAddress(pk2.PubKey().Address())

				// msgs := []cosmos.MsgSend{
				// 	{
				// 		FromAddress: cosmos.Address(addr),
				// 		ToAddress:   cosmos.Address(recipient1),
				// 		Amount: cosmos.Coins{
				// 			{
				// 				Denom:  "uluna",
				// 				Amount: pack.U64(1000000),
				// 			},
				// 		},
				// 	},
				// 	{
				// 		FromAddress: cosmos.Address(addr),
				// 		ToAddress:   cosmos.Address(recipient1),
				// 		Amount: cosmos.Coins{
				// 			{
				// 				Denom:  "uluna",
				// 				Amount: pack.U64(2000000),
				// 			},
				// 		},
				// 	},
				// }

				client := cosmos.NewClient(cosmos.DefaultClientOptions(), app.MakeCodec())
				account, err := client.Account(addr)
				Expect(err).NotTo(HaveOccurred())

				fmt.Printf("account = %v\n", account)

				txBuilder := terra.NewTxBuilder(cosmos.TxBuilderOptions{
					AccountNumber: account.AccountNumber,
					// SequenceNumber: account.SequenceNumber,
					// Gas:            200000,
					ChainID:   "testnet",
					CoinDenom: "uluna",
					Cdc:       app.MakeCodec(),
					// Memo:           "multichain",
					// Fees: cosmos.Coins{
					// 	{
					// 		Denom:  "uluna",
					// 		Amount: pack.U64(3000),
					// 	},
					// },
				})

				tx, err := txBuilder.BuildTx(
					multichain.Address(recipient1.String()),
					multichain.Address(addr.String()),
					pack.NewU256FromU64(pack.U64(2000000)), // value
					pack.NewU256FromU64(account.SequenceNumber),
					pack.NewU256FromU64(pack.U64(20000)), // gas limit
					pack.NewU256FromU64(pack.U64(300)),   // gas price,
					pack.NewBytes([]byte("multichain")),  // memo
				)
				Expect(err).NotTo(HaveOccurred())

				sighashes, err := tx.Sighashes()
				Expect(err).NotTo(HaveOccurred())
				sigBytes, err := pk.Sign([]byte(sighashes[0]))
				Expect(err).NotTo(HaveOccurred())

				pubKey := pk.PubKey().(secp256k1.PubKeySecp256k1)
				err = tx.Sign(
					[]pack.Bytes{pack.NewBytes(sigBytes)},
					pack.NewBytes(pubKey[:]),
				)
				Expect(err).NotTo(HaveOccurred())

				txHash := tx.Hash()
				err = client.SubmitTx(ctx, tx)
				Expect(err).NotTo(HaveOccurred())

				for {
					// Loop until the transaction has at least a few
					// confirmations. This implies that the transaction is
					// definitely valid, and the test has passed. We were
					// successfully able to use the multichain to construct and
					// submit a Bitcoin transaction!
					_, confs, err := client.Tx(ctx, txHash)
					if err == nil {
						break
					}

					if !strings.Contains(err.Error(), "not found") {
						Expect(err).NotTo(HaveOccurred())
					}

					Expect(confs).To(Equal(1))

					time.Sleep(2 * time.Second)
				}
			})
		})
	})
})

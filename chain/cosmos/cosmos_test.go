package cosmos_test

// import (
// 	"encoding/hex"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/tendermint/tendermint/crypto/secp256k1"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/terra-project/core/app"

// 	"github.com/renproject/multichain/chain/cosmos"
// 	"github.com/renproject/multichain/compat/cosmoscompat"
// 	"github.com/renproject/pack"

// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// )

// var _ = Describe("Cosmos", func() {
// 	Context("when submitting transactions", func() {
// 		Context("when sending LUNA to multiple addresses", func() {
// 			It("should work", func() {
// 				// Load private key, and assume that the associated address has
// 				// funds to spend. You can do this by setting TERRA_PK to the
// 				// value specified in the `./multichaindeploy/.env` file.
// 				pkEnv := os.Getenv("TERRA_PK")
// 				if pkEnv == "" {
// 					panic("TERRA_PK is undefined")
// 				}

// 				addrEnv := os.Getenv("TERRA_ADDRESS")
// 				if addrEnv == "" {
// 					panic("TERRA_ADDRESS is undefined")
// 				}

// 				// pkEnv := "a96e62ed3955e65be32703f12d87b6b5cf26039ecfa948dc5107a495418e5330"
// 				// addrEnv := "terra10s4mg25tu6termrk8egltfyme4q7sg3hl8s38u"

// 				pkBz, err := hex.DecodeString(pkEnv)
// 				Expect(err).ToNot(HaveOccurred())

// 				var pk secp256k1.PrivKeySecp256k1
// 				copy(pk[:], pkBz)

// 				addr := cosmoscompat.Address(pk.PubKey().Address())

// 				decoder := cosmos.NewAddressDecoder("terra")
// 				expectedAddr, err := decoder.DecodeAddress(pack.NewString(addrEnv))
// 				Expect(err).ToNot(HaveOccurred())
// 				Expect(addr).Should(Equal(expectedAddr))

// 				pk1 := secp256k1.GenPrivKey()
// 				pk2 := secp256k1.GenPrivKey()

// 				recipient1 := sdk.AccAddress(pk1.PubKey().Address())
// 				recipient2 := sdk.AccAddress(pk2.PubKey().Address())

// 				msgs := []cosmoscompat.MsgSend{
// 					{
// 						FromAddress: cosmoscompat.Address(addr),
// 						ToAddress:   cosmoscompat.Address(recipient1),
// 						Amount: cosmoscompat.Coins{
// 							{
// 								Denom:  "uluna",
// 								Amount: pack.U64(1000000),
// 							},
// 						},
// 					},
// 					{
// 						FromAddress: cosmoscompat.Address(addr),
// 						ToAddress:   cosmoscompat.Address(recipient2),
// 						Amount: cosmoscompat.Coins{
// 							{
// 								Denom:  "uluna",
// 								Amount: pack.U64(2000000),
// 							},
// 						},
// 					},
// 				}

// 				client := cosmoscompat.NewClient(cosmoscompat.DefaultClientOptions(), app.MakeCodec())
// 				account, err := client.Account(addr)
// 				Expect(err).NotTo(HaveOccurred())

// 				txBuilder := cosmos.NewTxBuilder(cosmoscompat.TxOptions{
// 					AccountNumber:  account.AccountNumber,
// 					SequenceNumber: account.SequenceNumber,
// 					Gas:            200000,
// 					ChainID:        "testnet",
// 					Memo:           "multichain",
// 					Fees: cosmoscompat.Coins{
// 						{
// 							Denom:  "uluna",
// 							Amount: pack.U64(3000),
// 						},
// 					},
// 				}).WithCodec(app.MakeCodec())

// 				tx, err := txBuilder.BuildTx(msgs)
// 				Expect(err).NotTo(HaveOccurred())

// 				sigBytes, err := pk.Sign(tx.SigBytes())
// 				Expect(err).NotTo(HaveOccurred())

// 				pubKey := pk.PubKey().(secp256k1.PubKeySecp256k1)
// 				err = tx.Sign([]cosmoscompat.StdSignature{
// 					{
// 						Signature: pack.NewBytes(sigBytes),
// 						PubKey:    pack.NewBytes(pubKey[:]),
// 					},
// 				})
// 				Expect(err).NotTo(HaveOccurred())

// 				txHash, err := client.SubmitTx(tx, pack.NewString("sync"))
// 				Expect(err).NotTo(HaveOccurred())

// 				for {
// 					// Loop until the transaction has at least a few
// 					// confirmations. This implies that the transaction is
// 					// definitely valid, and the test has passed. We were
// 					// successfully able to use the multichain to construct and
// 					// submit a Bitcoin transaction!
// 					_, err := client.Tx(txHash)
// 					if err == nil {
// 						break
// 					}

// 					if !strings.Contains(err.Error(), "not found") {
// 						Expect(err).NotTo(HaveOccurred())
// 					}

// 					time.Sleep(10 * time.Second)
// 				}
// 			})
// 		})
// 	})
// })

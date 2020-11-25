package icon_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	secphal "github.com/haltingstate/secp256k1-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/id"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/icon"
	"github.com/renproject/multichain/chain/icon/crypto"
	"github.com/renproject/pack"
	"github.com/renproject/surge"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ = Describe("Multichain", func() {
	// Create context to work within.
	ctx := context.Background()

	// Initialise the logger.
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := loggerConfig.Build()
	Expect(err).ToNot(HaveOccurred())

	//
	// ACCOUNT API
	//
	Context("Account API", func() {
		accountChainTable := []struct {
			senderEnv           func() (id.PrivKey, *id.PubKey, multichain.Address)
			privKeyToAddr       func(pk id.PrivKey) multichain.Address
			rpcURL              pack.String
			randomRecipientAddr func() multichain.Address
			initialise          func(pack.String) (multichain.AccountClient, multichain.AccountTxBuilder)
			txParams            func(multichain.AccountClient) (pack.U256, pack.U256, pack.U256, pack.U256, pack.Bytes)
			chain               multichain.Chain
		}{
			{
				func() (id.PrivKey, *id.PubKey, multichain.Address) {
					pkEnv := os.Getenv("ICON_PK")
					if pkEnv == "" {
						panic("ICON_PK is undefined")
					}
					pkDecode, err := hex.DecodeString(pkEnv)
					privKey, err := crypto.ParsePrivateKey(pkDecode)
					Expect(err).ToNot(HaveOccurred())
					addrEncoder := icon.NewAddressEncodeDecoder()
					senderAddress := icon.NewAccountAddressFromPublicKey(privKey.PublicKey())
					senderAddr, err := addrEncoder.EncodeAddress(senderAddress.Bytes())
					Expect(err).ToNot(HaveOccurred())
					senderPrivKey := id.PrivKey{}
					err = surge.FromBinary(&senderPrivKey, pkDecode)
					return senderPrivKey, senderPrivKey.PubKey(), senderAddr
				},
				func(privKey id.PrivKey) multichain.Address {
					pkBinary, err := surge.ToBinary(privKey)
					Expect(err).NotTo(HaveOccurred())
					pkDecode, err := hex.DecodeString(string(pkBinary))
					Expect(err).NotTo(HaveOccurred())
					privKeyTmp, err := crypto.ParsePrivateKey(pkDecode)
					addrEncoder := icon.NewAddressEncodeDecoder()
					tmpAddress := icon.NewAccountAddressFromPublicKey(privKeyTmp.PublicKey())
					addr, err := addrEncoder.EncodeAddress(tmpAddress.Bytes())
					Expect(err).NotTo(HaveOccurred())
					return addr
				},
				"http://127.0.0.1:9000/api/v3",
				func() multichain.Address {
					_, pub_key := crypto.GenerateKeyPair()
					recipient := icon.NewAccountAddressFromPublicKey(pub_key)
					addrEncoder := icon.NewAddressEncodeDecoder()
					recipientAddr, err := addrEncoder.EncodeAddress(recipient.Bytes())
					Expect(err).NotTo(HaveOccurred())
					return recipientAddr
				},
				func(rpcURL pack.String) (multichain.AccountClient, multichain.AccountTxBuilder) {
					client := icon.NewClient(rpcURL)
					txBuilder := icon.NewTxBuilder("0x3") //testnet
					return client, txBuilder
				},
				func(_ multichain.AccountClient) (pack.U256, pack.U256, pack.U256, pack.U256, pack.Bytes) {
					amount := pack.NewU256FromU64(pack.U64(0xde0b6b3a7640000))
					stepLimit := pack.NewU256FromU64(pack.U64(0x2fb60ca5))
					return amount, stepLimit, pack.NewU256FromU64(pack.U64(0)), pack.NewU256FromU64(pack.U64(0)), pack.Bytes(nil)
				},
				multichain.Icon,
			},
		}

		for _, accountChain := range accountChainTable {
			accountChain := accountChain
			Context(fmt.Sprintf("%v", accountChain.chain), func() {
				Specify("build, broadcast and fetch tx", func() {
					// Load private key and the associated address.
					senderPrivKey, _, senderAddr := accountChain.senderEnv()
					// Get a random recipient address.
					recipientAddr := accountChain.randomRecipientAddr()

					// Initialise the account chain's client, and possibly get a nonce for
					// the sender.
					accountClient, txBuilder := accountChain.initialise(accountChain.rpcURL)

					amount, stepLimit, _, _, payload := accountChain.txParams(accountClient)

					// Get the appropriate nonce for sender.
					// Build a transaction.
					accountTx, err := txBuilder.BuildTx(
						senderAddr,
						recipientAddr,
						amount, // amount,
						pack.NewU256FromU64(0),
						stepLimit,
						pack.NewU256FromU64(0),
						payload, // payload
					)
					Expect(err).NotTo(HaveOccurred())

					sighashes, err := accountTx.Sighashes()
					pkBinary, err := surge.ToBinary(senderPrivKey)
					fmt.Println(pkBinary)
					Expect(err).NotTo(HaveOccurred())
					pkEnv := os.Getenv("ICON_PK")
					if pkEnv == "" {
						panic("ICON_PK is undefined")
					}
					pkDecode, err := hex.DecodeString(pkEnv)
					sig, _ := crypto.ParseSignature(secphal.Sign(sighashes[0].Bytes(), pkDecode))
					signature, err := sig.SerializeRSV()
					sig65 := pack.Bytes65{}
					copy(sig65[:], signature)
					err = accountTx.Sign(
						[]pack.Bytes65{sig65},
						pack.NewBytes(nil),
					)
					Expect(err).NotTo(HaveOccurred())

					err = accountClient.SubmitTx(ctx, accountTx)
					Expect(err).NotTo(HaveOccurred())
					txHash := accountTx.Hash()
					logger.Debug("submit tx", zap.String("from", string(senderAddr)), zap.String("to", string(recipientAddr)), zap.Any("txHash", txHash))
					// Wait slightly before we query the chain's node.
					time.Sleep(time.Second)
					for {
						// Loop until the transaction has at least a few confirmations.
						tx, confs, err := accountClient.Tx(ctx, pack.Bytes("0x"+hex.EncodeToString(txHash)))
						if err == nil {
							Expect(confs.Uint64()).To(BeNumerically(">", 0))
							Expect(tx.From()).To(Equal(senderAddr))
							Expect(tx.To()).To(Equal(recipientAddr))
							Expect(tx.Value().String()).To(Equal("1000000000000000000"))
							break
						}

						// wait and retry querying for the transaction
						time.Sleep(5 * time.Second)
					}
				})
			})
		}
	})
})

package multichain_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/renproject/id"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/bitcoin"
	"github.com/renproject/multichain/chain/bitcoincash"
	"github.com/renproject/multichain/chain/digibyte"
	"github.com/renproject/multichain/chain/dogecoin"
	"github.com/renproject/multichain/chain/terra"
	"github.com/renproject/multichain/chain/zcash"
	"github.com/renproject/pack"
	"github.com/renproject/surge"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/terra-project/core/app"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Multichain", func() {
	// Create context to work within
	ctx := context.Background()

	// Initialise the logger
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := loggerConfig.Build()
	Expect(err).ToNot(HaveOccurred())

	Context("Address API", func() {
		It("should pass", func() {
			Fail("not implemented")
		})
	})

	Context("Account API", func() {
		accountChainTable := []struct {
			senderEnv           func() (id.PrivKey, *id.PubKey, multichain.Address)
			privKeyToAddr       func(pk id.PrivKey) multichain.Address
			rpcURL              pack.String
			randomRecipientAddr func() multichain.Address
			initialise          func() multichain.AccountClient
			txBuilder           multichain.AccountTxBuilder
			txParams            func() (pack.U256, pack.U256, pack.U256, pack.U256, pack.Bytes)
			chain               multichain.Chain
		}{
			{
				func() (id.PrivKey, *id.PubKey, multichain.Address) {
					pkEnv := os.Getenv("TERRA_PK")
					if pkEnv == "" {
						panic("TERRA_PK is undefined")
					}
					pkBytes, err := hex.DecodeString(pkEnv)
					Expect(err).NotTo(HaveOccurred())
					var pk secp256k1.PrivKeySecp256k1
					copy(pk[:], pkBytes)
					addrEncoder := terra.NewAddressEncoder("terra")
					senderAddr, err := addrEncoder.EncodeAddress(multichain.RawAddress(pack.Bytes(pk.PubKey().Address())))
					Expect(err).NotTo(HaveOccurred())
					senderPrivKey := id.PrivKey{}
					err = surge.FromBinary(&senderPrivKey, pkBytes)
					Expect(err).NotTo(HaveOccurred())
					return senderPrivKey, senderPrivKey.PubKey(), senderAddr
				},
				func(privKey id.PrivKey) multichain.Address {
					pkBytes, err := surge.ToBinary(privKey)
					Expect(err).NotTo(HaveOccurred())
					var pk secp256k1.PrivKeySecp256k1
					copy(pk[:], pkBytes)
					addrEncoder := terra.NewAddressEncoder("terra")
					addr, err := addrEncoder.EncodeAddress(multichain.RawAddress(pack.Bytes(pk.PubKey().Address())))
					Expect(err).NotTo(HaveOccurred())
					return addr
				},
				"http://127.0.0.1:26657",
				func() multichain.Address {
					recipientKey := secp256k1.GenPrivKey()
					addrEncoder := terra.NewAddressEncoder("terra")
					recipient, err := addrEncoder.EncodeAddress(multichain.RawAddress(pack.Bytes(recipientKey.PubKey().Address())))
					Expect(err).NotTo(HaveOccurred())
					return recipient
				},
				func() multichain.AccountClient {
					client := terra.NewClient(terra.DefaultClientOptions())
					return client
				},
				terra.NewTxBuilder(terra.TxBuilderOptions{
					AccountNumber: pack.NewU64(1),
					ChainID:       "testnet",
					CoinDenom:     "uluna",
					Cdc:           app.MakeCodec(),
				}),
				func() (pack.U256, pack.U256, pack.U256, pack.U256, pack.Bytes) {
					amount := pack.NewU256FromU64(pack.U64(2000000))
					nonce := pack.NewU256FromU64(0)
					gasLimit := pack.NewU256FromU64(pack.U64(300000))
					gasPrice := pack.NewU256FromU64(pack.U64(300))
					payload := pack.NewBytes([]byte("multichain"))
					return amount, nonce, gasLimit, gasPrice, payload
				},
				multichain.Terra,
			},
		}

		for _, accountChain := range accountChainTable {
			accountChain := accountChain
			FContext(fmt.Sprintf("%v", accountChain.chain), func() {
				Specify("build, broadcast and fetch tx", func() {
					// Load private key and the associated address.
					senderPrivKey, senderPubKey, senderAddr := accountChain.senderEnv()
					fmt.Printf("sender address   = %v\n", senderAddr)

					// Get a random recipient address.
					recipientAddr := accountChain.randomRecipientAddr()
					fmt.Printf("random recipient = %v\n", recipientAddr)

					// Initialise the account chain's client.
					accountClient := accountChain.initialise()

					// Build a transaction.
					amount, nonce, gasLimit, gasPrice, payload := accountChain.txParams()
					accountTx, err := accountChain.txBuilder.BuildTx(
						multichain.Address(senderAddr),
						recipientAddr,
						amount, nonce, gasLimit, gasPrice,
						payload,
					)
					Expect(err).NotTo(HaveOccurred())

					// Get the transaction bytes and sign them.
					sighashes, err := accountTx.Sighashes()
					Expect(err).NotTo(HaveOccurred())
					hash := id.Hash(sighashes[0])
					sig, err := senderPrivKey.Sign(&hash)
					Expect(err).NotTo(HaveOccurred())
					sigBytes, err := surge.ToBinary(sig)
					Expect(err).NotTo(HaveOccurred())
					txSignature := pack.Bytes65{}
					copy(txSignature[:], sigBytes)
					senderPubKeyBytes, err := surge.ToBinary(senderPubKey)
					Expect(err).NotTo(HaveOccurred())
					err = accountTx.Sign(
						[]pack.Bytes65{txSignature},
						pack.NewBytes(senderPubKeyBytes),
					)
					Expect(err).NotTo(HaveOccurred())
					ser, err := accountTx.Serialize()
					Expect(err).NotTo(HaveOccurred())
					fmt.Printf("tx serialised = %v\n", ser)

					// Submit the transaction to the account chain.
					txHash := accountTx.Hash()
					fmt.Printf("tx hash = %v\n", txHash)
					err = accountClient.SubmitTx(ctx, accountTx)
					Expect(err).NotTo(HaveOccurred())

					// Wait slightly before we query the chain's node.
					time.Sleep(time.Second)

					for {
						// Loop until the transaction has at least a few confirmations.
						tx, confs, err := accountClient.Tx(ctx, txHash)
						if err == nil {
							Expect(confs.Uint64()).To(Equal(uint64(1)))
							Expect(tx.Value()).To(Equal(amount))
							Expect(tx.From()).To(Equal(senderAddr))
							Expect(tx.To()).To(Equal(recipientAddr))
							break
						}

						// wait and retry querying for the transaction
						time.Sleep(5 * time.Second)
					}
				})
			})
		}
	})

	Context("UTXO API", func() {
		utxoChainTable := []struct {
			privKeyEnv    string
			newAddressPKH func([]byte) (btcutil.Address, error)
			newAddressSH  func([]byte) (btcutil.Address, error)
			rpcURL        pack.String
			initialise    func(pack.String, btcutil.Address) (multichain.UTXOClient, []multichain.UTXOutput, func(context.Context, pack.Bytes) (int64, error))
			txBuilder     multichain.UTXOTxBuilder
			chain         multichain.Chain
		}{
			{
				"BITCOIN_PK",
				func(pkh []byte) (btcutil.Address, error) {
					addr, err := btcutil.NewAddressPubKeyHash(pkh, &chaincfg.RegressionNetParams)
					return addr, err
				},
				func(script []byte) (btcutil.Address, error) {
					addr, err := btcutil.NewAddressScriptHash(script, &chaincfg.RegressionNetParams)
					return addr, err
				},
				pack.NewString("http://0.0.0.0:18443"),
				func(rpcURL pack.String, pkhAddr btcutil.Address) (multichain.UTXOClient, []multichain.UTXOutput, func(context.Context, pack.Bytes) (int64, error)) {
					client := bitcoin.NewClient(bitcoin.DefaultClientOptions())
					outputs, err := client.UnspentOutputs(ctx, 0, 999999999, multichain.Address(pkhAddr.EncodeAddress()))
					Expect(err).NotTo(HaveOccurred())
					return client, outputs, client.Confirmations
				},
				bitcoin.NewTxBuilder(&chaincfg.RegressionNetParams),
				multichain.Bitcoin,
			},
			{
				"BITCOINCASH_PK",
				func(pkh []byte) (btcutil.Address, error) {
					addr, err := bitcoincash.NewAddressPubKeyHash(pkh, &chaincfg.RegressionNetParams)
					return addr, err
				},
				func(script []byte) (btcutil.Address, error) {
					addr, err := bitcoincash.NewAddressScriptHash(script, &chaincfg.RegressionNetParams)
					return addr, err
				},
				pack.NewString("http://0.0.0.0:19443"),
				func(rpcURL pack.String, pkhAddr btcutil.Address) (multichain.UTXOClient, []multichain.UTXOutput, func(context.Context, pack.Bytes) (int64, error)) {
					client := bitcoincash.NewClient(bitcoincash.DefaultClientOptions())
					outputs, err := client.UnspentOutputs(ctx, 0, 999999999, multichain.Address(pkhAddr.EncodeAddress()))
					Expect(err).NotTo(HaveOccurred())
					return client, outputs, client.Confirmations
				},
				bitcoincash.NewTxBuilder(&chaincfg.RegressionNetParams),
				multichain.BitcoinCash,
			},
			{
				"DIGIBYTE_PK",
				func(pkh []byte) (btcutil.Address, error) {
					addr, err := btcutil.NewAddressPubKeyHash(pkh, &digibyte.RegressionNetParams)
					return addr, err
				},
				func(script []byte) (btcutil.Address, error) {
					addr, err := btcutil.NewAddressScriptHash(script, &digibyte.RegressionNetParams)
					return addr, err
				},
				pack.NewString("http://0.0.0.0:20443"),
				func(rpcURL pack.String, pkhAddr btcutil.Address) (multichain.UTXOClient, []multichain.UTXOutput, func(context.Context, pack.Bytes) (int64, error)) {
					client := digibyte.NewClient(digibyte.DefaultClientOptions())
					outputs, err := client.UnspentOutputs(ctx, 0, 999999999, multichain.Address(pkhAddr.EncodeAddress()))
					Expect(err).NotTo(HaveOccurred())
					return client, outputs, client.Confirmations
				},
				digibyte.NewTxBuilder(&digibyte.RegressionNetParams),
				multichain.DigiByte,
			},
			{
				"DOGECOIN_PK",
				func(pkh []byte) (btcutil.Address, error) {
					addr, err := btcutil.NewAddressPubKeyHash(pkh, &dogecoin.RegressionNetParams)
					return addr, err
				},
				func(script []byte) (btcutil.Address, error) {
					addr, err := btcutil.NewAddressScriptHash(script, &dogecoin.RegressionNetParams)
					return addr, err
				},
				pack.NewString("http://0.0.0.0:18332"),
				func(rpcURL pack.String, pkhAddr btcutil.Address) (multichain.UTXOClient, []multichain.UTXOutput, func(context.Context, pack.Bytes) (int64, error)) {
					client := dogecoin.NewClient(dogecoin.DefaultClientOptions())
					outputs, err := client.UnspentOutputs(ctx, 0, 999999999, multichain.Address(pkhAddr.EncodeAddress()))
					Expect(err).NotTo(HaveOccurred())
					return client, outputs, client.Confirmations
				},
				dogecoin.NewTxBuilder(&dogecoin.RegressionNetParams),
				multichain.Dogecoin,
			},
			{
				"ZCASH_PK",
				func(pkh []byte) (btcutil.Address, error) {
					addr, err := zcash.NewAddressPubKeyHash(pkh, &zcash.RegressionNetParams)
					return addr, err
				},
				func(script []byte) (btcutil.Address, error) {
					addr, err := zcash.NewAddressScriptHash(script, &zcash.RegressionNetParams)
					return addr, err
				},
				pack.String("http://0.0.0.0:18232"),
				func(rpcURL pack.String, pkhAddr btcutil.Address) (multichain.UTXOClient, []multichain.UTXOutput, func(context.Context, pack.Bytes) (int64, error)) {
					client := zcash.NewClient(zcash.DefaultClientOptions())
					outputs, err := client.UnspentOutputs(ctx, 0, 999999999, multichain.Address(pkhAddr.EncodeAddress()))
					Expect(err).NotTo(HaveOccurred())
					return client, outputs, client.Confirmations
				},
				zcash.NewTxBuilder(&zcash.RegressionNetParams, 1000000),
				multichain.Zcash,
			},
		}

		for _, utxoChain := range utxoChainTable {
			utxoChain := utxoChain
			Context(fmt.Sprintf("%v", utxoChain.chain), func() {
				Specify("build, broadcast and fetch tx", func() {
					// Load private key.
					pkEnv := os.Getenv(utxoChain.privKeyEnv)
					if pkEnv == "" {
						panic(fmt.Sprintf("%v is undefined", utxoChain.privKeyEnv))
					}
					wif, err := btcutil.DecodeWIF(pkEnv)
					Expect(err).NotTo(HaveOccurred())

					// Get the PKH address from the loaded private key.
					pkhAddr, err := utxoChain.newAddressPKH(btcutil.Hash160(wif.PrivKey.PubKey().SerializeCompressed()))
					Expect(err).NotTo(HaveOccurred())

					// Recipient
					pkhAddrUncompressed, err := utxoChain.newAddressPKH(btcutil.Hash160(wif.PrivKey.PubKey().SerializeUncompressed()))
					Expect(err).ToNot(HaveOccurred())

					// Initialise the UTXO client and fetch the unspent outputs. Also get a
					// function to query the number of block confirmations for a transaction.
					utxoClient, unspentOutputs, confsFn := utxoChain.initialise(utxoChain.rpcURL, pkhAddr)
					Expect(len(unspentOutputs)).To(BeNumerically(">", 0))
					output := unspentOutputs[0]

					// Build a transaction
					inputs := []multichain.UTXOInput{
						{Output: multichain.UTXOutput{
							Outpoint: multichain.UTXOutpoint{
								Hash:  output.Outpoint.Hash[:],
								Index: output.Outpoint.Index,
							},
							PubKeyScript: output.PubKeyScript,
							Value:        output.Value,
						}},
					}
					recipients := []multichain.UTXORecipient{
						{
							To:    multichain.Address(pkhAddr.EncodeAddress()),
							Value: pack.NewU256FromU64(pack.NewU64((output.Value.Int().Uint64() - 1000) / 2)),
						},
						{
							To:    multichain.Address(pkhAddrUncompressed.EncodeAddress()),
							Value: pack.NewU256FromU64(pack.NewU64((output.Value.Int().Uint64() - 1000) / 2)),
						},
					}
					utxoTx, err := utxoChain.txBuilder.BuildTx(inputs, recipients)
					Expect(err).NotTo(HaveOccurred())

					// Get the sighashes that need to be signed, and sign them.
					sighashes, err := utxoTx.Sighashes()
					signatures := make([]pack.Bytes65, len(sighashes))
					Expect(err).ToNot(HaveOccurred())
					for i := range sighashes {
						hash := id.Hash(sighashes[i])
						privKey := (*id.PrivKey)(wif.PrivKey)
						signature, err := privKey.Sign(&hash)
						Expect(err).ToNot(HaveOccurred())
						signatures[i] = pack.NewBytes65(signature)
					}
					Expect(utxoTx.Sign(signatures, pack.NewBytes(wif.SerializePubKey()))).To(Succeed())

					// Submit the signed transaction to the UTXO chain's node.
					txHash, err := utxoTx.Hash()
					Expect(err).ToNot(HaveOccurred())
					err = utxoClient.SubmitTx(ctx, utxoTx)
					Expect(err).ToNot(HaveOccurred())

					// Check confirmations after waiting for the transaction to be in the
					// mempool.
					time.Sleep(time.Second)

					for {
						// Loop until the transaction has at least a few
						// confirmations.
						confs, err := confsFn(ctx, txHash)
						Expect(err).ToNot(HaveOccurred())
						logger.Debug(fmt.Sprintf("[%v] confirming", utxoChain.chain), zap.Uint64("current", uint64(confs)))
						if confs >= 1 {
							break
						}
						time.Sleep(10 * time.Second)
					}

					// Load the output and verify that it is equal to the original output.
					output2, _, err := utxoClient.Output(ctx, output.Outpoint)
					Expect(err).ToNot(HaveOccurred())
					Expect(reflect.DeepEqual(output, output2)).To(BeTrue())
				})
			})
		}
	})
})

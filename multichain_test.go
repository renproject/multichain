package multichain_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing/quick"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	cosmossdk "github.com/cosmos/cosmos-sdk/types"
	filaddress "github.com/filecoin-project/go-address"
	filtypes "github.com/filecoin-project/lotus/chain/types"
	"github.com/renproject/id"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/bitcoin"
	"github.com/renproject/multichain/chain/bitcoincash"

	// "github.com/renproject/multichain/chain/digibyte"
	"github.com/renproject/multichain/chain/dogecoin"
	"github.com/renproject/multichain/chain/filecoin"
	"github.com/renproject/multichain/chain/terra"
	"github.com/renproject/multichain/chain/zcash"
	"github.com/renproject/pack"
	"github.com/renproject/surge"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testBTC  = flag.Bool("btc", false, "Pass this flag to test Bitcoin")
	testBCH  = flag.Bool("bch", false, "Pass this flag to test Bitcoincash")
	testDOGE = flag.Bool("doge", false, "Pass this flag to test Dogecoin")
	testFIL  = flag.Bool("fil", false, "Pass this flag to test Filecoin")
	testLUNA = flag.Bool("luna", false, "Pass this flag to test Terra")
	testZEC  = flag.Bool("zec", false, "Pass this flag to test Zcash")
)

var _ = Describe("Multichain", func() {
	// new randomness
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create context to work within.
	ctx := context.Background()

	// Initialise the logger.
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := loggerConfig.Build()
	Expect(err).ToNot(HaveOccurred())

	// Populate the test flags by underlying asset chain.
	testFlags := map[multichain.Chain]bool{}
	testFlags[multichain.Bitcoin] = *testBTC
	testFlags[multichain.BitcoinCash] = *testBCH
	testFlags[multichain.Dogecoin] = *testDOGE
	testFlags[multichain.Filecoin] = *testFIL
	testFlags[multichain.Terra] = *testLUNA
	testFlags[multichain.Zcash] = *testZEC

	//
	// Multichain Configs
	//
	Context("Multichain Declarations", func() {
		Context("All supporting chains/assets are declared", func() {
			accountChains := []struct {
				chain multichain.Chain
				asset multichain.Asset
			}{
				{
					multichain.Fantom,
					multichain.FTM,
				},
				{
					multichain.Filecoin,
					multichain.FIL,
				},
				{
					multichain.Moonbeam,
					multichain.GLMR,
				},
				{
					multichain.Polygon,
					multichain.MATIC,
				},
				{
					multichain.Solana,
					multichain.SOL,
				},
				{
					multichain.Terra,
					multichain.LUNA,
				},
			}
			utxoChains := []struct {
				chain multichain.Chain
				asset multichain.Asset
			}{
				{
					multichain.Bitcoin,
					multichain.BTC,
				},
				{
					multichain.BitcoinCash,
					multichain.BCH,
				},
				{
					multichain.DigiByte,
					multichain.DGB,
				},
				{
					multichain.Dogecoin,
					multichain.DOGE,
				},
				{
					multichain.Zcash,
					multichain.ZEC,
				},
			}

			for _, accountChain := range accountChains {
				accountChain := accountChain
				Specify(fmt.Sprintf("Chain=%v, Asset=%v should be supported", accountChain.chain, accountChain.asset), func() {
					Expect(accountChain.chain.IsAccountBased()).To(BeTrue())
					Expect(accountChain.chain.NativeAsset()).To(Equal(accountChain.asset))
					Expect(accountChain.asset.OriginChain()).To(Equal(accountChain.chain))
				})
			}
			for _, utxoChain := range utxoChains {
				utxoChain := utxoChain
				Specify(fmt.Sprintf("Chain=%v, Asset=%v should be supported", utxoChain.chain, utxoChain.asset), func() {
					Expect(utxoChain.chain.IsUTXOBased()).To(BeTrue())
					Expect(utxoChain.chain.NativeAsset()).To(Equal(utxoChain.asset))
					Expect(utxoChain.asset.OriginChain()).To(Equal(utxoChain.chain))
				})
			}
		})
	})

	//
	// ADDRESS API
	//
	Context("Address API", func() {
		chainTable := []struct {
			chain            multichain.Chain
			newEncodeDecoder func() multichain.AddressEncodeDecoder
			newAddress       func() multichain.Address
			newRawAddress    func() multichain.RawAddress
			newSHAddress     func() multichain.Address
			newSHRawAddress  func() multichain.RawAddress
		}{
			{
				multichain.Bitcoin,
				func() multichain.AddressEncodeDecoder {
					addrEncodeDecoder := bitcoin.NewAddressEncodeDecoder(&chaincfg.RegressionNetParams)
					return addrEncodeDecoder
				},
				func() multichain.Address {
					// Generate a random SECP256K1 private key.
					pk := id.NewPrivKey()
					// Get bitcoin WIF private key with the pub key configured to be in
					// the compressed form.
					wif, err := btcutil.NewWIF((*btcec.PrivateKey)(pk), &chaincfg.RegressionNetParams, true)
					Expect(err).NotTo(HaveOccurred())
					addrPubKeyHash, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.RegressionNetParams)
					Expect(err).NotTo(HaveOccurred())
					// Return the human-readable encoded bitcoin address in base58 format.
					return multichain.Address(addrPubKeyHash.EncodeAddress())
				},
				func() multichain.RawAddress {
					// Generate a random SECP256K1 private key.
					pk := id.NewPrivKey()
					// Get bitcoin WIF private key with the pub key configured to be in
					// the compressed form.
					wif, err := btcutil.NewWIF((*btcec.PrivateKey)(pk), &chaincfg.RegressionNetParams, true)
					Expect(err).NotTo(HaveOccurred())
					// Get the address pubKey hash. This is the most commonly used format
					// for a bitcoin address.
					addrPubKeyHash, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.RegressionNetParams)
					Expect(err).NotTo(HaveOccurred())
					// Encode into the checksummed base58 format.
					encoded := addrPubKeyHash.EncodeAddress()
					return multichain.RawAddress(pack.Bytes(base58.Decode(encoded)))
				},
				func() multichain.Address {
					// Random bytes of script.
					script := make([]byte, r.Intn(100))
					r.Read(script)
					// Create address script hash from the random script bytes.
					addrScriptHash, err := btcutil.NewAddressScriptHash(script, &chaincfg.RegressionNetParams)
					Expect(err).NotTo(HaveOccurred())
					// Return in human-readable encoded form.
					return multichain.Address(addrScriptHash.EncodeAddress())
				},
				func() multichain.RawAddress {
					// Random bytes of script.
					script := make([]byte, r.Intn(100))
					r.Read(script)
					// Create address script hash from the random script bytes.
					addrScriptHash, err := btcutil.NewAddressScriptHash(script, &chaincfg.RegressionNetParams)
					Expect(err).NotTo(HaveOccurred())
					// Encode to the checksummed base58 format.
					encoded := addrScriptHash.EncodeAddress()
					return multichain.RawAddress(pack.Bytes(base58.Decode(encoded)))
				},
			},
			{
				multichain.Filecoin,
				func() multichain.AddressEncodeDecoder {
					return filecoin.NewAddressEncodeDecoder()
				},
				func() multichain.Address {
					pubKey := make([]byte, 64)
					r.Read(pubKey)
					addr, err := filaddress.NewSecp256k1Address(pubKey)
					Expect(err).NotTo(HaveOccurred())
					return multichain.Address(addr.String())
				},
				func() multichain.RawAddress {
					rawAddr := make([]byte, 20)
					r.Read(rawAddr)
					formattedRawAddr := append([]byte{byte(filaddress.SECP256K1)}, rawAddr[:]...)
					return multichain.RawAddress(pack.NewBytes(formattedRawAddr[:]))
				},
				func() multichain.Address {
					return multichain.Address("")
				},
				func() multichain.RawAddress {
					return multichain.RawAddress([]byte{})
				},
			},
			{
				multichain.Terra,
				func() multichain.AddressEncodeDecoder {
					return terra.NewAddressEncodeDecoder()
				},
				func() multichain.Address {
					pk := secp256k1.GenPrivKey()
					addr := cosmossdk.AccAddress(pk.PubKey().Address())
					return multichain.Address(addr.String())
				},
				func() multichain.RawAddress {
					pk := secp256k1.GenPrivKey()
					rawAddr := pk.PubKey().Address()
					return multichain.RawAddress(pack.Bytes(rawAddr))
				},
				func() multichain.Address {
					return multichain.Address("")
				},
				func() multichain.RawAddress {
					return multichain.RawAddress([]byte{})
				},
			},
			{
				multichain.BitcoinCash,
				func() multichain.AddressEncodeDecoder {
					addrEncodeDecoder := bitcoincash.NewAddressEncodeDecoder(&chaincfg.RegressionNetParams)
					return addrEncodeDecoder
				},
				func() multichain.Address {
					pk := id.NewPrivKey()
					wif, err := btcutil.NewWIF((*btcec.PrivateKey)(pk), &chaincfg.RegressionNetParams, true)
					Expect(err).NotTo(HaveOccurred())
					addrPubKeyHash, err := bitcoincash.NewAddressPubKeyHash(btcutil.Hash160(wif.PrivKey.PubKey().SerializeUncompressed()), &chaincfg.RegressionNetParams)
					Expect(err).NotTo(HaveOccurred())
					return multichain.Address(addrPubKeyHash.EncodeAddress())
				},
				func() multichain.RawAddress {
					pk := id.NewPrivKey()
					wif, err := btcutil.NewWIF((*btcec.PrivateKey)(pk), &chaincfg.RegressionNetParams, true)
					Expect(err).NotTo(HaveOccurred())
					addrPubKeyHash, err := bitcoincash.NewAddressPubKeyHash(btcutil.Hash160(wif.PrivKey.PubKey().SerializeUncompressed()), &chaincfg.RegressionNetParams)
					Expect(err).NotTo(HaveOccurred())

					addrBytes := addrPubKeyHash.ScriptAddress()
					addrBytes = append([]byte{0x00}, addrBytes...)
					return multichain.RawAddress(pack.Bytes(addrBytes))
				},
				func() multichain.Address {
					script := make([]byte, r.Intn(100))
					r.Read(script)
					addrScriptHash, err := bitcoincash.NewAddressScriptHash(script, &chaincfg.RegressionNetParams)
					Expect(err).NotTo(HaveOccurred())
					return multichain.Address(addrScriptHash.EncodeAddress())
				},
				func() multichain.RawAddress {
					script := make([]byte, r.Intn(100))
					r.Read(script)
					addrScriptHash, err := bitcoincash.NewAddressScriptHash(script, &chaincfg.RegressionNetParams)
					Expect(err).NotTo(HaveOccurred())

					addrBytes := addrScriptHash.ScriptAddress()
					addrBytes = append([]byte{8}, addrBytes...)
					return multichain.RawAddress(pack.Bytes(addrBytes))
				},
			},
			{
				multichain.Zcash,
				func() multichain.AddressEncodeDecoder {
					addrEncodeDecoder := zcash.NewAddressEncodeDecoder(&zcash.RegressionNetParams)
					return addrEncodeDecoder
				},
				func() multichain.Address {
					pk := id.NewPrivKey()
					wif, err := btcutil.NewWIF((*btcec.PrivateKey)(pk), zcash.RegressionNetParams.Params, true)
					Expect(err).NotTo(HaveOccurred())
					addrPubKeyHash, err := zcash.NewAddressPubKeyHash(btcutil.Hash160(wif.PrivKey.PubKey().SerializeUncompressed()), &zcash.RegressionNetParams)
					Expect(err).NotTo(HaveOccurred())
					return multichain.Address(addrPubKeyHash.EncodeAddress())
				},
				func() multichain.RawAddress {
					pk := id.NewPrivKey()
					wif, err := btcutil.NewWIF((*btcec.PrivateKey)(pk), zcash.RegressionNetParams.Params, true)
					Expect(err).NotTo(HaveOccurred())
					addrPubKeyHash, err := zcash.NewAddressPubKeyHash(btcutil.Hash160(wif.PrivKey.PubKey().SerializeUncompressed()), &zcash.RegressionNetParams)
					Expect(err).NotTo(HaveOccurred())
					return multichain.RawAddress(pack.Bytes(base58.Decode(addrPubKeyHash.EncodeAddress())))
				},
				func() multichain.Address {
					script := make([]byte, r.Intn(100))
					r.Read(script)
					addrScriptHash, err := zcash.NewAddressScriptHash(script, &zcash.RegressionNetParams)
					Expect(err).NotTo(HaveOccurred())
					return multichain.Address(addrScriptHash.EncodeAddress())
				},
				func() multichain.RawAddress {
					script := make([]byte, r.Intn(100))
					r.Read(script)
					addrScriptHash, err := zcash.NewAddressScriptHash(script, &zcash.RegressionNetParams)
					Expect(err).NotTo(HaveOccurred())
					return multichain.RawAddress(pack.Bytes(base58.Decode(addrScriptHash.EncodeAddress())))
				},
			},
		}

		for _, chain := range chainTable {
			chain := chain
			Context(fmt.Sprintf("%v", chain.chain), func() {
				encodeDecoder := chain.newEncodeDecoder()

				It("should encode a raw address correctly", func() {
					rawAddr := chain.newRawAddress()
					encodedAddr, err := encodeDecoder.EncodeAddress(rawAddr)
					Expect(err).NotTo(HaveOccurred())
					decodedRawAddr, err := encodeDecoder.DecodeAddress(encodedAddr)
					Expect(err).NotTo(HaveOccurred())
					Expect(decodedRawAddr).To(Equal(rawAddr))
				})

				It("should decode an address correctly", func() {
					addr := chain.newAddress()
					decodedRawAddr, err := encodeDecoder.DecodeAddress(addr)
					Expect(err).NotTo(HaveOccurred())
					encodedAddr, err := encodeDecoder.EncodeAddress(decodedRawAddr)
					Expect(err).NotTo(HaveOccurred())
					Expect(encodedAddr).To(Equal(addr))
				})

				if chain.chain.IsUTXOBased() {
					It("should encoded a raw script address correctly", func() {
						rawScriptAddr := chain.newSHRawAddress()
						encodedAddr, err := encodeDecoder.EncodeAddress(rawScriptAddr)
						Expect(err).NotTo(HaveOccurred())
						decodedRawAddr, err := encodeDecoder.DecodeAddress(encodedAddr)
						Expect(err).NotTo(HaveOccurred())
						Expect(decodedRawAddr).To(Equal(rawScriptAddr))
					})

					It("should decode a script address correctly", func() {
						scriptAddr := chain.newSHAddress()
						decodedRawAddr, err := encodeDecoder.DecodeAddress(scriptAddr)
						Expect(err).NotTo(HaveOccurred())
						encodedAddr, err := encodeDecoder.EncodeAddress(decodedRawAddr)
						Expect(err).NotTo(HaveOccurred())
						Expect(encodedAddr).To(Equal(scriptAddr))
					})
				}

				if chain.chain == multichain.Bitcoin {
					mainnetEncodeDecoder := bitcoin.NewAddressEncodeDecoder(&chaincfg.MainNetParams)

					It("should decode a Bech32 address correctly", func() {
						segwitAddrs := []string{
							"bc1qp3gcp95e85rupv9zgj57j0lvsqnzcehawzaax3",
							"bc1qh6fjfx39ae4ahvusc4eggyrwjm65zyu83mzwlx",
							"bc1q3zqxadsagdwjp2fpddn8dk5ge8lf0nn0p750ar",
							"bc1q2lthuszmh0mynte4nzsfqtjjseu6fdrmeffr62",
							"bc1qdqkfrt2hpgncqwut88809he6wxysfw8w3cgsh4",
							"bc1qna5zwwuqcst3dqqx8rmwa66jpa45w28tlypg54",
							"bc1qjk2ytl6uctuxfsyf8dn6ptwfsthfat4hd78l0m",
							"bc1qyg6zhg9dhmkj0wz4svsdz6g0ujll225v0wc5hx",
							"bc1quvtmmjccre6plqslujw7qcy820fycg2q2a73an",
							"bc1qztxl2qc3k90uud846qfeawqzz3aedhq48vv3lu",
							"bc1qvkknfkfhfr0axql478klvjs6sanwj6njym5wf2",
							"bc1qya5t2pj7hqpezcnwh72k69h4cgg3srqwtd0e6w",
						}
						for _, segwitAddr := range segwitAddrs {
							decodedRawAddr, err := mainnetEncodeDecoder.DecodeAddress(multichain.Address(segwitAddr))
							Expect(err).NotTo(HaveOccurred())
							encodedAddr, err := mainnetEncodeDecoder.EncodeAddress(decodedRawAddr)
							Expect(err).NotTo(HaveOccurred())
							Expect(string(encodedAddr)).To(Equal(segwitAddr))
						}
					})

					It("should encode a Bech32 address correctly", func() {
						loop := func() bool {
							l := 21
							if r.Intn(2) == 1 {
								l = 33
							}
							randBytes := make([]byte, l)
							r.Read(randBytes)
							randBytes[0] = byte(0)
							rawAddr := multichain.RawAddress(randBytes)
							encodedAddr, err := mainnetEncodeDecoder.EncodeAddress(rawAddr)
							Expect(err).NotTo(HaveOccurred())
							decodedRawAddr, err := mainnetEncodeDecoder.DecodeAddress(encodedAddr)
							Expect(err).NotTo(HaveOccurred())
							Expect(decodedRawAddr).To(Equal(rawAddr))
							return true
						}
						Expect(quick.Check(loop, nil)).To(Succeed())
					})
				}
			})
		}
	})

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
					pkEnv := os.Getenv("TERRA_PK")
					if pkEnv == "" {
						panic("TERRA_PK is undefined")
					}
					pkBytes, err := hex.DecodeString(pkEnv)
					Expect(err).NotTo(HaveOccurred())
					var pk secp256k1.PrivKeySecp256k1
					copy(pk[:], pkBytes)
					addrEncoder := terra.NewAddressEncoder()
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
					addrEncoder := terra.NewAddressEncoder()
					addr, err := addrEncoder.EncodeAddress(multichain.RawAddress(pack.Bytes(pk.PubKey().Address())))
					Expect(err).NotTo(HaveOccurred())
					return addr
				},
				"http://127.0.0.1:26657",
				func() multichain.Address {
					recipientKey := secp256k1.GenPrivKey()
					addrEncoder := terra.NewAddressEncoder()
					recipient, err := addrEncoder.EncodeAddress(multichain.RawAddress(pack.Bytes(recipientKey.PubKey().Address())))
					Expect(err).NotTo(HaveOccurred())
					return recipient
				},
				func(rpcURL pack.String) (multichain.AccountClient, multichain.AccountTxBuilder) {
					client := terra.NewClient(
						terra.DefaultClientOptions().
							WithHost(rpcURL).
							WithCoinDenom("uluna"),
					)
					txBuilder := terra.NewTxBuilder(
						terra.DefaultTxBuilderOptions().
							WithChainID("testnet"),
						client,
					)

					return client, txBuilder
				},
				func(_ multichain.AccountClient) (pack.U256, pack.U256, pack.U256, pack.U256, pack.Bytes) {
					amount := pack.NewU256FromU64(pack.U64(2000000))
					gasLimit := pack.NewU256FromU64(pack.U64(100000))
					gasPrice := pack.NewU256FromU64(pack.U64(1))
					gasCap := pack.NewU256FromInt(gasPrice.Int())
					payload := pack.NewBytes([]byte("multichain"))
					return amount, gasLimit, gasPrice, gasCap, payload
				},
				multichain.Terra,
			},
			{
				func() (id.PrivKey, *id.PubKey, multichain.Address) {
					pkEnv := os.Getenv("FILECOIN_PK")
					if pkEnv == "" {
						panic("FILECOIN_PK is undefined")
					}
					var ki filtypes.KeyInfo
					data, err := hex.DecodeString(pkEnv)
					Expect(err).NotTo(HaveOccurred())
					err = json.Unmarshal(data, &ki)
					Expect(err).NotTo(HaveOccurred())
					privKey := id.PrivKey{}
					err = surge.FromBinary(&privKey, ki.PrivateKey)
					Expect(err).NotTo(HaveOccurred())
					pubKey := privKey.PubKey()

					// FIXME: add method in renproject/id to get uncompressed pubkey bytes
					pubKeyCompressed, err := surge.ToBinary(pubKey)
					Expect(err).NotTo(HaveOccurred())
					/*addr*/ _, err = filaddress.NewSecp256k1Address(pubKeyCompressed)
					Expect(err).NotTo(HaveOccurred())
					addrStr := os.Getenv("FILECOIN_ADDRESS")
					if addrStr == "" {
						panic("FILECOIN_ADDRESS is undefined")
					}

					return privKey, pubKey, multichain.Address(pack.String(addrStr))
				},
				func(privKey id.PrivKey) multichain.Address {
					pubKey := privKey.PubKey()
					pubKeyCompressed, err := surge.ToBinary(pubKey)
					Expect(err).NotTo(HaveOccurred())
					addr, err := filaddress.NewSecp256k1Address(pubKeyCompressed)
					Expect(err).NotTo(HaveOccurred())
					return multichain.Address(pack.String(addr.String()))
				},
				"http://127.0.0.1:1234/rpc/v0",
				func() multichain.Address {
					pk := id.NewPrivKey()
					pubKey := pk.PubKey()
					pubKeyCompressed, err := surge.ToBinary(pubKey)
					Expect(err).NotTo(HaveOccurred())
					addr, err := filaddress.NewSecp256k1Address(pubKeyCompressed)
					Expect(err).NotTo(HaveOccurred())
					return multichain.Address(pack.String(addr.String()))
				},
				func(rpcURL pack.String) (multichain.AccountClient, multichain.AccountTxBuilder) {
					// dirty hack to fetch auth token
					client, err := filecoin.NewClient(
						filecoin.DefaultClientOptions().
							WithRPCURL(rpcURL).
							WithAuthToken(fetchAuthToken()),
					)
					Expect(err).NotTo(HaveOccurred())

					txBuilder := filecoin.NewTxBuilder()

					return client, txBuilder
				},
				func(client multichain.AccountClient) (pack.U256, pack.U256, pack.U256, pack.U256, pack.Bytes) {
					amount := pack.NewU256FromU64(pack.NewU64(100000000))
					gasLimit := pack.NewU256FromU64(pack.NewU64(2189560))

					// Fetch gas price and gas cap using the gas estimator.
					filecoinClient := client.(*filecoin.Client)
					gasPrice, gasCap, err := filecoin.NewGasEstimator(filecoinClient, gasLimit.Int().Int64()).
						EstimateGas(context.Background())
					Expect(err).NotTo(HaveOccurred())

					payload := pack.Bytes(nil)
					return amount, gasLimit, gasPrice, gasCap, payload
				},
				multichain.Filecoin,
			},
		}

		for _, accountChain := range accountChainTable {
			accountChain := accountChain
			if !testFlags[accountChain.chain] {
				continue
			}

			Context(fmt.Sprintf("%v", accountChain.chain), func() {
				Specify("build, broadcast and fetch tx", func() {
					// Load private key and the associated address.
					senderPrivKey, senderPubKey, senderAddr := accountChain.senderEnv()

					// Get a random recipient address.
					recipientAddr := accountChain.randomRecipientAddr()

					// Initialise the account chain's client, and possibly get a nonce for
					// the sender.
					accountClient, txBuilder := accountChain.initialise(accountChain.rpcURL)

					// Get the appropriate nonce for sender.
					nonce, err := accountClient.AccountNonce(ctx, senderAddr)
					Expect(err).NotTo(HaveOccurred())

					// Build a transaction.
					amount, gasLimit, gasPrice, gasCap, payload := accountChain.txParams(accountClient)

					accountTx, err := txBuilder.BuildTx(
						ctx,
						multichain.Address(senderAddr),
						recipientAddr,
						amount, nonce, gasLimit, gasPrice, gasCap,
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

					// Submit the transaction to the account chain.
					txHash := accountTx.Hash()
					err = accountClient.SubmitTx(ctx, accountTx)
					Expect(err).NotTo(HaveOccurred())
					logger.Debug("submit tx", zap.String("from", string(senderAddr)), zap.String("to", string(recipientAddr)), zap.Any("txHash", txHash))

					// Wait slightly before we query the chain's node.
					time.Sleep(time.Second)

					for {
						// Loop until the transaction has at least a few confirmations.
						tx, confs, err := accountClient.Tx(ctx, txHash)
						if err == nil {
							Expect(confs.Uint64()).To(BeNumerically(">", 0))
							Expect(tx.Value()).To(Equal(amount))
							Expect(tx.From()).To(Equal(senderAddr))
							Expect(tx.To()).To(Equal(recipientAddr))
							Expect(tx.Value()).To(Equal(amount))
							break
						}

						// wait and retry querying for the transaction
						time.Sleep(5 * time.Second)
					}
				})

				It("should be able to fetch the latest block", func() {
					// Initialise client
					accountClient, _ := accountChain.initialise(accountChain.rpcURL)

					latestBlock, err := accountClient.LatestBlock(ctx)
					Expect(err).NotTo(HaveOccurred())
					Expect(uint64(latestBlock)).To(BeNumerically(">", 1))
				})
			})
		}
	})

	//
	// UTXO API
	//
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
			/*
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
			*/
		}

		for _, utxoChain := range utxoChainTable {
			utxoChain := utxoChain
			if !testFlags[utxoChain.chain] {
				continue
			}

			Context(fmt.Sprintf("%v", utxoChain.chain), func() {
				Specify("(P2PKH) build, broadcast and fetch tx", func() {
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

					// Recipient 1
					pkhAddrUncompressed, err := utxoChain.newAddressPKH(btcutil.Hash160(wif.PrivKey.PubKey().SerializeUncompressed()))
					Expect(err).ToNot(HaveOccurred())

					// Recipient 2
					recipientPrivKey := id.NewPrivKey()
					recipientPubKey := recipientPrivKey.PubKey()
					recipientPubKeyCompressed, err := surge.ToBinary(recipientPubKey)
					Expect(err).NotTo(HaveOccurred())
					recipientPkhAddr, err := utxoChain.newAddressPKH(btcutil.Hash160(((*btcec.PublicKey)(recipientPubKey)).SerializeCompressed()))
					Expect(err).NotTo(HaveOccurred())

					// Initialise the UTXO client and fetch the unspent outputs. Also get a
					// function to query the number of block confirmations for a transaction.
					utxoClient, unspentOutputs, confsFn := utxoChain.initialise(utxoChain.rpcURL, pkhAddr)
					Expect(len(unspentOutputs)).To(BeNumerically(">", 0))
					var output multichain.UTXOutput
					thresholdValue := pack.NewU256FromU64(pack.NewU64(2500))
					for _, unspentOutput := range unspentOutputs {
						if unspentOutput.Value.GreaterThan(thresholdValue) {
							output = unspentOutput
							break
						}
					}

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
					utxoValue1 := pack.NewU256FromU64(pack.NewU64((output.Value.Int().Uint64() - 1000) / 4))
					utxoValue2 := pack.NewU256FromU64(pack.NewU64((output.Value.Int().Uint64() - 1000) * 3 / 4))
					recipients := []multichain.UTXORecipient{
						{
							To:    multichain.Address(pkhAddrUncompressed.EncodeAddress()),
							Value: utxoValue1,
						},
						{
							To:    multichain.Address(recipientPkhAddr.EncodeAddress()),
							Value: utxoValue2,
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

					// Load the first output and verify the value.
					output3, _, err := utxoClient.Output(ctx, multichain.UTXOutpoint{
						Hash:  txHash,
						Index: pack.NewU32(0),
					})
					Expect(err).ToNot(HaveOccurred())
					Expect(output3.Value).To(Equal(utxoValue1))

					// Load the second output and verify the value.
					output4, _, err := utxoClient.Output(ctx, multichain.UTXOutpoint{
						Hash:  txHash,
						Index: pack.NewU32(1),
					})
					Expect(err).ToNot(HaveOccurred())
					Expect(output4.Value).To(Equal(utxoValue2))

					// Construct UTXO to be signed by invalid key. This UTXO should fail
					// when submitted to the network, since the signer doesn't have the
					// right to spend it.
					// We submit the invalid signed UTXO (which should fail), and wait
					// for a maximum of 5 seconds.
					inputs2 := []multichain.UTXOInput{{
						Output: output4,
					}}
					recipients2 := []multichain.UTXORecipient{{
						To:    multichain.Address(pkhAddr.EncodeAddress()),
						Value: output4.Value.Sub(pack.NewU256FromU64(pack.U64(500))),
					}}
					utxoTx2, err := utxoChain.txBuilder.BuildTx(inputs2, recipients2)
					Expect(err).NotTo(HaveOccurred())
					sighashes2, err := utxoTx2.Sighashes()
					signatures2 := make([]pack.Bytes65, len(sighashes2))
					for i := range sighashes2 {
						hash := id.Hash(sighashes2[i])
						privKey := (*id.PrivKey)(wif.PrivKey)
						signature, err := privKey.Sign(&hash)
						Expect(err).ToNot(HaveOccurred())
						signatures2[i] = pack.NewBytes65(signature)
					}
					Expect(utxoTx2.Sign(signatures2, pack.NewBytes(wif.SerializePubKey()))).To(Succeed())
					failingCtx, failingCancelFn := context.WithTimeout(ctx, 5*time.Second)
					Expect(utxoClient.SubmitTx(failingCtx, utxoTx2)).To(HaveOccurred())
					failingCancelFn()

					// Try to spend UTXO from valid key. We should be able to successfully
					// submit the signed UTXO to the network.
					utxoTx3, err := utxoChain.txBuilder.BuildTx(inputs2, recipients2)
					Expect(err).NotTo(HaveOccurred())
					sighashes3, err := utxoTx3.Sighashes()
					signatures3 := make([]pack.Bytes65, len(sighashes3))
					for i := range sighashes3 {
						hash := id.Hash(sighashes3[i])
						signature, err := recipientPrivKey.Sign(&hash)
						Expect(err).ToNot(HaveOccurred())
						signatures3[i] = pack.NewBytes65(signature)
					}
					Expect(utxoTx3.Sign(signatures3, pack.NewBytes(recipientPubKeyCompressed))).To(Succeed())
					Expect(utxoClient.SubmitTx(ctx, utxoTx3)).NotTo(HaveOccurred())
				})

				Specify("(P2SH)  build, broadcast and fetch tx", func() {
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
					recipientPrivKey := id.NewPrivKey()
					recipientPubKey := recipientPrivKey.PubKey()
					recipientPubKeyCompressed, err := surge.ToBinary(recipientPubKey)
					Expect(err).NotTo(HaveOccurred())
					pubKey := pack.Bytes(((*btcec.PublicKey)(recipientPubKey)).SerializeCompressed())
					script, err := getScript(pubKey)
					Expect(err).NotTo(HaveOccurred())
					pubKeyScript, err := getPubKeyScript(pubKey)
					Expect(err).NotTo(HaveOccurred())
					recipientP2SH, err := utxoChain.newAddressSH(script)
					Expect(err).NotTo(HaveOccurred())

					// Initialise the UTXO client and fetch the unspent outputs. Also get a
					// function to query the number of block confirmations for a transaction.
					utxoClient, unspentOutputs, confsFn := utxoChain.initialise(utxoChain.rpcURL, pkhAddr)
					Expect(len(unspentOutputs)).To(BeNumerically(">", 0))
					var output multichain.UTXOutput
					thresholdValue := pack.NewU256FromU64(pack.NewU64(2500))
					for _, unspentOutput := range unspentOutputs {
						if unspentOutput.Value.GreaterThan(thresholdValue) {
							output = unspentOutput
							break
						}
					}

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
							To:    multichain.Address(recipientP2SH.EncodeAddress()),
							Value: output.Value.Sub(pack.NewU256FromU64(pack.U64(500))),
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
					logger.Debug("[P2KH -> P2SH] submit tx", zap.String("from", pkhAddr.EncodeAddress()), zap.String("to", recipientP2SH.EncodeAddress()), zap.String("txHash", string(txHashToHex(txHash))))

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

					// Load the output and verify that the pub key script is as calculated
					// initially.
					output2, _, err := utxoClient.Output(ctx, multichain.UTXOutpoint{
						Hash:  txHash,
						Index: pack.NewU32(0),
					})
					Expect(err).ToNot(HaveOccurred())
					Expect(output2.PubKeyScript.Equal(pubKeyScript)).To(BeTrue())

					// Validate that the output2 is spendable
					sigScript, err := getScript(pubKey)
					Expect(err).NotTo(HaveOccurred())
					inputs2 := []multichain.UTXOInput{{
						Output:    output2,
						SigScript: sigScript,
					}}
					recipients2 := []multichain.UTXORecipient{{
						To:    multichain.Address(pkhAddr.EncodeAddress()),
						Value: output2.Value.Sub(pack.NewU256FromU64(pack.U64(500))),
					}}
					utxoTx2, err := utxoChain.txBuilder.BuildTx(inputs2, recipients2)
					Expect(err).NotTo(HaveOccurred())

					// Create another transaction using the same inputs, which we will
					// sign with the original user's address. Validate that none other
					// than the recipient's signature can spend this UTXO.
					utxoTx3, err := utxoChain.txBuilder.BuildTx(inputs2, recipients2)
					Expect(err).NotTo(HaveOccurred())

					// Get the sighashes that need to be signed, and sign them.
					sighashes2, err := utxoTx2.Sighashes()
					signatures2 := make([]pack.Bytes65, len(sighashes2))
					signatures3 := make([]pack.Bytes65, len(sighashes2))
					Expect(err).ToNot(HaveOccurred())
					for i := range sighashes2 {
						hash := id.Hash(sighashes2[i])
						signature, err := recipientPrivKey.Sign(&hash)
						Expect(err).ToNot(HaveOccurred())
						signatures2[i] = pack.NewBytes65(signature)
					}
					for i := range sighashes2 {
						hash := id.Hash(sighashes2[i])
						privKey := (*id.PrivKey)(wif.PrivKey)
						signature, err := privKey.Sign(&hash)
						Expect(err).ToNot(HaveOccurred())
						signatures3[i] = pack.NewBytes65(signature)
					}
					Expect(utxoTx2.Sign(signatures2, pack.NewBytes(recipientPubKeyCompressed))).To(Succeed())
					Expect(utxoTx3.Sign(signatures3, pack.NewBytes(wif.SerializePubKey()))).To(Succeed())

					// Try to submit tx signed by invalid spender. This should fail since
					failingCtx, failingCancelFn := context.WithTimeout(ctx, 5*time.Second)
					Expect(utxoClient.SubmitTx(failingCtx, utxoTx3)).To(HaveOccurred())
					failingCancelFn()

					// Submit the signed transaction to the UTXO chain's node.
					txHash2, err := utxoTx2.Hash()
					Expect(err).ToNot(HaveOccurred())
					err = utxoClient.SubmitTx(ctx, utxoTx2)
					Expect(err).ToNot(HaveOccurred())
					logger.Debug("[P2SH -> P2KH] submit tx", zap.String("from", recipientP2SH.EncodeAddress()), zap.String("to", pkhAddr.EncodeAddress()), zap.String("txHash", string(txHashToHex(txHash2))))

					// Check confirmations after waiting for the transaction to be in the
					// mempool.
					time.Sleep(time.Second)

					for {
						// Loop until the transaction has at least a few
						// confirmations.
						confs, err := confsFn(ctx, txHash2)
						Expect(err).ToNot(HaveOccurred())
						logger.Debug(fmt.Sprintf("[%v] confirming", utxoChain.chain), zap.Uint64("current", uint64(confs)))
						if confs >= 1 {
							break
						}
						time.Sleep(10 * time.Second)
					}
				})

				It("should be able to fetch the latest block", func() {
					// get a random address
					randAddr := make([]byte, 20)
					r.Read(randAddr)
					pkhAddr, err := utxoChain.newAddressPKH(randAddr)
					Expect(err).NotTo(HaveOccurred())

					// initialise client
					utxoClient, _, _ := utxoChain.initialise(utxoChain.rpcURL, pkhAddr)

					latestBlock, err := utxoClient.LatestBlock(ctx)
					Expect(err).NotTo(HaveOccurred())
					Expect(uint64(latestBlock)).To(BeNumerically(">", 1))
				})
			})
		}
	})
})

func txHashToHex(txHash pack.Bytes) pack.String {
	// bitcoin's msgTx is a byte-reversed hash
	// https://github.com/btcsuite/btcd/blob/master/chaincfg/chainhash/hash.go#L27-L28
	txHashCopy := make([]byte, len(txHash))
	copy(txHashCopy[:], txHash)
	hashSize := len(txHashCopy)
	for i := 0; i < hashSize/2; i++ {
		txHashCopy[i], txHashCopy[hashSize-1-i] = txHashCopy[hashSize-1-i], txHashCopy[i]
	}
	return pack.String(hex.EncodeToString(txHashCopy))
}

func fetchAuthToken() pack.String {
	// fetch the auth token from filecoin's running docker container
	cmd := exec.Command("docker", "exec", "infra_filecoin_1", "/bin/bash", "-c", "/app/lotus auth api-info --perm admin")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		panic(fmt.Sprintf("could not run command: %v", err))
	}
	tokenWithSuffix := strings.TrimPrefix(out.String(), "FULLNODE_API_INFO=")
	authToken := strings.Split(tokenWithSuffix, ":/")
	return pack.NewString(fmt.Sprintf("Bearer %s", authToken[0]))
}

func getScript(pubKey pack.Bytes) (pack.Bytes, error) {
	pubKeyHash160 := btcutil.Hash160(pubKey)
	return txscript.NewScriptBuilder().
		AddOp(txscript.OP_DUP).
		AddOp(txscript.OP_HASH160).
		AddData(pubKeyHash160).
		AddOp(txscript.OP_EQUALVERIFY).
		AddOp(txscript.OP_CHECKSIG).
		Script()
}

func getPubKeyScript(pubKey pack.Bytes) (pack.Bytes, error) {
	script, err := getScript(pubKey)
	if err != nil {
		return nil, fmt.Errorf("invalid script: %v", err)
	}
	pubKeyScript, err := txscript.NewScriptBuilder().
		AddOp(txscript.OP_HASH160).
		AddData(btcutil.Hash160(script)).
		AddOp(txscript.OP_EQUAL).
		Script()
	if err != nil {
		return nil, fmt.Errorf("invalid pubkeyscript: %v", err)
	}
	return pubKeyScript, nil
}

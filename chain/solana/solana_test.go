package solana_test

import (
	"context"
	"encoding/hex"
	"os"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/bitcoin"
	"github.com/renproject/multichain/chain/solana"
	"github.com/renproject/multichain/chain/solana/solana-ffi/cgo"
	"github.com/renproject/pack"
	"github.com/renproject/surge"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Solana", func() {
	// Setup logger.
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := loggerConfig.Build()
	Expect(err).ToNot(HaveOccurred())

	Context("mint and burn", func() {
		It("should succeed", func() {
			// Base58 address of the RenBridge program that is deployed to Solana.
			program := multichain.Address("9TaQuUfNMC5rFvdtzhHPk84WaFH3SFnweZn4tw9RriDP")

			// Construct user's keypair path (~/.config/solana/id.json).
			userHomeDir, err := os.UserHomeDir()
			Expect(err).NotTo(HaveOccurred())
			keypairPath := userHomeDir + "/.config/solana/id.json"

			// RenVM secret and corresponding authority (20-byte Ethereum address).
			renVmSecret := "0000000000000000000000000000000000000000000000000000000000000001"
			renVmAuthority := "7E5F4552091A69125d5DfCb7b8C2659029395Bdf"
			renVmAuthorityBytes, err := hex.DecodeString(renVmAuthority)
			Expect(err).NotTo(HaveOccurred())

			// Initialize RenBridge program.
			initializeSig := cgo.RenBridgeInitialize(keypairPath, solana.DefaultClientRPCURL, renVmAuthorityBytes)
			logger.Debug("Initialize", zap.String("tx signature", string(initializeSig)))

			// Initialize RenBTC token.
			time.Sleep(10 * time.Second)
			selector := "BTC/toSolana"
			initializeTokenSig := cgo.RenBridgeInitializeToken(keypairPath, solana.DefaultClientRPCURL, selector)
			logger.Debug("InitializeToken", zap.String("tx signature", string(initializeTokenSig)))

			// Initialize a new token account.
			time.Sleep(10 * time.Second)
			initializeAccountSig := cgo.RenBridgeInitializeAccount(keypairPath, solana.DefaultClientRPCURL, selector)
			logger.Debug("InitializeAccount", zap.String("tx signature", string(initializeAccountSig)))

			// Mint some tokens.
			time.Sleep(10 * time.Second)
			mintAmount := uint64(10000000000) // 10 tokens.
			mintSig := cgo.RenBridgeMint(keypairPath, solana.DefaultClientRPCURL, renVmSecret, selector, mintAmount)
			logger.Debug("Mint", zap.String("tx signature", string(mintSig)))

			// Burn some tokens.
			time.Sleep(10 * time.Second)
			recipient := multichain.Address("mwjUmhAW68zCtgZpW5b1xD5g7MZew6xPV4")
			bitcoinAddrEncodeDecoder := bitcoin.NewAddressEncodeDecoder(&chaincfg.RegressionNetParams)
			recipientRawAddr, err := bitcoinAddrEncodeDecoder.DecodeAddress(recipient)
			Expect(err).NotTo(HaveOccurred())
			burnCount := cgo.RenBridgeGetBurnCount(solana.DefaultClientRPCURL)
			burnAmount := uint64(5000000000) // 5 tokens.
			burnSig := cgo.RenBridgeBurn(keypairPath, solana.DefaultClientRPCURL, selector, burnCount, burnAmount, []byte(recipientRawAddr))
			logger.Debug("Burn", zap.String("tx signature", string(burnSig)))

			// Fetch burn log.
			time.Sleep(20 * time.Second)
			client := solana.NewClient(solana.DefaultClientOptions())
			contractCallInput := solana.BurnCallContractInput{Nonce: pack.NewU64(burnCount)}
			calldata, err := surge.ToBinary(contractCallInput)
			Expect(err).ToNot(HaveOccurred())
			burnLogBytes, err := client.CallContract(context.Background(), program, calldata)
			Expect(err).NotTo(HaveOccurred())
			burnLog := solana.BurnCallContractOutput{}
			err = surge.FromBinary(&burnLog, burnLogBytes)
			Expect(err).NotTo(HaveOccurred())
			Expect(burnLog.Amount).To(Equal(pack.U64(burnAmount)))
			Expect(burnLog.Recipient).To(Equal(recipientRawAddr))
		})
	})
})

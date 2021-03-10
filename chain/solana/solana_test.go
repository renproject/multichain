package solana_test

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"os"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/bitcoin"
	"github.com/renproject/multichain/chain/solana"
	"github.com/renproject/solana-ffi/cgo"
	"github.com/renproject/pack"
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
			// Base58 address of the Gateway program that is deployed to Solana.
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

			// Initialize Gateway for the RenBTC token.
			selector := "BTC/toSolana"
			initializeSig := cgo.GatewayInitialize(keypairPath, solana.DefaultClientRPCURL, renVmAuthorityBytes, selector)
			logger.Debug("Initialize", zap.String("tx signature", string(initializeSig)))

			// Initialize a new token account.
			time.Sleep(10 * time.Second)
			initializeAccountSig := cgo.GatewayInitializeAccount(keypairPath, solana.DefaultClientRPCURL, selector)
			logger.Debug("InitializeAccount", zap.String("tx signature", string(initializeAccountSig)))

			// Mint some tokens.
			time.Sleep(10 * time.Second)
			mintAmount := uint64(10000000000) // 10 tokens.
			mintSig := cgo.GatewayMint(keypairPath, solana.DefaultClientRPCURL, renVmSecret, selector, mintAmount)
			logger.Debug("Mint", zap.String("tx signature", string(mintSig)))

			// Burn some tokens.
			time.Sleep(10 * time.Second)
			recipient := multichain.Address("mwjUmhAW68zCtgZpW5b1xD5g7MZew6xPV4")
			bitcoinAddrEncodeDecoder := bitcoin.NewAddressEncodeDecoder(&chaincfg.RegressionNetParams)
			recipientRawAddr, err := bitcoinAddrEncodeDecoder.DecodeAddress(recipient)
			Expect(err).NotTo(HaveOccurred())
			burnCount := cgo.GatewayGetBurnCount(solana.DefaultClientRPCURL)
			burnAmount := uint64(5000000000) // 5 tokens.
			burnSig := cgo.GatewayBurn(keypairPath, solana.DefaultClientRPCURL, selector, burnCount, burnAmount, uint32(len(recipientRawAddr)), []byte(recipientRawAddr))
			logger.Debug("Burn", zap.String("tx signature", string(burnSig)))

			// Fetch burn log.
			time.Sleep(20 * time.Second)
			client := solana.NewClient(solana.DefaultClientOptions())
			calldata := make([]byte, 8)
			binary.LittleEndian.PutUint64(calldata, burnCount)
			data, err := client.CallContract(context.Background(), program, multichain.ContractCallData(calldata))
			Expect(err).NotTo(HaveOccurred())
			Expect(len(data)).To(Equal(41))

			fetchedAmount := binary.LittleEndian.Uint64(data[:8])
			recipientLen := uint8(data[8:9][0])
			fetchedRecipient := pack.Bytes(data[9 : 9+int(recipientLen)])
			Expect(fetchedAmount).To(Equal(burnAmount))
			Expect([]byte(fetchedRecipient)).To(Equal([]byte(recipientRawAddr)))
		})
	})
})

package solana_test

import (
	"context"
	"encoding/binary"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/near/borsh-go"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/solana"
	"github.com/renproject/pack"
	"github.com/renproject/solana-ffi/cgo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Bytes32 is an alias for [32]byte
type Bytes32 = [32]byte

// GatewayRegistry defines the state of gateway registry, serialized and
// deserialized by the borsh schema.
type GatewayRegistry struct {
	IsInitialised uint8
	Owner         Bytes32
	Count         uint64
	Selectors     []Bytes32
	Gateways      []Bytes32
}

var _ = Describe("Solana", func() {
	// Setup logger.
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := loggerConfig.Build()
	Expect(err).ToNot(HaveOccurred())

	Context("When minting and burning", func() {
		It("should succeed", func() {
			// Base58 address of the Gateway program that is deployed to Solana.
			program := multichain.Address("FDdKRjbBeFtyu5c66cZghJsTTjDTT1aD3zsgTWMTpaif")

			// Construct user's keypair path (~/.config/solana/id.json).
			userHomeDir, err := os.UserHomeDir()
			Expect(err).NotTo(HaveOccurred())
			keypairPath := userHomeDir + "/.config/solana/id.json"

			// RenVM secret and the selector for this gateway.
			renVmSecret := "0000000000000000000000000000000000000000000000000000000000000001"
			selector := "BTC/toSolana"

			// Mint some tokens.
			time.Sleep(10 * time.Second)
			mintAmount := uint64(1000000000) // 10 tokens.
			nilSlice := make([]byte, 32)
			mintSig := cgo.GatewayMint(keypairPath, solana.DefaultClientRPCURL, renVmSecret, selector, mintAmount, nilSlice, nilSlice)
			logger.Debug("Mint", zap.String("tx signature", string(mintSig)))

			// Burn some tokens.
			time.Sleep(10 * time.Second)
			recipient := []byte("mwjUmhAW68zCtgZpW5b1xD5g7MZew6xPV4")
			Expect(err).NotTo(HaveOccurred())
			burnCount := cgo.GatewayGetBurnCount(solana.DefaultClientRPCURL, solana.DefaultClientRPCURL)
			burnAmount := uint64(500000000) // 5 tokens.
			burnSig := cgo.GatewayBurn(keypairPath, solana.DefaultClientRPCURL, selector, burnCount, burnAmount, uint32(len(recipient)), recipient)
			logger.Debug("Burn", zap.String("tx signature", string(burnSig)))

			// Fetch burn log.
			time.Sleep(20 * time.Second)
			client := solana.NewClient(solana.DefaultClientOptions())
			calldata := make([]byte, 8)
			binary.LittleEndian.PutUint64(calldata, burnCount)
			data, err := client.CallContract(context.Background(), program, multichain.ContractCallData(calldata))
			Expect(err).NotTo(HaveOccurred())

			Expect(len(data)).To(Equal(97))
			fetchedAmount := [32]byte{}
			copy(fetchedAmount[:], data[0:32])
			recipientLen := uint8(data[32:33][0])
			fetchedRecipient := pack.Bytes(data[33 : 33+int(recipientLen)])
			Expect(pack.NewU256(fetchedAmount)).To(Equal(pack.NewU256FromUint64(burnAmount)))
			Expect([]byte(fetchedRecipient)).To(Equal(recipient))
		})
	})

	Context("When getting Gateways from Registry", func() {
		It("should deserialize successfully", func() {
			// Solana client using default client options.
			client := solana.NewClient(solana.DefaultClientOptions())

			// Base58 address of the Gateway registry program deployed to Solana.
			registryProgram := multichain.Address("DHpzwsdvAzq61PN9ZwQWg2hzwX8gYNfKAdsNKKtdKDux")
			seeds := []byte("GatewayRegistryState")
			registryState := solana.ProgramDerivedAddress(pack.Bytes(seeds), registryProgram)

			// Fetch account data at gateway registry's state
			accountData, err := client.GetAccountData(registryState)
			Expect(err).NotTo(HaveOccurred())

			// Deserialize the account data into registry state's structure.
			registry := GatewayRegistry{}
			err = borsh.Deserialize(&registry, []byte(accountData))
			Expect(err).NotTo(HaveOccurred())

			// The registry (in the CI test environment) is pre-populated with gateway
			// addresses for BTC/toSolana selector.
			btcSelectorHash := [32]byte{}
			copy(btcSelectorHash[:], crypto.Keccak256([]byte("BTC/toSolana")))
			zero := pack.NewU256FromU8(pack.U8(0)).Bytes32()

			addrEncodeDecoder := solana.NewAddressEncodeDecoder()
			expectedBtcGateway, _ := addrEncodeDecoder.DecodeAddress(multichain.Address("FDdKRjbBeFtyu5c66cZghJsTTjDTT1aD3zsgTWMTpaif"))

			Expect(registry.Count).To(Equal(uint64(1)))
			Expect(registry.Selectors[0]).To(Equal(btcSelectorHash))
			Expect(registry.Selectors[1]).To(Equal(zero))
			Expect(len(registry.Selectors)).To(Equal(32))
			Expect(registry.Gateways[0][:]).To(Equal([]byte(expectedBtcGateway)))
			Expect(registry.Gateways[1]).To(Equal(zero))
			Expect(len(registry.Gateways)).To(Equal(32))
		})
	})
})

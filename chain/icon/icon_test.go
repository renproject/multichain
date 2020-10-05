package icon_test

import (
	"fmt"

	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/icon"
	"github.com/renproject/multichain/chain/icon/crypto"
	"github.com/renproject/pack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Multichain", func() {
	// Create context to work within.
	//ctx := context.Background()

	// Initialise the logger.
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := loggerConfig.Build()
	Expect(err).ToNot(HaveOccurred())

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
				multichain.Icon,
				func() multichain.AddressEncodeDecoder {
					return icon.NewAddressEncodeDecoder()
				},
				func() multichain.Address {
					priv_key, pub_key := crypto.GenerateKeyPair()
					addr := icon.NewAccountAddressFromPublicKey(pub_key)
					logger.Debug("Private Key: " + priv_key.String())
					//logger.Debug("Public Key: " + pub_key.String())
					logger.Debug("Public Address: " + addr.String())
					return multichain.Address(addr.String())
				},
				func() multichain.RawAddress {
					priv_key, pub_key := crypto.GenerateKeyPair()
					addr := icon.NewAccountAddressFromPublicKey(pub_key)
					rawAddr := addr.Bytes()
					logger.Debug("Private Key: " + priv_key.String())
					//logger.Debug("Public Key: " + pub_key.String())
					logger.Debug("Public Address: " + addr.String())
					return multichain.RawAddress(pack.Bytes(rawAddr))
				},
				func() multichain.Address {
					return multichain.Address("")
				},
				func() multichain.RawAddress {
					return multichain.RawAddress([]byte{})
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
					//logger.Debug(hex.EncodeToString(decodedRawAddr))
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
			})
		}
	})
})

package harmony_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/chain/harmony"
)

var _ = Describe("Address", func() {
	Context("when decoding a valid address", func() {
		It("should work without errors", func() {
			addrs := []string{
				"an83characterlonghumanreadablepartthatcontainsthenumber1andtheexcludedcharactersbio1tt5tgs",
				"A12UEL5L",
				"abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw",
				"11qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqc8247j",
				"split1checkupstagehandshakeupstreamerranterredcaperred2y9e3w",
			}
			encoderDecoder := harmony.NewEncoderDecoder()
			for _, addr := range addrs {
				_, err := encoderDecoder.DecodeAddress(address.Address(addr))
				Expect(err).ToNot(HaveOccurred())
			}
		})
	})

	Context("when decoding an invalid address", func() {
		It("should work without errors", func() {
			addrs := []string{
				"split1checkupstagehandshakeupstreamerranterredcaperred2y9e2w",
				"s lit1checkupstagehandshakeupstreamerranterredcaperredp8hs2p",
				"spl" + string(rune(127)) + "t1checkupstagehandshakeupstreamerranterredcaperred2y9e3w",
				"split1cheo2y9e2w",
				"split1a2y9w",
				"1checkupstagehandshakeupstreamerranterredcaperred2y9e3w",
				"11qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqsqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqc8247j",
			}
			for _, addr := range addrs {
				_, err := harmony.NewEncoderDecoder().DecodeAddress(address.Address(addr))
				Expect(err).To(HaveOccurred())
			}
		})
	})

	Context("when encoding a valid address", func() {
		It("should work without errors", func() {
			key, _ := crypto.GenerateKey()
			ethAddr := crypto.PubkeyToAddress(key.PublicKey).String()
			encoderDecoder := harmony.NewEncoderDecoder()
			addr, err := encoderDecoder.EncodeAddress(common.HexToAddress(ethAddr).Bytes())
			Expect(err).ToNot(HaveOccurred())

			rawAddr, err := encoderDecoder.DecodeAddress(addr)
			Expect(err).ToNot(HaveOccurred())
			for i, b := range rawAddr {
				Expect(b).To(Equal(common.HexToAddress(ethAddr).Bytes()[i]))
			}
		})
	})

	Context("when encoding/decoding a valid address", func() {
		It("should work without errors", func() {
			addr := address.Address("one1zksj3evekayy90xt4psrz8h6j2v3hla4qwz4ur")
			encoderDecoder := harmony.NewEncoderDecoder()
			rawAddr, err := encoderDecoder.DecodeAddress(addr)
			Expect(err).ToNot(HaveOccurred())

			conv, err := encoderDecoder.EncodeAddress(rawAddr)
			Expect(err).ToNot(HaveOccurred())
			Expect(conv).To(Equal(addr))
		})
	})
})
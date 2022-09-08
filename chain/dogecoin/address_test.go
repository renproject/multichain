package dogecoin_test

import (
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/dogecoin"
	"github.com/renproject/multichain/chain/utxochain"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dogecoin", func() {
	Context("when decoding segwit address", func() {
		Context("when decoding an address from a different network ", func() {
			It("should return an error ", func() {
				// A valid utxo segwit address which is not a valid doge address
				addr := multichain.Address("bc1qk6yk2ctcu2pmtxfzhya692h774562vlv2g7dvl")
				decoder := utxochain.NewAddressDecoder(&dogecoin.MainNetParams)
				_, err := decoder.DecodeAddress(addr)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

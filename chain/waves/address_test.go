package waves

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/multichain/api/address"
)

var _ = Describe("Address", func() {
	Context("when decoding and encoding", func() {
		It("should equal itself", func() {
			addr := address.Address("3PMtE788h78hf1DFVPPXKBVa58sjt3QLxwT")

			dec, err := AddressEncodeDecoder{}.DecodeAddress(addr)
			Expect(err).ToNot(HaveOccurred())

			addr2, err := AddressEncodeDecoder{}.EncodeAddress(dec)
			Expect(err).ToNot(HaveOccurred())

			Expect(addr).Should(Equal(addr2))
		})
	})
})

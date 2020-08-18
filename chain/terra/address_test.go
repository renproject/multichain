package terra_test

import (
	"github.com/renproject/multichain/chain/terra"
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Terra", func() {
	Context("when decoding address", func() {
		Context("when decoding Terra address", func() {
			It("should work", func() {
				decoder := terra.NewAddressDecoder()
			
				addrStr := "terra1ztez03dp94y2x55fkhmrvj37ck204geq33msma"
				addr, err := decoder.DecodeAddress(pack.NewString(addrStr))
	
				Expect(err).ToNot(HaveOccurred())
				Expect(addr.AccAddress().String()).Should(Equal(addrStr))
			})
		})
	})
})
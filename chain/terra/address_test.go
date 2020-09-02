package terra_test

import (
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/terra"
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Terra", func() {
	Context("when decoding address", func() {
		Context("when decoding Terra address", func() {
			It("should work", func() {
				decoder := terra.NewAddressDecoder("terra")

				addrStr := "terra1ztez03dp94y2x55fkhmrvj37ck204geq33msma"
				_, err := decoder.DecodeAddress(multichain.Address(pack.NewString(addrStr)))
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})

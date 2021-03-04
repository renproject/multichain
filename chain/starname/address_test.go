package starname_test

import (
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/starname"
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Starname", func() {
	Context("when decoding address", func() {
		Context("when decoding Starname address", func() {
			It("should work", func() {
				decoder := starname.NewAddressDecoder("star")
				addrStr := "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk"
				_, err := decoder.DecodeAddress(multichain.Address(pack.NewString(addrStr)))
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})

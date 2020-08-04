package ontology

import (
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ontology", func() {
	Context("when decoding address", func() {
		Context("when decoding Ontology address", func() {
			It("should work", func() {
				decoder := NewAddressDecoder()
			
				addrStr := "AeeYDwUjR2r5Fm3bqsgwCf8y42wcVJHviQ"
				_, err := decoder.DecodeAddress(pack.NewString(addrStr))
	
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
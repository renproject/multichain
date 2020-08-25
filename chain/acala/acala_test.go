package acala_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/renproject/multichain/chain/acala"
)

var _ = Describe("Substrate client", func() {
	Context("when verifying burns", func() {
		It("should verify a valid burn", func() {
			_, err := acala.NewClient(acala.DefaultClientOptions())
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

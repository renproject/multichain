package solana_test

import (
	"context"

	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/chain/solana"
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Solana", func() {
	Context("...", func() {
		It("...", func() {
			client := solana.NewClient(solana.DefaultClientOptions())
			_, err := client.CallContract(
				context.Background(),
				address.Address(pack.NewString("JBUjNGPApBQ3gw6w2UQPYr1978rkFEGqH1Zs3PZBrHec")),
				pack.NewBytes([]byte{}),
			)
			Expect(err).ToNot(HaveOccurred())
			// Expect(value).To(Equal())
		})
	})
})

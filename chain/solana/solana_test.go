package solana_test

import (
	"context"
	"fmt"

	"github.com/renproject/multichain/chain/solana"
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Solana", func() {
	Context("...", func() {
		It("...", func() {

			client := solana.NewClient(solana.ClientOptions{
				URL: "http://localhost:8899",
			})
			value, err := client.CallContract(
				context.Background(),
				"JBUjNGPApBQ3gw6w2UQPYr1978rkFEGqH1Zs3PZBrHec",
				pack.NewBytes([]byte{}),
				pack.NewBytes([]byte{}).Type(),
			)

			fmt.Printf("account data: %#v", value)

			Expect(err).ToNot(HaveOccurred())
			// Expect(value).To(Equal())
		})
	})
})

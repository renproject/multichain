package bitcoin_test

import (
	"context"

	"github.com/renproject/multichain/chain/bitcoin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gas", func() {
	Context("when estimating bitcoin network fee", func() {
		It("should work", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			client := bitcoin.NewClient(bitcoin.DefaultClientOptions())

			// estimate fee to include tx within 1 block.
			gasEstimator1 := bitcoin.NewGasEstimator(client, 1)
			gasPrice1, _, err := gasEstimator1.EstimateGasPrice(ctx)
			Expect(err).NotTo(HaveOccurred())

			// estimate fee to include tx within 10 blocks.
			gasEstimator2 := bitcoin.NewGasEstimator(client, 10)
			gasPrice2, _, err := gasEstimator2.EstimateGasPrice(ctx)
			Expect(err).NotTo(HaveOccurred())

			// estimate fee to include tx within 100 blocks.
			gasEstimator3 := bitcoin.NewGasEstimator(client, 100)
			gasPrice3, _, err := gasEstimator3.EstimateGasPrice(ctx)
			Expect(err).NotTo(HaveOccurred())

			// expect fees in this order at the very least.
			Expect(gasPrice1.GreaterThanEqual(gasPrice2)).To(BeTrue())
			Expect(gasPrice2.GreaterThanEqual(gasPrice3)).To(BeTrue())
		})
	})
})

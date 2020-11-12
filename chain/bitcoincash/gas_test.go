package bitcoincash_test

import (
	"context"

	"github.com/renproject/multichain/chain/bitcoincash"
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gas", func() {
	Context("when estimating bitcoincash network fee", func() {
		It("should work", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			client := bitcoincash.NewClient(bitcoincash.DefaultClientOptions())

			fallbackGas := uint64(123)
			gasEstimator := bitcoincash.NewGasEstimator(client, pack.NewU256FromUint64(fallbackGas))
			gasPrice, _, err := gasEstimator.EstimateGas(ctx)
			if err != nil {
				Expect(gasPrice).To(Equal(pack.NewU256FromUint64(fallbackGas)))
			} else {
				Expect(gasPrice.Int().Uint64()).To(BeNumerically(">", 0))
			}
		})
	})
})

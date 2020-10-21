package zcash_test

import (
	"context"

	"github.com/renproject/multichain/chain/zcash"
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gas", func() {
	Context("when estimating zcash network fee", func() {
		It("should work", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			client := zcash.NewClient(zcash.DefaultClientOptions())

			// estimate fee to include tx within 1 block.
			fallback1 := uint64(123)
			gasEstimator1 := zcash.NewGasEstimator(client, 1, pack.NewU256FromUint64(fallback1))
			gasPrice1, _, err := gasEstimator1.EstimateGas(ctx)
			if err != nil {
				Expect(gasPrice1).To(Equal(pack.NewU256FromUint64(fallback1)))
			}

			// estimate fee to include tx within 10 blocks.
			fallback2 := uint64(234)
			gasEstimator2 := zcash.NewGasEstimator(client, 10, pack.NewU256FromUint64(fallback2))
			gasPrice2, _, err := gasEstimator2.EstimateGas(ctx)
			if err != nil {
				Expect(gasPrice2).To(Equal(pack.NewU256FromUint64(fallback2)))
			}

			// estimate fee to include tx within 100 blocks.
			fallback3 := uint64(345)
			gasEstimator3 := zcash.NewGasEstimator(client, 100, pack.NewU256FromUint64(fallback3))
			gasPrice3, _, err := gasEstimator3.EstimateGas(ctx)
			if err != nil {
				Expect(gasPrice3).To(Equal(pack.NewU256FromUint64(fallback3)))
			}

			// expect fees in this order at the very least.
			if err == nil {
				Expect(gasPrice1.GreaterThanEqual(gasPrice2)).To(BeTrue())
				Expect(gasPrice2.GreaterThanEqual(gasPrice3)).To(BeTrue())
			}
		})
	})
})

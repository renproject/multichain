package cosmos_test

import (
	"context"
	"testing/quick"

	"github.com/renproject/multichain/chain/cosmos"
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gas", func() {
	Context("when estimating gas parameters", func() {
		It("should work", func() {
			f := func(gasPerByte pack.U256) bool {
				gasEstimator := cosmos.NewGasEstimator(gasPerByte)
				gasPrice, _, err := gasEstimator.EstimateGas(context.Background())
				Expect(err).NotTo(HaveOccurred())
				Expect(gasPrice).To(Equal(gasPerByte))
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})
})

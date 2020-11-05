package cosmos_test

import (
	"context"
	"math/rand"
	"testing/quick"
	"time"

	"github.com/renproject/multichain/chain/cosmos"
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gas", func() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	Context("when estimating gas parameters", func() {
		It("should work", func() {
			f := func() bool {
				gasPriceMicro := r.Float64()
				gasEstimator := cosmos.NewGasEstimator(gasPriceMicro)
				gasPricePico, _, err := gasEstimator.EstimateGas(context.Background())
				Expect(err).NotTo(HaveOccurred())
				expectedGasPrice := pack.NewU256FromUint64(uint64(gasPriceMicro * 1000000))
				Expect(gasPricePico).To(Equal(expectedGasPrice))
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})
})

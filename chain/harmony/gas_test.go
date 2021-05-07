package harmony_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/multichain/chain/harmony"
	"github.com/renproject/pack"
)

var _ = Describe("Gas", func() {
	Context("when estimating gas", func() {
		It("should work without errors", func() {
			e := harmony.Estimator{}
			gas, err := e.EstimateGasPrice(context.TODO())
			Expect(err).NotTo(HaveOccurred())
			Expect(gas).To(Equal(pack.NewU256FromU64(pack.NewU64(1))))
		})
	})
})

package solana_test

import (
	"github.com/renproject/multichain/chain/solana"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Solana", func() {
	Context("...", func() {
		It("...", func() {
			_ = solana.NewClient(solana.DefaultClientOptions())
		})
	})
})

package digibyte_test

import (
	"context"
	"encoding/hex"

	"github.com/renproject/multichain/chain/digibyte"
	
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DigiByte", func() {
	Context("when decoding an address", func() {
		It("should work without errors", func() {
			rt := runtime.NewRuntime(nil, nil, nil, nil, nil, nil, nil, nil)
			// This test needs some work
			_, err := rt.BitcoinDecodeAddress(multichain.DigiByte, pack.NewString("DCo1dbnnwWB4cucwSduXMdTV1tDErZHNfx"))
			Expect(err).To.Not(HaveOccurred())
		})
	})
})

package digibyte_test

import (
	"fmt"

	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/bitcoin"
	"github.com/renproject/multichain/chain/digibyte"
	"github.com/renproject/multichain/runtime"
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DigiByte", func() {

	Context("when creating an address", func() {
		It("should work without errors", func() {
			rt := runtime.NewRuntime(runtime.AddressDecoders{
				multichain.DigiByte: bitcoin.NewAddressDecoder(&digibyte.MainNetParams),
			}, nil, nil, nil, nil, nil)

			// Encode PKH into DigiByte Address
			val, err := rt.DecodeAddress(multichain.DigiByte, pack.NewString("DBLsEv4FdFPGrMWzcagDQvoKgUL2CikhMf"))
			Expect(err).NotTo(HaveOccurred())

			fmt.Println(val)
		})
	})
})

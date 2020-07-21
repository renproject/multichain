package digibyte_test

import (
	"fmt"

	"github.com/renproject/multichain"
	"github.com/renproject/multichain/runtime"
	"github.com/renproject/multichain/chain/digibyte"
	"github.com/renproject/pack"
	
	"github.com/btcsuite/btcd/chaincfg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DigiByte", func() {

	Context("when creating an address", func() {
		It("should work without errors", func() {
			digibyte.Init()
			
			rt := runtime.NewRuntime(nil, runtime.BitcoinCompatAddressDecoders{
				multichain.DigiByte: digibyte.NewAddressDecoder(&chaincfg.MainNetParams),
			}, nil, nil, nil, nil, nil, nil)

			// Encode PKH into DigiByte Address
			val, err := rt.BitcoinDecodeAddress(multichain.DigiByte, pack.NewString("DBLsEv4FdFPGrMWzcagDQvoKgUL2CikhMf"))
			Expect(err).NotTo(HaveOccurred())

			fmt.Println(val)
		})
	})
})

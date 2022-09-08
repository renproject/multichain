package digibyte_test

import (
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/chain/digibyte"
	"github.com/renproject/multichain/chain/utxochain"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DigiByte", func() {
	Context("when decoding an address", func() {
		It("should work without errors", func() {
			_, err := utxochain.NewAddressDecoder(&digibyte.MainNetParams).DecodeAddress(address.Address("DBLsEv4FdFPGrMWzcagDQvoKgUL2CikhMf"))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

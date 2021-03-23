package terra_test

import (
	"fmt"

	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/terra"
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Terra", func() {
	Context("when decoding address", func() {
		decoder := terra.NewAddressDecoder()
		Context("when decoding a valid address", func() {
			It("should work", func() {
				addrStr := "terra1ztez03dp94y2x55fkhmrvj37ck204geq33msma"
				_, err := decoder.DecodeAddress(multichain.Address(pack.NewString(addrStr)))
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Context("when decoding an address with invalid prefix", func() {
			It("should fail", func() {
				addrStr := "cosmosztez03dp94y2x55fkhmrvj37ck204geq33msma"
				_, err := decoder.DecodeAddress(multichain.Address(pack.NewString(addrStr)))
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when decoding an invalid address", func() {
			It("should fail", func() {
				addrStr := "terra1ztez03dp94y2x55fkhmrvj37ck204geq33msm"
				_, err := decoder.DecodeAddress(multichain.Address(pack.NewString(addrStr)))
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when decoding an empty address", func() {
			It("should fail", func() {
				addrStr := ""
				_, err := decoder.DecodeAddress(multichain.Address(pack.NewString(addrStr)))
				Expect(err).Should(MatchError(fmt.Errorf("unexpected address length: want=20, got=0")))
			})
		})
	})
})

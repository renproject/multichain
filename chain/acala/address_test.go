package acala_test

import (
	"github.com/centrifuge/go-substrate-rpc-client/signature"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/chain/acala"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Address", func() {
	addrEncodeDecoder := acala.NewAddressEncodeDecoder(acala.AddressTypeDefault)

	Context("when encoding raw address", func() {
		It("should match the human-readable address", func() {
			addr, err := addrEncodeDecoder.EncodeAddress(address.RawAddress(signature.TestKeyringPairAlice.PublicKey))
			Expect(err).NotTo(HaveOccurred())
			Expect(addr).To(Equal(address.Address(signature.TestKeyringPairAlice.Address)))
			rawAddr, err := addrEncodeDecoder.DecodeAddress(addr)
			Expect(err).NotTo(HaveOccurred())
			Expect(rawAddr).To(Equal(address.RawAddress(signature.TestKeyringPairAlice.PublicKey)))
		})
	})

	Context("when decoding human-readable address", func() {
		It("should match the raw address", func() {
			rawAddr, err := addrEncodeDecoder.DecodeAddress(address.Address(signature.TestKeyringPairAlice.Address))
			Expect(err).NotTo(HaveOccurred())
			Expect(rawAddr).To(Equal(address.RawAddress(signature.TestKeyringPairAlice.PublicKey)))
			addr, err := addrEncodeDecoder.EncodeAddress(rawAddr)
			Expect(err).NotTo(HaveOccurred())
			Expect(addr).To(Equal(address.Address(signature.TestKeyringPairAlice.Address)))
		})
	})
})

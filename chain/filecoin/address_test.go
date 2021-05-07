package filecoin_test

import (
	"math/rand"
	"testing/quick"
	"time"

	filaddress "github.com/filecoin-project/go-address"
	"github.com/multiformats/go-varint"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/chain/filecoin"
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Address", func() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	encoderDecoder := filecoin.NewAddressEncodeDecoder()

	Context("when encoding and decoding", func() {
		Context("for ID protocol", func() {
			It("should behave correctly without errors", func() {
				f := func() bool {
					x := varint.ToUvarint(uint64(r.Int63()))
					x = append([]byte{byte(filaddress.ID)}, x...)

					rawAddr := address.RawAddress(pack.NewBytes(x[:]))
					addr, err := encoderDecoder.AddressEncoder.EncodeAddress(rawAddr)
					Expect(err).ToNot(HaveOccurred())

					decodedAddr, err := encoderDecoder.AddressDecoder.DecodeAddress(addr)
					Expect(err).ToNot(HaveOccurred())
					Expect(decodedAddr).To(Equal(rawAddr))
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})

		Context("for Sepc protocol", func() {
			It("should behave correctly without errors", func() {
				f := func(x [filaddress.PayloadHashLength]byte) bool {
					y := append([]byte{byte(filaddress.SECP256K1)}, x[:]...)

					rawAddr := address.RawAddress(pack.NewBytes(y[:]))
					addr, err := encoderDecoder.AddressEncoder.EncodeAddress(rawAddr)
					Expect(err).ToNot(HaveOccurred())

					decodedAddr, err := encoderDecoder.AddressDecoder.DecodeAddress(addr)
					Expect(err).ToNot(HaveOccurred())
					Expect(decodedAddr).To(Equal(rawAddr))
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})

		Context("for Actor protocol", func() {
			It("should behave correctly without errors", func() {
				f := func(x [filaddress.PayloadHashLength]byte) bool {
					y := append([]byte{byte(filaddress.Actor)}, x[:]...)

					rawAddr := address.RawAddress(pack.NewBytes(y[:]))
					addr, err := encoderDecoder.AddressEncoder.EncodeAddress(rawAddr)
					Expect(err).ToNot(HaveOccurred())

					decodedAddr, err := encoderDecoder.AddressDecoder.DecodeAddress(addr)
					Expect(err).ToNot(HaveOccurred())
					Expect(decodedAddr).To(Equal(rawAddr))
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})

		Context("for BLS protocol", func() {
			It("should behave correctly without errors", func() {
				f := func(x [filaddress.BlsPublicKeyBytes]byte) bool {
					y := append([]byte{byte(filaddress.BLS)}, x[:]...)

					rawAddr := address.RawAddress(pack.NewBytes(y[:]))
					addr, err := encoderDecoder.AddressEncoder.EncodeAddress(rawAddr)
					Expect(err).ToNot(HaveOccurred())

					decodedAddr, err := encoderDecoder.AddressDecoder.DecodeAddress(addr)
					Expect(err).ToNot(HaveOccurred())
					Expect(decodedAddr).To(Equal(rawAddr))
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})
	})
})

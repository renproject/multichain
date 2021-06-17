package bsc_test

import (
	"encoding/hex"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/multichain/chain/bsc"
	"github.com/renproject/surge"
	"testing/quick"
)

var _ = Describe("Address", func() {
	Context("when unmarshaling and unmarshaling", func() {
		It("should equal itself", func() {
			f := func(x [20]byte) bool {
				addr := bsc.Address(x)
				Expect(addr.SizeHint()).To(Equal(20))

				bytes, err := surge.ToBinary(addr)
				Expect(err).ToNot(HaveOccurred())

				var newAddr bsc.Address
				err = surge.FromBinary(&newAddr, bytes)
				Expect(err).ToNot(HaveOccurred())

				Expect(addr).To(Equal(newAddr))
				return true
			}

			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when unmarshaling and unmarshaling to/from JSON", func() {
		It("should equal itself", func() {
			f := func(x [20]byte) bool {
				addr := bsc.Address(x)

				bytes, err := json.Marshal(addr)
				Expect(err).ToNot(HaveOccurred())

				var newAddr bsc.Address
				err = json.Unmarshal(bytes, &newAddr)
				Expect(err).ToNot(HaveOccurred())

				Expect(addr).To(Equal(newAddr))
				return true
			}

			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the address is invalid hex", func() {
			It("should return an error", func() {
				f := func(x [40]byte) bool {
					bytes, err := json.Marshal(string(x[:]))
					Expect(err).ToNot(HaveOccurred())

					var newAddr bsc.Address
					err = json.Unmarshal(bytes, &newAddr)
					Expect(err).To(HaveOccurred())
					return true
				}

				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the address is invalid length", func() {
			It("should return an error", func() {
				f := func(x [10]byte) bool {
					addr := hex.EncodeToString(x[:])
					bytes, err := json.Marshal(addr)
					Expect(err).ToNot(HaveOccurred())

					var newAddr bsc.Address
					err = json.Unmarshal(bytes, &newAddr)
					Expect(err).To(HaveOccurred())
					return true
				}

				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when unmarshalling random data", func() {
		It("should not panic", func() {
			f := func(x []byte) bool {
				var addr bsc.Address
				Expect(func() { addr.Unmarshal(x, surge.MaxBytes) }).ToNot(Panic())
				Expect(func() { json.Unmarshal(x, &addr) }).ToNot(Panic())
				return true
			}

			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

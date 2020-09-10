package ethereum_test

import (
	"encoding/hex"
	"fmt"
	"math"
	"testing/quick"

	"github.com/renproject/multichain/chain/ethereum"
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Encoding", func() {
	Context("when encoding bytes", func() {
		It("should return the correct result", func() {
			f := func(x []byte) bool {
				arg := pack.NewBytes(x)

				resBytes := ethereum.Encode(arg)
				resString := hex.EncodeToString(resBytes)

				expectedBytes := make([]byte, int(math.Ceil(float64(len(x))/32)*32))
				copy(expectedBytes, x)
				// Note: since the first parameter has a dynamic length, the
				// first 32 bytes instead contain a pointer to the data.
				expectedString := fmt.Sprintf("%064x", 32) + fmt.Sprintf("%064x", len(x)) + hex.EncodeToString(expectedBytes)

				Expect(resString).To(Equal(expectedString))
				return true
			}

			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when encoding 32 bytes", func() {
		It("should return the correct result", func() {
			f := func(x [32]byte) bool {
				arg := pack.NewBytes32(x)

				resBytes := ethereum.Encode(arg)
				resString := hex.EncodeToString(resBytes)
				expectedString := hex.EncodeToString(x[:])

				Expect(resString).To(Equal(expectedString))
				return true
			}

			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when encoding 8-bit unsigned integers", func() {
		It("should return the correct result", func() {
			f := func(x uint8) bool {
				arg := pack.NewU8(x)

				resBytes := ethereum.Encode(arg)
				resString := hex.EncodeToString(resBytes)
				expectedString := fmt.Sprintf("%064x", x)

				Expect(resString).To(Equal(expectedString))
				return true
			}

			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when encoding 16-bit unsigned integers", func() {
		It("should return the correct result", func() {
			f := func(x uint16) bool {
				arg := pack.NewU16(x)

				resBytes := ethereum.Encode(arg)
				resString := hex.EncodeToString(resBytes)
				expectedString := fmt.Sprintf("%064x", x)

				Expect(resString).To(Equal(expectedString))
				return true
			}

			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when encoding 32-bit unsigned integers", func() {
		It("should return the correct result", func() {
			f := func(x uint32) bool {
				arg := pack.NewU32(x)

				resBytes := ethereum.Encode(arg)
				resString := hex.EncodeToString(resBytes)
				expectedString := fmt.Sprintf("%064x", x)

				Expect(resString).To(Equal(expectedString))
				return true
			}

			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when encoding 64-bit unsigned integers", func() {
		It("should return the correct result", func() {
			f := func(x uint64) bool {
				arg := pack.NewU64(x)

				resBytes := ethereum.Encode(arg)
				resString := hex.EncodeToString(resBytes)
				expectedString := fmt.Sprintf("%064x", x)

				Expect(resString).To(Equal(expectedString))
				return true
			}

			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when encoding 128-bit unsigned integers", func() {
		It("should return the correct result", func() {
			f := func(x [16]byte) bool {
				arg := pack.NewU128(x)

				resBytes := ethereum.Encode(arg)
				resString := hex.EncodeToString(resBytes)
				expectedString := fmt.Sprintf("%064x", x)

				Expect(resString).To(Equal(expectedString))
				return true
			}

			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when encoding 256-bit unsigned integers", func() {
		It("should return the correct result", func() {
			f := func(x [32]byte) bool {
				arg := pack.NewU256(x)

				resBytes := ethereum.Encode(arg)
				resString := hex.EncodeToString(resBytes)
				expectedString := fmt.Sprintf("%064x", x)

				Expect(resString).To(Equal(expectedString))
				return true
			}

			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when encoding Ethereum addresses", func() {
		It("should return the correct result", func() {
			f := func(x [20]byte) bool {
				arg := ethereum.Address(x)

				resBytes := ethereum.Encode(arg)
				resString := hex.EncodeToString(resBytes)

				expectedBytes := make([]byte, 32)
				copy(expectedBytes, x[:])
				expectedString := hex.EncodeToString(expectedBytes)

				Expect(resString).To(Equal(expectedString))
				return true
			}

			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when encoding an unsupported type", func() {
		It("should panic", func() {
			f := func(x bool) bool {
				arg := pack.NewBool(x)
				Expect(func() { ethereum.Encode(arg) }).To(Panic())
				return true
			}

			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	type testCase struct {
		addr   string
		amount uint64
		hash   string
		result string
	}

	testCases := []testCase{
		{
			addr:   "797522Fb74d42bB9fbF6b76dEa24D01A538d5D66",
			amount: 10000,
			hash:   "702826c3977ee72158db2ce1fb758075ee2799db65fb27b5d0952f860a8084ed",
			result: "797522fb74d42bb9fbf6b76dea24d01a538d5d660000000000000000000000000000000000000000000000000000000000000000000000000000000000002710702826c3977ee72158db2ce1fb758075ee2799db65fb27b5d0952f860a8084ed",
		},
		{
			addr:   "58afb504ef2444a267b8c7ce57279417f1377ceb",
			amount: 50000000000000000,
			hash:   "dabff9ceb1b3dabb696d143326fdb98a8c7deb260e65d08a294b16659d573f93",
			result: "58afb504ef2444a267b8c7ce57279417f1377ceb00000000000000000000000000000000000000000000000000000000000000000000000000b1a2bc2ec50000dabff9ceb1b3dabb696d143326fdb98a8c7deb260e65d08a294b16659d573f93",
		},
		{
			addr:   "0000000000000000000000000000000000000000",
			amount: 0,
			hash:   "0000000000000000000000000000000000000000000000000000000000000000",
			result: "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
	}

	DescribeTable("when encoding args",
		func(test testCase) {
			addrBytes, err := hex.DecodeString(test.addr)
			Expect(err).ToNot(HaveOccurred())

			var addr ethereum.Address
			copy(addr[:], addrBytes)

			hashBytes32 := [32]byte{}
			hashBytes, err := hex.DecodeString(test.hash)
			Expect(err).ToNot(HaveOccurred())
			copy(hashBytes32[:], hashBytes)

			args := []interface{}{
				addr,
				pack.NewU64(test.amount),
				pack.NewBytes32(hashBytes32),
			}
			result := ethereum.Encode(args...)
			Expect(hex.EncodeToString(result)).To(Equal(test.result))
		},

		Entry("should return the same result as solidity for small transactions", testCases[0]),
		Entry("should return the same result as solidity for large transactions", testCases[1]),
		Entry("should return the same result as solidity for empty transactions", testCases[2]),
	)
})

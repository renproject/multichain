package ethereumcompat_test

import (
	"encoding/hex"
	"fmt"
	"math"
	"testing/quick"

	"github.com/renproject/multichain/compat/ethereumcompat"
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

				resBytes := ethereumcompat.Encode(arg)
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

				resBytes := ethereumcompat.Encode(arg)
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

				resBytes := ethereumcompat.Encode(arg)
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

				resBytes := ethereumcompat.Encode(arg)
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

				resBytes := ethereumcompat.Encode(arg)
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

				resBytes := ethereumcompat.Encode(arg)
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

				resBytes := ethereumcompat.Encode(arg)
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

				resBytes := ethereumcompat.Encode(arg)
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
				arg := ethereumcompat.Address(x)

				resBytes := ethereumcompat.Encode(arg)
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
				Expect(func() { ethereumcompat.Encode(arg) }).To(Panic())
				return true
			}

			err := quick.Check(f, nil)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when encoding complex calldata", func() {
		It("should return the expected bytes", func() {
			phashHex := "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
			phashBytes, err := hex.DecodeString(phashHex)
			Expect(err).ToNot(HaveOccurred())
			var phash pack.Bytes32
			copy(phash[:], phashBytes)
			tokenHex := "0x675000eed287586fcb53d676f4ab1c15e8be314d"
			tokenAddr, err := ethereumcompat.NewAddressFromHex(tokenHex)
			Expect(err).ToNot(HaveOccurred())
			toHex := "0xc7ddb84d0d8f70dbd8453f88a5967681cd3c9830"
			toAddr, err := ethereumcompat.NewAddressFromHex(toHex)
			Expect(err).ToNot(HaveOccurred())
			nhashHex := "5f2515e44d1bf07d591b78f5c896c87b7dd7441fbf3395f3f0c88578ef7bdc04"
			nhashBytes, err := hex.DecodeString(nhashHex)
			Expect(err).ToNot(HaveOccurred())
			var nhash pack.Bytes32
			copy(nhash[:], nhashBytes)
			args := []interface{}{
				phash,
				pack.NewU256FromU64(pack.NewU64(9711004)),
				ethereumcompat.Address(tokenAddr),
				ethereumcompat.Address(toAddr),
				nhash,
			}
			result := ethereumcompat.Encode(args...)
			Expect(hex.EncodeToString(result)).To(Equal("c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a4700000000000000000000000000000000000000000000000000000000000942d9c000000000000000000000000675000eed287586fcb53d676f4ab1c15e8be314d000000000000000000000000c7ddb84d0d8f70dbd8453f88a5967681cd3c98305f2515e44d1bf07d591b78f5c896c87b7dd7441fbf3395f3f0c88578ef7bdc04"))
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

			var addr ethereumcompat.Address
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
			result := ethereumcompat.Encode(args...)
			Expect(hex.EncodeToString(result)).To(Equal(test.result))
		},

		Entry("should return the same result as solidity for small transactions", testCases[0]),
		Entry("should return the same result as solidity for large transactions", testCases[1]),
		Entry("should return the same result as solidity for empty transactions", testCases[2]),
	)
})

package solana_test

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/renproject/multichain/chain/solana"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Solana FFI", func() {
	Context("FFI", func() {
		It("should say hello", func() {
			Expect(solana.Hello("Stranger")).To(Equal("Hello Stranger!"))
			Expect(solana.Hello("Friend")).To(Equal("Hello Friend!"))
		})

		It("should create a unique pubkey", func() {
			key1 := solana.UniquePubkey()
			key2 := solana.UniquePubkey()
			Expect(key1).NotTo(Equal(key2))
		})
	})

	Context("Program Derived Address", func() {
		It("should correctly compute 1", func() {
			program := "6kAHanNCT1LKFoMn3fBdyvJuvHLcWhLpJbTpbHpqRiG4"
			seeds := []byte("RenBridgeState")
			programDerivedAddress := solana.ProgramDerivedAddress(seeds, program)
			expectedDerivedAddress := base58.Decode("7gMf4XXqunXaagnMVf8c3KSKnANTjhDvn2HgVTxMb4ZD")
			Expect(programDerivedAddress[:]).To(Equal(expectedDerivedAddress))
		})

		It("should correctly compute 2", func() {
			program := "6kAHanNCT1LKFoMn3fBdyvJuvHLcWhLpJbTpbHpqRiG4"
			selector := "BTC/toSolana"
			selectorHash := crypto.Keccak256([]byte(selector))
			programDerivedAddress := solana.ProgramDerivedAddress(selectorHash, program)
			expectedDerivedAddress := base58.Decode("6SPY5x3tmjLZ9SWcZFKhwpANrhYJagNNF4Sa4LAwtbCn")
			Expect(programDerivedAddress[:]).To(Equal(expectedDerivedAddress))
		})
	})
})

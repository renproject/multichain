package solana_test

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/chain/solana"
	"github.com/renproject/pack"
	"github.com/renproject/solana-ffi/cgo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Solana FFI", func() {
	Context("FFI", func() {
		It("should create a unique pubkey", func() {
			key1 := solana.UniquePubkey()
			key2 := solana.UniquePubkey()
			Expect(key1).NotTo(Equal(key2))
		})
	})

	Context("Program Derived Address", func() {
		It("should correctly compute 1", func() {
			program := address.Address("6kAHanNCT1LKFoMn3fBdyvJuvHLcWhLpJbTpbHpqRiG4")
			seeds := []byte("GatewayState")
			programDerivedAddress := solana.ProgramDerivedAddress(pack.Bytes(seeds), program)
			expectedDerivedAddress := address.Address("APthNc29MGRJRkKahDRNrSNA2o1e8p6aFAJNRV8ZdJaV")
			Expect(programDerivedAddress[:]).To(Equal(expectedDerivedAddress))
		})

		It("should correctly compute 2", func() {
			program := address.Address("6kAHanNCT1LKFoMn3fBdyvJuvHLcWhLpJbTpbHpqRiG4")
			selector := "BTC/toSolana"
			selectorHash := crypto.Keccak256([]byte(selector))
			programDerivedAddress := solana.ProgramDerivedAddress(pack.Bytes(selectorHash), program)
			expectedDerivedAddress := address.Address("6SPY5x3tmjLZ9SWcZFKhwpANrhYJagNNF4Sa4LAwtbCn")
			Expect(programDerivedAddress[:]).To(Equal(expectedDerivedAddress))
		})
	})

	Context("Associated Token Account", func() {
		It("should correctly calculate", func() {
			walletAddress := "fYq3qkHoVogcPnkxFWAwiJGJs29Xtg4FZ6xcAHWd51w"
			selector := "BTC/toSolana"
			assTokenAccount := cgo.AssociatedTokenAccount(walletAddress, selector)
			expectedAssTokenAccount := "GxMKqib75YSD5RegZP8A7ZkSv8uBFmfNsNXzGptBdqdo"
			Expect(assTokenAccount).To(Equal(expectedAssTokenAccount))
		})
	})
})

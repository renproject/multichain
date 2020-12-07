package solana_test

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/bitcoin"
	"github.com/renproject/multichain/chain/solana"
	"github.com/renproject/pack"
	"github.com/renproject/surge"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Solana", func() {
	Context("fetch burn log", func() {
		It("should be able to fetch", func() {
			client := solana.NewClient(solana.DefaultClientOptions())

			burnCount := pack.NewU64(uint64(1))
			contractCallInput := solana.BurnCallContractInput{Nonce: burnCount}
			calldata, err := surge.ToBinary(contractCallInput)
			Expect(err).ToNot(HaveOccurred())
			program := multichain.Address("9TaQuUfNMC5rFvdtzhHPk84WaFH3SFnweZn4tw9RriDP")

			burnLogBytes, err := client.CallContract(context.Background(), program, calldata)
			Expect(err).NotTo(HaveOccurred())

			burnLog := solana.BurnCallContractOutput{}
			err = surge.FromBinary(&burnLog, burnLogBytes)
			Expect(err).NotTo(HaveOccurred())
			Expect(burnLog.Amount).To(Equal(pack.U64(2000000000)))
			expectedRecipient := multichain.Address("mwjUmhAW68zCtgZpW5b1xD5g7MZew6xPV4")
			bitcoinAddrEncodeDecoder := bitcoin.NewAddressEncodeDecoder(&chaincfg.RegressionNetParams)
			rawAddr, err := bitcoinAddrEncodeDecoder.DecodeAddress(expectedRecipient)
			Expect(err).NotTo(HaveOccurred())
			Expect(burnLog.Recipient).To(Equal(rawAddr))
		})
	})
})

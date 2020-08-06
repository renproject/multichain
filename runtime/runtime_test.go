package runtime_test

import (
	"context"

	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/bitcoin"
	"github.com/renproject/multichain/compat/bitcoincompat"
	"github.com/renproject/multichain/runtime"
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bitcoin-compat", func() {
	Context("when querying outputs", func() {
		Context("when the chain is not supported", func() {
			It("should return an error", func() {
				rt := runtime.NewRuntime(nil, nil, nil, nil, nil)
				_, err := rt.BitcoinOutput(context.Background(), multichain.Bitcoin, multichain.BTC, bitcoincompat.Outpoint{})
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the chain is supported", func() {
			It("should return the output", func() {
				// TODO: Implement.
			})
		})
	})

	Context("when querying gas-per-byte", func() {
		Context("when the chain is not supported", func() {
			It("should return an error", func() {
				rt := runtime.NewRuntime(nil, nil, nil, nil, nil)
				_, err := rt.BitcoinGasPerByte(context.Background(), multichain.Bitcoin)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the chain is supported", func() {
			It("should return the output", func() {
				rt := runtime.NewRuntime(nil, nil, nil, runtime.BitcoinCompatGasEstimators{
					multichain.Bitcoin: bitcoincompat.NewGasEstimator(pack.NewU64(10000)),
				}, nil)
				gasPerByte, err := rt.BitcoinGasPerByte(context.Background(), multichain.Bitcoin)
				Expect(err).ToNot(HaveOccurred())
				Expect(gasPerByte).To(Equal(pack.NewU64(10000)))
			})
		})
	})

	Context("when building transactions", func() {
		Context("when the chain is not supported", func() {
			It("should return an error", func() {
				rt := runtime.NewRuntime(nil, nil, nil, nil, nil)
				_, err := rt.BitcoinBuildTx(context.Background(), multichain.Bitcoin, multichain.BTC, []bitcoincompat.Input{}, []bitcoincompat.Recipient{})
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the chain is supported", func() {
			It("should return the output", func() {
				// TODO: Implement.
			})
		})
	})

	Context("when submitting transcations", func() {
		Context("when the chain is not supported", func() {
			It("should return an error", func() {
				rt := runtime.NewRuntime(nil, nil, nil, nil, nil)
				_, err := rt.BitcoinSubmitTx(context.Background(), multichain.Bitcoin, &bitcoin.Tx{})
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the chain is supported", func() {
			It("should return the output", func() {
				// TODO: Implement.
			})
		})
	})
})

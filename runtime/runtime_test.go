package runtime_test

import (
	"context"
	"encoding/hex"

	"github.com/btcsuite/btcd/chaincfg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/bitcoin"
	"github.com/renproject/multichain/compat/bitcoincompat"
	"github.com/renproject/multichain/compat/ethereumcompat"
	"github.com/renproject/multichain/compat/substratecompat"
	"github.com/renproject/multichain/runtime"
	"github.com/renproject/pack"
)

var _ = Describe("Bitcoin-compat", func() {
	Context("when decoding addresses", func() {
		Context("when the chain is not supported", func() {
			It("should return an error", func() {
				rt := runtime.NewRuntime(nil, nil, nil, nil, nil, nil, nil, nil)
				_, err := rt.BitcoinDecodeAddress(multichain.Bitcoin, pack.NewString("mwjUmhAW68zCtgZpW5b1xD5g7MZew6xPV4"))
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the chain is supported", func() {
			It("should return an address", func() {
				rt := runtime.NewRuntime(nil, runtime.BitcoinCompatAddressDecoders{
					multichain.Bitcoin: bitcoin.NewAddressDecoder(&chaincfg.RegressionNetParams),
				}, nil, nil, nil, nil, nil, nil)
				addr, err := rt.BitcoinDecodeAddress(multichain.Bitcoin, pack.NewString("mwjUmhAW68zCtgZpW5b1xD5g7MZew6xPV4"))
				Expect(err).ToNot(HaveOccurred())
				Expect(addr.EncodeAddress()).To(Equal("mwjUmhAW68zCtgZpW5b1xD5g7MZew6xPV4"))
			})
		})
	})

	Context("when querying outputs", func() {
		Context("when the chain is not supported", func() {
			It("should return an error", func() {
				rt := runtime.NewRuntime(nil, nil, nil, nil, nil, nil, nil, nil)
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
				rt := runtime.NewRuntime(nil, nil, nil, nil, nil, nil, nil, nil)
				_, err := rt.BitcoinGasPerByte(context.Background(), multichain.Bitcoin)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the chain is supported", func() {
			It("should return the output", func() {
				rt := runtime.NewRuntime(nil, nil, nil, runtime.BitcoinCompatGasEstimators{
					multichain.Bitcoin: bitcoincompat.NewGasEstimator(pack.NewU64(10000)),
				}, nil, nil, nil, nil)
				gasPerByte, err := rt.BitcoinGasPerByte(context.Background(), multichain.Bitcoin)
				Expect(err).ToNot(HaveOccurred())
				Expect(gasPerByte).To(Equal(pack.NewU64(10000)))
			})
		})
	})

	Context("when building transactions", func() {
		Context("when the chain is not supported", func() {
			It("should return an error", func() {
				rt := runtime.NewRuntime(nil, nil, nil, nil, nil, nil, nil, nil)
				_, err := rt.BitcoinBuildTx(context.Background(), multichain.Bitcoin, multichain.BTC, []bitcoincompat.Output{}, []bitcoincompat.Recipient{})
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
				rt := runtime.NewRuntime(nil, nil, nil, nil, nil, nil, nil, nil)
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

var _ = Describe("Ethereum-compat", func() {
	Context("when decoding addresses", func() {
		Context("when the chain is not supported", func() {
			It("should return an error", func() {
				rt := runtime.NewRuntime(nil, nil, nil, nil, nil, nil, nil, nil)
				_, err := rt.EthereumDecodeAddress(multichain.Ethereum, pack.NewString("mwjUmhAW68zCtgZpW5b1xD5g7MZew6xPV4"))
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the chain is supported", func() {
			It("should return an address", func() {
				rt := runtime.NewRuntime(nil, nil, nil, nil, nil, runtime.EthereumCompatAddressDecoders{
					multichain.Ethereum: ethereumcompat.NewAddressDecoder(),
				}, nil, nil)
				addr, err := rt.EthereumDecodeAddress(multichain.Ethereum, pack.NewString("0x0123456789012345678901234567890123456789"))
				Expect(err).ToNot(HaveOccurred())
				Expect(addr.String()).To(Equal("0123456789012345678901234567890123456789"))
			})
		})
	})

	Context("when querying burn events", func() {
		Context("when the chain is not supported", func() {
			It("should return an error", func() {
				rt := runtime.NewRuntime(nil, nil, nil, nil, nil, nil, nil, nil)
				_, _, err := rt.EthereumBurnEvent(context.Background(), multichain.Ethereum, multichain.BTC, pack.Bytes32([32]byte{}))
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the chain is supported", func() {
			It("should return the burn event", func() {
				// TODO: Implement.
			})
		})
	})
})

var _ = Describe("Substrate-compat", func() {
	Context("when decoding addresses", func() {
		Context("when the chain is not supported", func() {
			It("should return an error", func() {
				rt := runtime.NewRuntime(nil, nil, nil, nil, nil, nil, nil, nil)
				_, err := rt.SubstrateDecodeAddress(multichain.Acala, pack.NewString("mwjUmhAW68zCtgZpW5b1xD5g7MZew6xPV4"))
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the chain is supported", func() {
			It("should return an address", func() {
				rt := runtime.NewRuntime(nil, nil, nil, nil, nil, nil, nil, runtime.SubstrateCompatAddressDecoders{
					multichain.Acala: substratecompat.NewAddressDecoder(),
				})
				addr, err := rt.SubstrateDecodeAddress(multichain.Acala, pack.NewString("5Hp67orXVehS6dnHJbcuaZjGRbJhg5YqXvR6BgJPE1JZtQtP"))
				Expect(err).ToNot(HaveOccurred())
				Expect(hex.EncodeToString(addr[:])).To(Equal("2afe43a82f10bb2eff7e730d994be2c6f2db0637a3b6e6eb260d8c796b19bb300bf5f8"))
			})
		})
	})

	Context("when querying burn events", func() {
		Context("when the chain is not supported", func() {
			It("should return an error", func() {
				rt := runtime.NewRuntime(nil, nil, nil, nil, nil, nil, nil, nil)
				_, _, err := rt.SubstrateBurnEvent(context.Background(), multichain.Acala, multichain.BTC, pack.Bytes32([32]byte{}))
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the chain is supported", func() {
			It("should return the burn event", func() {
				// TODO: Implement.
			})
		})
	})
})

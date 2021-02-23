package zcash_test

import (
	"bytes"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/renproject/id"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/chain/zcash"
)

var _ = Describe("Zcash Address", func() {
	Context("address", func() {
		addrEncodeDecoder := zcash.NewAddressEncodeDecoder(&zcash.RegressionNetParams)

		It("addr pub key hash", func() {
			pk := id.NewPrivKey()
			wif, err := btcutil.NewWIF((*btcec.PrivateKey)(pk), zcash.RegressionNetParams.Params, true)
			Expect(err).NotTo(HaveOccurred())
			addrPubKeyHash, err := zcash.NewAddressPubKeyHash(btcutil.Hash160(wif.PrivKey.PubKey().SerializeUncompressed()), &zcash.RegressionNetParams)
			Expect(err).NotTo(HaveOccurred())
			addr := address.Address(addrPubKeyHash.EncodeAddress())

			decodedRawAddr, err := addrEncodeDecoder.DecodeAddress(addr)
			Expect(err).NotTo(HaveOccurred())
			encodedAddr, err := addrEncodeDecoder.EncodeAddress(decodedRawAddr)
			Expect(err).NotTo(HaveOccurred())
			Expect(encodedAddr).To(Equal(addr))
		})

		It("addr script hash", func() {
			script := make([]byte, rand.Intn(100))
			rand.Read(script)
			addrScriptHash, err := zcash.NewAddressScriptHash(script, &zcash.RegressionNetParams)
			Expect(err).NotTo(HaveOccurred())
			addr := address.Address(addrScriptHash.EncodeAddress())

			decodedRawAddr, err := addrEncodeDecoder.DecodeAddress(addr)
			Expect(err).NotTo(HaveOccurred())
			encodedAddr, err := addrEncodeDecoder.EncodeAddress(decodedRawAddr)
			Expect(err).NotTo(HaveOccurred())
			Expect(encodedAddr).To(Equal(addr))
		})
	})

	Context("AddressEncodeDecoder", func() {
		It("should give an error when decoding address on different network", func() {
			params := []zcash.Params{
				zcash.MainNetParams,
				zcash.TestNet3Params,
				zcash.RegressionNetParams,
			}

			for i, param := range params {
				// Generate a P2PKH address with the params
				pk := id.NewPrivKey()
				wif, err := btcutil.NewWIF((*btcec.PrivateKey)(pk), param.Params, true)
				Expect(err).NotTo(HaveOccurred())
				addrPubKeyHash, err := zcash.NewAddressPubKeyHash(btcutil.Hash160(wif.PrivKey.PubKey().SerializeUncompressed()), &param)
				Expect(err).NotTo(HaveOccurred())
				p2pkhAddr := address.Address(addrPubKeyHash.EncodeAddress())

				// Generate a P2SH address with the params
				script := make([]byte, rand.Intn(100))
				rand.Read(script)
				addrScriptHash, err := zcash.NewAddressScriptHash(script, &param)
				Expect(err).NotTo(HaveOccurred())
				p2shAddr := address.Address(addrScriptHash.EncodeAddress())

				// Try decode the address using decoders with different network params
				for j := range params {
					addrEncodeDecoder := zcash.NewAddressEncodeDecoder(&params[j])
					_, err := addrEncodeDecoder.DecodeAddress(p2pkhAddr)
					// Check the prefix in the params instead of comparing the network directly
					// because testnet and regression network has the same prefix.
					if bytes.Equal(params[i].P2PKHPrefix, params[j].P2PKHPrefix) {
						Expect(err).NotTo(HaveOccurred())
					} else {
						Expect(err).To(HaveOccurred())
					}

					_, err = addrEncodeDecoder.DecodeAddress(p2shAddr)
					if bytes.Equal(params[i].P2PKHPrefix, params[j].P2PKHPrefix) {
						Expect(err).NotTo(HaveOccurred())
					} else {
						Expect(err).To(HaveOccurred())
					}
				}
			}
		})
	})
})

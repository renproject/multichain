package zcash_test

import (
	"math/rand"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
})

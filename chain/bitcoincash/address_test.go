package bitcoincash_test

import (
	"math/rand"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/id"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/chain/bitcoincash"
)

var _ = Describe("Bitcoin Cash Address", func() {
	Context("address", func() {
		addrEncodeDecoder := bitcoincash.NewAddressEncodeDecoder(&chaincfg.RegressionNetParams)

		It("addr pub key hash", func() {
			pk := id.NewPrivKey()
			wif, err := btcutil.NewWIF((*btcec.PrivateKey)(pk), &chaincfg.RegressionNetParams, true)
			Expect(err).NotTo(HaveOccurred())
			addrPubKeyHash, err := bitcoincash.NewAddressPubKeyHash(btcutil.Hash160(wif.PrivKey.PubKey().SerializeUncompressed()), &chaincfg.RegressionNetParams)
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
			addrScriptHash, err := bitcoincash.NewAddressScriptHash(script, &chaincfg.RegressionNetParams)
			Expect(err).NotTo(HaveOccurred())
			addr := address.Address(addrScriptHash.EncodeAddress())

			decodedRawAddr, err := addrEncodeDecoder.DecodeAddress(addr)
			Expect(err).NotTo(HaveOccurred())
			encodedAddr, err := addrEncodeDecoder.EncodeAddress(decodedRawAddr)
			Expect(err).NotTo(HaveOccurred())
			Expect(encodedAddr).To(Equal(addr))
		})

		It("legacy addr", func() {
			pk := id.NewPrivKey()
			wif, err := btcutil.NewWIF((*btcec.PrivateKey)(pk), &chaincfg.RegressionNetParams, true)
			Expect(err).NotTo(HaveOccurred())
			addrPubKeyHash, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.RegressionNetParams)
			Expect(err).NotTo(HaveOccurred())
			addr := address.Address(addrPubKeyHash.EncodeAddress())

			decodedRawAddr, err := addrEncodeDecoder.DecodeAddress(addr)
			Expect(err).NotTo(HaveOccurred())
			encodedAddr, err := addrEncodeDecoder.EncodeAddress(decodedRawAddr)
			Expect(err).NotTo(HaveOccurred())
			Expect(encodedAddr).To(Equal(addr))
		})
	})
})

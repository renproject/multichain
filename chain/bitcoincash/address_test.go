package bitcoincash_test

import (
	"fmt"
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

var _ = Describe("Bitcoin Cash", func() {
	Context("Address Encode/Decode", func() {
		addrEncodeDecoders := []struct {
			network       *chaincfg.Params
			encodeDecoder bitcoincash.AddressEncodeDecoder
		}{
			{
				&chaincfg.MainNetParams,
				bitcoincash.NewAddressEncodeDecoder(&chaincfg.MainNetParams),
			},
			{
				&chaincfg.TestNet3Params,
				bitcoincash.NewAddressEncodeDecoder(&chaincfg.TestNet3Params),
			},
			{
				&chaincfg.RegressionNetParams,
				bitcoincash.NewAddressEncodeDecoder(&chaincfg.RegressionNetParams),
			},
		}

		for _, addrEncodeDecoder := range addrEncodeDecoders {
			addrEncodeDecoder := addrEncodeDecoder
			Context(fmt.Sprintf("Encode/Decode for %v network", addrEncodeDecoder.network.Name), func() {
				Specify("AddressPubKeyHash", func() {
					pk := id.NewPrivKey()
					wif, err := btcutil.NewWIF((*btcec.PrivateKey)(pk), addrEncodeDecoder.network, true)
					Expect(err).NotTo(HaveOccurred())
					addrPubKeyHash, err := bitcoincash.NewAddressPubKeyHash(btcutil.Hash160(wif.PrivKey.PubKey().SerializeUncompressed()), addrEncodeDecoder.network)
					Expect(err).NotTo(HaveOccurred())
					addr := address.Address(addrPubKeyHash.EncodeAddress())

					decodedRawAddr, err := addrEncodeDecoder.encodeDecoder.DecodeAddress(addr)
					Expect(err).NotTo(HaveOccurred())
					encodedAddr, err := addrEncodeDecoder.encodeDecoder.EncodeAddress(decodedRawAddr)
					Expect(err).NotTo(HaveOccurred())
					Expect(encodedAddr).To(Equal(addr))
				})

				Specify("AddressScriptHash", func() {
					script := make([]byte, rand.Intn(100))
					rand.Read(script)
					addrScriptHash, err := bitcoincash.NewAddressScriptHash(script, addrEncodeDecoder.network)
					Expect(err).NotTo(HaveOccurred())
					addr := address.Address(addrScriptHash.EncodeAddress())

					decodedRawAddr, err := addrEncodeDecoder.encodeDecoder.DecodeAddress(addr)
					Expect(err).NotTo(HaveOccurred())
					encodedAddr, err := addrEncodeDecoder.encodeDecoder.EncodeAddress(decodedRawAddr)
					Expect(err).NotTo(HaveOccurred())
					Expect(encodedAddr).To(Equal(addr))
				})

				Specify("AddressLegacy", func() {
					pk := id.NewPrivKey()
					wif, err := btcutil.NewWIF((*btcec.PrivateKey)(pk), addrEncodeDecoder.network, true)
					Expect(err).NotTo(HaveOccurred())
					addrPubKeyHash, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), addrEncodeDecoder.network)
					Expect(err).NotTo(HaveOccurred())
					addr := address.Address(addrPubKeyHash.EncodeAddress())

					decodedRawAddr, err := addrEncodeDecoder.encodeDecoder.DecodeAddress(addr)
					Expect(err).NotTo(HaveOccurred())
					encodedAddr, err := addrEncodeDecoder.encodeDecoder.EncodeAddress(decodedRawAddr)
					Expect(err).NotTo(HaveOccurred())
					Expect(encodedAddr).To(Equal(addr))
				})
			})
		}
	})
})

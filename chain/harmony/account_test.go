package harmony_test

import (
	"context"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/chain/harmony"
	"github.com/renproject/pack"
	"time"
)

var _ = Describe("Harmony", func() {
	Context("when broadcasting a tx", func() {
		It("should work", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := harmony.NewClient(harmony.DefaultClientOptions())
			chainId, err := c.ChainId(ctx)
			Expect(err).NotTo(HaveOccurred())

			x := "1f84c95ac16e6a50f08d44c7bde7aff8742212fda6e4321fde48bf83bef266dc"
			senderKey, err := crypto.HexToECDSA(x)
			Expect(err).NotTo(HaveOccurred())
			addrBytes, err := bech32.ConvertBits(crypto.PubkeyToAddress(senderKey.PublicKey).Bytes(), 8, 5, true)
			Expect(err).NotTo(HaveOccurred())
			senderAddr, err := bech32.Encode(harmony.Bech32AddressHRP, addrBytes)
			Expect(err).NotTo(HaveOccurred())

			toKey, _ := crypto.GenerateKey()
			toAddrBytes, err := bech32.ConvertBits(crypto.PubkeyToAddress(toKey.PublicKey).Bytes(), 8, 5, true)
			Expect(err).NotTo(HaveOccurred())
			toAddr, err := bech32.Encode(harmony.Bech32AddressHRP, toAddrBytes)
			Expect(err).NotTo(HaveOccurred())

			nonce, err := c.AccountNonce(ctx, address.Address(senderAddr))
			Expect(err).NotTo(HaveOccurred())

			gasLimit := uint64(80000000)
			gas, err := harmony.Estimator{}.EstimateGasPrice(ctx)
			Expect(err).NotTo(HaveOccurred())

			amount := pack.NewU256FromU64(pack.NewU64(100000000))

			txBuilder := harmony.NewTxBuilder(chainId)
			tx, err := txBuilder.BuildTx(ctx, address.Address(senderAddr), address.Address(toAddr), amount, nonce, pack.NewU256FromU64(pack.NewU64(gasLimit)), gas, gas, pack.Bytes(nil))
			Expect(err).NotTo(HaveOccurred())

			sigHash, err := tx.Sighashes()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(sigHash)).To(Equal(1))

			sig, err := crypto.Sign(sigHash[0][:], senderKey)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(sig)).To(Equal(65))

			var signature [65]byte
			copy(signature[:], sig)
			err = tx.Sign([]pack.Bytes65{pack.NewBytes65(signature)}, pack.Bytes(nil))
			Expect(err).ToNot(HaveOccurred())

			err = c.SubmitTx(ctx, tx)
			Expect(err).ToNot(HaveOccurred())

			time.Sleep(time.Second)
			for {
				txResp, _, err := c.Tx(ctx, tx.Hash())
				if err == nil && txResp != nil {
					break
				}
				// wait and retry querying for the transaction
				time.Sleep(5 * time.Second)
			}

			updatedBalance, err := c.AccountBalance(ctx, address.Address(toAddr))
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedBalance).To(Equal(amount))

		})
	})
})

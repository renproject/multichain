package waves_test

import (
	"context"
	"fmt"
	"time"

	"github.com/mr-tron/base58"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/renproject/multichain/api/address"

	"github.com/renproject/multichain/chain/waves"
	"github.com/renproject/pack"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

var _ = Describe("Waves", func() {
	Context("when submitting transactions", func() {
		Context("when sending waves to address", func() {
			It("should work", func() {
				defer func() {
					fmt.Println("Waves test ended.")
				}()
				client := waves.NewClient("127.0.0.1:6870", 'I')
				// Hardcoded seed with balance.
				seed, _ := base58.Decode("CApJrVsZ6AY5zbunL2nqgrb7MkJF9rPiFz63RtaRPyna")
				from, _ := proto.NewKeyPair(seed)

				// Random seed.
				to, _ := proto.NewKeyPair([]byte("test4"))
				toAddr, _ := to.Addr('I')

				builder := waves.NewTxBuilder('I')
				value := pack.NewU256FromU64(pack.NewU64(100))

				tx, err := builder.BuildTx(from.Public.String(), address.Address(toAddr.String()), value, pack.U256{}, nil)
				Expect(err).ToNot(HaveOccurred())

				err = tx.Sign(from.Secret.Bytes())
				Expect(err).ToNot(HaveOccurred())

				err = client.SubmitTx(context.Background(), tx)
				Expect(err).ToNot(HaveOccurred())

				timeout := time.After(5 * time.Minute)
				for {
					select {
					case <-timeout:
						panic("not found")
					default:
					}
					id := tx.Hash()
					rs, confirmations, _ := client.Tx(context.Background(), id)
					if confirmations == 0 {
						<-time.After(10 * time.Second)
						continue
					}
					fmt.Println("confirmations: ", confirmations)
					Expect(rs.Hash()).To(Equal(id))
					if confirmations >= 3 {
						break
					}
					<-time.After(10 * time.Second)
				}
			})
		})
	})
})

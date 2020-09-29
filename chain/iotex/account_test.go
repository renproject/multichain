package iotex_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/chain/iotex"
	"github.com/renproject/pack"
	"google.golang.org/grpc"
)

var _ = Describe("IoTeX", func() {
	Context("when decoding address", func() {
		Context("when decoding IoTeX address", func() {
			It("should work", func() {
				decoder := iotex.NewAddressDecoder()
				addrStr := "io17ch0jth3dxqa7w9vu05yu86mqh0n6502d92lmp"
				_, err := decoder.DecodeAddress(address.Address(pack.NewString(addrStr)))
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
var _ = Describe("IoTeX", func() {
	Context("when submitting transactions", func() {
		Context("when sending IOTX", func() {
			It("should work", func() {
				pkEnv := os.Getenv("pk")
				if pkEnv == "" {
					panic("pk is undefined")
				}
				gopts := []grpc.DialOption{}
				gopts = append(gopts, grpc.WithInsecure())
				endpoint := "api.testnet.iotex.one:80"
				conn, err := grpc.Dial(endpoint, gopts...)
				c := iotexapi.NewAPIServiceClient(conn)
				request := &iotexapi.GetAccountRequest{Address: "io1vdtfpzkwpyngzvx7u2mauepnzja7kd5rryp0sg"}
				res, err := c.GetAccount(context.Background(), request)
				Expect(err).ToNot(HaveOccurred())

				opts := iotex.ClientOptions{
					Endpoint: endpoint,
					Secure:   false,
				}
				client := iotex.NewClient(opts)
				builder := iotex.TxBuilder{}
				gasPrice, _ := new(big.Int).SetString("1000000000000", 10)
				tx, err := builder.BuildTx("io1vdtfpzkwpyngzvx7u2mauepnzja7kd5rryp0sg", "io1vdtfpzkwpyngzvx7u2mauepnzja7kd5rryp0sg", pack.NewU256FromU64(pack.NewU64(1)), pack.NewU256FromU64(pack.NewU64(res.AccountMeta.PendingNonce)), pack.NewU256FromInt(gasPrice), pack.NewU256FromU64(pack.NewU64(1000000)), nil)
				Expect(err).ToNot(HaveOccurred())
				sh, err := tx.Sighashes()
				Expect(err).ToNot(HaveOccurred())
				sk, err := crypto.HexStringToPrivateKey(pkEnv)
				Expect(err).ToNot(HaveOccurred())
				sig, err := sk.Sign(sh[0][:])
				Expect(err).ToNot(HaveOccurred())
				var sig65 [65]byte
				copy(sig65[:], sig[:])
				tx.Sign([]pack.Bytes65{pack.NewBytes65(sig65)}, sk.PublicKey().Bytes())
				sigHash, err := tx.Sighashes()
				Expect(err).ToNot(HaveOccurred())

				sHash := hex.EncodeToString(sigHash[0][:])
				Expect(err).ToNot(HaveOccurred())
				fmt.Println("sig hash:", sHash)
				fmt.Println("public key:", hex.EncodeToString(sk.PublicKey().Bytes()))
				fmt.Println("transaction hash:", hex.EncodeToString(tx.Hash()))

				err = client.SubmitTx(context.Background(), tx)
				Expect(err).ToNot(HaveOccurred())

				// We wait for 10 s before beginning to check transaction.
				time.Sleep(10 * time.Second)
				returnedTx, n, err := client.Tx(context.Background(), tx.Hash())
				Expect(err).ToNot(HaveOccurred())
				Expect(n).To(Equal(pack.NewU64(1)))
				Expect(returnedTx.Nonce()).To(Equal(pack.NewU256FromU64(pack.NewU64(res.AccountMeta.PendingNonce))))
			})
		})
	})
})

package harmony_test

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/harmony-one/harmony/rpc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/chain/harmony"
)

var _ = Describe("Harmony", func() {
	Context("when calling a contract", func() {
		It("should work", func() {
			contractAddr := address.Address("one155jp2y76nazx8uw5sa94fr0m4s5aj8e5xm6fu3")
			rawAddr, err := harmony.NewEncoderDecoder().DecodeAddress(contractAddr)
			bech32Addr := common.BytesToAddress(rawAddr)
			callData := rpc.CallArgs{
				To: &bech32Addr,
			}
			params := harmony.Params{
				CallArgs: callData,
				Block:    37000,
			}
			marshalledData, err := json.Marshal(params)
			Expect(err).NotTo(HaveOccurred())

			c := harmony.NewClient(harmony.DefaultClientOptions())
			_, err = c.CallContract(context.TODO(), contractAddr, marshalledData)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
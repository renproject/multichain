package icon_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/multichain"
	"github.com/renproject/multichain/chain/icon"
	"github.com/renproject/pack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ = Describe("Icon Blockchain", func() {
	ctx := context.Background()

	// Initialise the logger.
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	Context("Gas Estimation", func() {
		accountChainTable := []struct {
			rpcURL pack.String
			chain  multichain.Chain
		}{
			{
				"http://127.0.0.1:9000/api/v3",
				multichain.Icon,
			},
		}

		for _, accountChain := range accountChainTable {
			accountChain := accountChain
			Context(fmt.Sprintf("%v", accountChain.chain), func() {
				Specify("get step limit", func() {
					iconClient := icon.NewClient(accountChain.rpcURL)

					es := icon.Estimator{
						Client: *iconClient,
					}

					resultEs, err := es.EstimateGasPrice(ctx)
					Expect(err).NotTo(HaveOccurred())
					Expect(resultEs.String()).To(Equal("2500000000"))
				})
			})
		}
	})
})

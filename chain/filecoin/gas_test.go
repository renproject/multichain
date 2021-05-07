package filecoin_test

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/renproject/multichain/chain/filecoin"
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gas", func() {
	Context("when estimating gas parameters", func() {
		It("should work", func() {
			// create context for the test
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// instantiate the client
			client, err := filecoin.NewClient(
				filecoin.DefaultClientOptions().
					WithAuthToken(fetchAuthToken()),
			)
			Expect(err).ToNot(HaveOccurred())

			// instantiate the gas estimator
			gasEstimator := filecoin.NewGasEstimator(client, 2000000)

			// estimate gas price
			_, _, err = gasEstimator.EstimateGas(ctx)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

func fetchAuthToken() pack.String {
	// fetch the auth token from filecoin's running docker container
	cmd := exec.Command("docker", "exec", "infra_filecoin_1", "/bin/bash", "-c", "/app/lotus auth api-info --perm admin")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		panic(fmt.Sprintf("could not run command: %v", err))
	}
	tokenWithSuffix := strings.TrimPrefix(out.String(), "FULLNODE_API_INFO=")
	authToken := strings.Split(tokenWithSuffix, ":/")
	return pack.NewString(fmt.Sprintf("Bearer %s", authToken[0]))
}

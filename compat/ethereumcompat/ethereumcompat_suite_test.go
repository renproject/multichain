package ethereumcompat_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEthereumCompat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ethereum Compat Suite")
}

package litecoin_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLitecoin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Litecoin Suite")
}

package multichain_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMultichain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Multichain Suite")
}

package solana_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSolana(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Solana Suite")
}

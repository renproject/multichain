package utxochain_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBitcoin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bitcoin Suite")
}

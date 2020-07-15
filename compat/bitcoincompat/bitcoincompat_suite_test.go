package bitcoincompat_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBitcoinCompat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bitcoin Compat Suite")
}

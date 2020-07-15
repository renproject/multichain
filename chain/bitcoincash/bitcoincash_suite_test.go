package bitcoincash_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBitcoinCash(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bitcoin Cash Suite")
}

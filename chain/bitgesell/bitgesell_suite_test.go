package bitgesell_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBitgesell(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bitgesell Suite")
}

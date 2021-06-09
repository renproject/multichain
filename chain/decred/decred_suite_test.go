package decred_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDecred(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Decred Suite")
}

package lbry_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLBRY(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LBRY Suite")
}

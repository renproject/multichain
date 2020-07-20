package zcoin_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestZcoin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Zcoin Suite")
}

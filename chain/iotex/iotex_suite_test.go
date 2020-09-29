package iotex_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEthereum(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IoTeX Suite")
}

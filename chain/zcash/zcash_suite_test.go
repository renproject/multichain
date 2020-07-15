package zcash_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestZcash(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Zcash Suite")
}

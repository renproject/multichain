package nix_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNIX(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NIX Suite")
}

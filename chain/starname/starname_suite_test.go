package starname_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStarname(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Starname (IOV) Suite")
}

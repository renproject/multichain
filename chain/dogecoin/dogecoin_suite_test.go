package dogecoin_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDogecoin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dogecoin Suite")
}

package filecoin_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFilecoin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filecoin Suite")
}

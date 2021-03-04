package qtum_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQtum(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Qtum Suite")
}

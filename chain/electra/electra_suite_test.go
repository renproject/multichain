package electra_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestElectra(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Electra Suite")
} 
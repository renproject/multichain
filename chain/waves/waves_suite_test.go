package waves_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWaves(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Waves Suite")
}

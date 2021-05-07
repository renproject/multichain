package harmony

import "testing"
import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHarmony(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Harmony Suite")
}
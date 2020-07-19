package digibyte_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDigiByte(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DigiByte Suite")
}

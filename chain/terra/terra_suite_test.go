package terra_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTerra(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Terra Suite")
}

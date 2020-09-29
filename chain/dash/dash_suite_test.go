package dash_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDash(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dash Suite")
}

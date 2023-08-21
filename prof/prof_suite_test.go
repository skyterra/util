package prof

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestProf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Prof Suite")
}

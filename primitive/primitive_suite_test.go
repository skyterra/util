package primitive_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPrimitive(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Primitive Suite")
}

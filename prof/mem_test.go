package prof

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Mem", func() {
	Context("Memory Usage", func() {
		It("should be succeed", func() {
			PrintMemUsage()
			b := make([]byte, 100<<20)
			b[1<<10] = 255
			PrintMemUsage("after alloc")
		})
	})
})

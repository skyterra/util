package pprof

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Mem", func() {
	Context("Memory Usage", func() {
		It("should be succeed", func() {
			PrintMemUsage()
		})
	})
})

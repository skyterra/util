package pprof

import (
	"fmt"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Pprof", func() {
	Context("Pprof", func() {
		It("CPU", func() {
			f := ProfCPUStart()

			p, q := uint64(0), uint64(1)
			for i := 0; i < 1000000; i++ {
				p, q = q, p+q
			}

			fmt.Println(p)
			ProfCPUStop(f)
		})

		It("Memory", func() {
			s := make([]byte, 0, 2*GB)
			fmt.Println(len(s))

			ProfMemory()
			PrintMemUsage()
		})
	})
})

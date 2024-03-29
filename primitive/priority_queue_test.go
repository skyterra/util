package primitive_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math/rand"
	"sort"
	"sync"
	"testing"
	"util/primitive"
)

type Number int

func (n Number) Priority() int64 {
	return int64(n)
}

var g = primitive.NewPriorityQueue(10)

var _ = Describe("PriorityQueue", func() {
	Context("SafePush & SafePop", func() {
		It("normal", func() {
			var n Number

			pq := primitive.NewPriorityQueue(10)
			size := 100

			var expect []Number
			for i := 0; i < size; i++ {
				n = Number(rand.Int())
				expect = append(expect, n)
				pq.SafePush(n)
			}

			sort.Slice(expect, func(i, j int) bool {
				return expect[i] < expect[j]
			})

			var real []Number
			for i := 0; i < size; i++ {
				real = append(real, pq.SafePop().(Number))
			}

			Expect(len(expect) == len(real)).Should(BeTrue())

			for i := 0; i < len(expect); i++ {
				Expect(expect[i] == real[i]).Should(BeTrue())
			}
		})

		It("pop empty", func() {
			pq := primitive.NewPriorityQueue(0)
			Expect(pq.SafePop() == nil).Should(BeTrue())
		})

		It("pop all", func() {
			var n Number

			pq := primitive.NewPriorityQueue(10)
			size := 100

			var expect []Number
			for i := 0; i < size; i++ {
				n = Number(rand.Int())
				expect = append(expect, n)
				pq.SafePush(n)
			}

			sort.Slice(expect, func(i, j int) bool {
				return expect[i] < expect[j]
			})

			real := pq.SafePopAll()
			Expect(len(expect) == len(real)).Should(BeTrue())

			for i := 0; i < len(expect); i++ {
				Expect(expect[i] == real[i]).Should(BeTrue())
			}
		})
	})

	Context("concurrent", func() {
		It("push", func() {
			finish := sync.WaitGroup{}
			finish.Add(2)

			for i := 0; i < 10000; i++ {
				g.SafePush(Number(rand.Int()))
			}

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				wg.Wait()
				for i := 0; i < 10000; i++ {
					g.SafePush(Number(rand.Int()))
				}

				finish.Done()
			}()

			go func() {
				wg.Wait()
				for g.SafePop() != nil {

				}

				finish.Done()
			}()

			wg.Done()
			finish.Wait()
			Expect(g.Len() == 0).Should(BeTrue())
		})
	})
})

func BenchmarkPriorityQueuePush(b *testing.B) {
	p := primitive.NewPriorityQueue(1024)
	for i := 0; i < b.N; i++ {
		p.SafePush(Number(i))
	}
}

func BenchmarkPriorityQueuePop(b *testing.B) {
	p := primitive.NewPriorityQueue(1024)
	for i := 0; i < b.N; i++ {
		p.SafePush(Number(i))
	}

	for p.SafePop() != nil {

	}
}

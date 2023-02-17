package primitive_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
	. "util/primitive"
)

var _ = Describe("BlockingQueue", func() {
	Context("push", func() {
		It("block", func() {
			q := NewBlockingQueue(2, 6)
			exit := false
			go func() {
				for i := 0; i < 10; i++ {
					q.Push(Number(i))
				}
				exit = true
			}()

			time.Sleep(10 * time.Millisecond)
			Expect(q.Len() == 6).Should(BeTrue())

			q.EndWait()
			time.Sleep(10 * time.Millisecond)
			Expect(exit).Should(BeTrue())
		})

		It("try push", func() {
			q := NewBlockingQueue(2, 6)
			for i := 0; i < 5; i++ {
				q.Push(Number(i))
			}

			Expect(q.TryPush(Number(100))).Should(Succeed())
			Expect(q.TryPush(Number(101))).ShouldNot(Succeed())

			Expect(q.Len() == 6).Should(BeTrue())
		})

		It("pop", func() {
			q := NewBlockingQueue(2, 6)
			exit := false

			go func() {
				e := q.Pop().(Number)
				Expect(e == 10).Should(BeTrue())
				exit = true
			}()

			time.Sleep(10 * time.Millisecond)

			q.Push(Number(10))
			time.Sleep(10 * time.Millisecond)

			Expect(exit).Should(BeTrue())
		})

		It("pop all", func() {
			q := NewBlockingQueue(2, 6)
			for i := 0; i < 5; i++ {
				q.Push(Number(i))
			}

			Expect(len(q.PopAll()) == 5).Should(BeTrue())

			exit := false
			go func() {
				Expect(len(q.PopAll()) == 1).Should(BeTrue())
				exit = true
			}()

			time.Sleep(10 * time.Millisecond)

			q.Push(Number(100))
			time.Sleep(10 * time.Millisecond)

			Expect(exit).Should(BeTrue())
		})

		It("try pop", func() {
			q := NewBlockingQueue(2, 6)
			Expect(q.TryPop() == nil).Should(BeTrue())

			q.Push(Number(10))
			Expect(q.TryPop().(Number) == 10).Should(BeTrue())
		})

		It("try pop all", func() {
			q := NewBlockingQueue(2, 6)
			Expect(len(q.TryPopAll()) == 0).Should(BeTrue())

			for i := 0; i < 5; i++ {
				q.Push(Number(i))
			}

			Expect(len(q.TryPopAll()) == 5).Should(BeTrue())
		})

		It("initSize & maxSize", func() {
			q := NewBlockingQueue(6, 2)
			for i := 0; i < 6; i++ {
				q.Push(Number(i))
			}
			Expect(q.Len() == 6).Should(BeTrue())
		})

		It("no limit", func() {
			q := NewBlockingQueue(1024, 0)
			for i := 0; i < 1024; i++ {
				q.Push(Number(i))
			}

			Expect(q.Len() == 1024).Should(BeTrue())
		})
	})
})

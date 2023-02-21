package schedule_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	"math/rand"
	. "util/schedule"
)

type Task struct {
	id      int
	handler func()
}

func (t *Task) Do() {
	fmt.Println("run task", t.id)
}

func (t *Task) Priority() int64 {
	return int64(t.id)
}

func (t *Task) OnPanic(err error) {
	fmt.Println(err.Error())
}

var _ = Describe("ThreadPool", func() {
	Context("push", func() {
		It("should be succeed", func() {
			p := NewThreadPool(2, 6, 10)
			p.Start()

			for i := 0; i < 100; i++ {
				p.Push(&Task{id: rand.Int()})
			}
		})
	})
})

package schedule

import (
	"sync"
	"util/primitive"
)

type ITheadTask interface {
	Do() error
}

type ThreadPool struct {
	q             *primitive.PriorityQueue
	wg            sync.WaitGroup
	maxGoroutines int
}

func (p *ThreadPool) Push(t ITheadTask) {

}

func (p *ThreadPool) Start() {
	for i := 0; i < p.maxGoroutines; i++ {
		go func() {

		}()
	}
}

func NewThreadPool(maxGoroutines int) *ThreadPool {
	return &ThreadPool{
		q:             primitive.NewPriorityQueue(1024),
		maxGoroutines: maxGoroutines,
	}
}

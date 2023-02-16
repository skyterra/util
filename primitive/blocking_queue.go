package primitive

import (
	"errors"
	"sync"
)

/*
 * 构建优先级队列阻塞队列
 */

type BlockingQueue struct {
	q *PriorityQueue

	endWait bool
	maxSize int
	cond    *sync.Cond
}

func (bq *BlockingQueue) Len() int {
	return bq.q.Len()
}

func (bq *BlockingQueue) TryPush(element IPriorityElement) error {
	if bq.endWait {
		return errors.New("blocking queue has been aborted")
	}

	bq.cond.L.Lock()
	defer bq.cond.L.Unlock()

	if bq.maxSize > 0 && bq.q.Len() == bq.maxSize {
		return errors.New("blocking queue is full")
	}

	bq.q.Push(element)

	return nil
}

func (bq *BlockingQueue) TryPop() IPriorityElement {
	if bq.endWait {
		return nil
	}

	bq.cond.L.Lock()
	defer bq.cond.L.Unlock()

	return bq.q.Pop()
}

func (bq *BlockingQueue) TryPopAll() []IPriorityElement {
	if bq.endWait {
		return nil
	}

	bq.cond.L.Lock()
	defer bq.cond.L.Unlock()

	return bq.q.PopAll()
}

func (bq *BlockingQueue) Push(element IPriorityElement) {
	if bq.endWait {
		return
	}

	bq.cond.L.Lock()
	if bq.maxSize > 0 && bq.q.Len() == bq.maxSize {
		bq.cond.Wait()
	}

	bq.q.Push(element)
	bq.cond.L.Unlock()
}

func (bq *BlockingQueue) Pop() IPriorityElement {
	bq.cond.L.Lock()
	if !bq.endWait && bq.q.Len() == 0 {
		bq.cond.Wait()
	}

	r := bq.q.Pop()
	bq.cond.L.Unlock()

	return r
}

func (bq *BlockingQueue) PopAll() []IPriorityElement {
	bq.cond.L.Lock()
	if !bq.endWait && bq.q.Len() == 0 {
		bq.cond.Wait()
	}

	r := bq.q.PopAll()
	bq.cond.L.Unlock()

	return r
}

func (bq *BlockingQueue) EndWait() {
	bq.endWait = true
	bq.cond.Broadcast()
}

// NewBlockingQueue 创建阻塞队列，size为初始队列大小，maxSize为队列大小上限，如果maxSize为0，表示
// 阻塞队列无上限
func NewBlockingQueue(size int, maxSize int) *BlockingQueue {
	return &BlockingQueue{
		q:    NewPriorityQueue(size),
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

package primitive

import (
	"errors"
	"sync"
)

/*
 * 构建优先级队列阻塞队列
 */

// 定义阻塞队列
type BlockingQueue struct {
	q *PriorityQueue

	maxSize int
	cond    *sync.Cond
}

// Len 返回队列长度
func (bq *BlockingQueue) Len() int {
	return bq.q.Len()
}

// TryPush 尝试向队列添加元素（非阻塞）
func (bq *BlockingQueue) TryPush(element IPriorityElement) error {
	bq.cond.L.Lock()
	defer bq.cond.L.Unlock()

	if bq.maxSize > 0 && bq.q.Len() == bq.maxSize {
		return errors.New("blocking queue is full")
	}

	bq.q.Push(element)
	bq.cond.Signal()

	return nil
}

// TryPop 尝试弹出元素（非阻塞）
func (bq *BlockingQueue) TryPop() interface{} {
	bq.cond.L.Lock()
	defer bq.cond.L.Unlock()

	return bq.q.Pop()
}

// TryPopAll 尝试弹出所有元素（非阻塞）
func (bq *BlockingQueue) TryPopAll() []interface{} {
	bq.cond.L.Lock()
	defer bq.cond.L.Unlock()

	return bq.q.PopAll()
}

// Push 向队列通添加元素，如果设置了maxSize且队列已满，则进入等待
func (bq *BlockingQueue) Push(element IPriorityElement) {
	bq.cond.L.Lock()
	if bq.maxSize > 0 && bq.q.Len() == bq.maxSize {
		bq.cond.Wait()
	}

	bq.q.Push(element)
	bq.cond.L.Unlock()

	bq.cond.Signal()
}

// Push 从队列中弹出元素，如果队列为空，进入等待
func (bq *BlockingQueue) Pop() interface{} {
	bq.cond.L.Lock()
	if bq.q.Len() == 0 {
		bq.cond.Wait()
	}

	r := bq.q.Pop()
	bq.cond.L.Unlock()

	return r
}

// PushAll 从队列中弹出所有元素，如果队列为空，进入等待
func (bq *BlockingQueue) PopAll() []interface{} {
	bq.cond.L.Lock()
	if bq.q.Len() == 0 {
		bq.cond.Wait()
	}

	r := bq.q.PopAll()
	bq.cond.L.Unlock()

	return r
}

// EndWait 结束所有等待
func (bq *BlockingQueue) EndWait() {
	bq.cond.Broadcast()
}

// NewBlockingQueue 创建阻塞队列，initSize为初始队列大小，maxSize为队列大小上限，如果maxSize为0，表示
// 阻塞队列无上限
func NewBlockingQueue(initSize int, maxSize int) *BlockingQueue {
	if maxSize > 0 && maxSize < initSize {
		maxSize = initSize
	}

	return &BlockingQueue{
		q:       NewPriorityQueue(initSize),
		cond:    sync.NewCond(&sync.Mutex{}),
		maxSize: maxSize,
	}
}

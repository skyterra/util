package primitive

import (
	"container/heap"
	"sync"
)

/*
 * 使用小顶堆构建优先级队列，优先级值越小代表越先执行
 */

// 定义优先级元素接口
type IPriorityElement interface {
	GetPriority() int64
}

type queue []IPriorityElement

func (q queue) Len() int {
	return len(q)
}

func (q *queue) Push(x interface{}) {
	*q = append(*q, x.(IPriorityElement))
}

func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)

	x := old[n-1]
	*q = old[0 : n-1]
	return x
}

func (q queue) Less(i, j int) bool {
	return q[i].GetPriority() < q[j].GetPriority()
}

func (q queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

// 定义线程安全的优先级队列
type PriorityQueue struct {
	locker sync.Mutex
	q      queue
}

// Push 向队列中添加元素
func (pq *PriorityQueue) Push(element IPriorityElement) {
	pq.locker.Lock()
	defer pq.locker.Unlock()

	heap.Push(&pq.q, element)
}

// Pop 按照优先级从小到达顺序弹出元素
func (pq *PriorityQueue) Pop() IPriorityElement {
	pq.locker.Lock()
	defer pq.locker.Unlock()

	if pq.Len() > 0 {
		return heap.Pop(&pq.q).(IPriorityElement)
	}

	return nil
}

// Len 获取优先级队列长度
func (pq *PriorityQueue) Len() int {
	return len(pq.q)
}

// NewPriorityQueue 构建优先级队列，size为队列容量
func NewPriorityQueue(size int) *PriorityQueue {
	return &PriorityQueue{
		q: make([]IPriorityElement, 0, size),
	}
}

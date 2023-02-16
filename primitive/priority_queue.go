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

// 定义协程安全的优先级队列
type PriorityQueue struct {
	locker sync.Mutex
	q      queue
}

// Push 向队列中添加元素
func (pq *PriorityQueue) Push(element IPriorityElement) {
	heap.Push(&pq.q, element)
}

// Pop 按照优先级从小到达顺序弹出元素
func (pq *PriorityQueue) Pop() IPriorityElement {
	if len(pq.q) == 0 {
		return nil
	}

	return heap.Pop(&pq.q).(IPriorityElement)
}

// PopAll 按照优先级从小到达顺序弹出所有元素
func (pq *PriorityQueue) PopAll() []IPriorityElement {
	elements := make([]IPriorityElement, 0, len(pq.q))

	for len(pq.q) > 0 {
		elements = append(elements, heap.Pop(&pq.q).(IPriorityElement))
	}

	return elements
}

// SafePush 协程安全，向队列中添加元素
func (pq *PriorityQueue) SafePush(element IPriorityElement) {
	pq.locker.Lock()
	defer pq.locker.Unlock()

	pq.Push(element)
}

// SafePop 协程安全，按照优先级从小到达顺序弹出元素
func (pq *PriorityQueue) SafePop() IPriorityElement {
	pq.locker.Lock()
	defer pq.locker.Unlock()

	return pq.Pop()
}

// SafePopAll 协程安全，按照优先级从小到达顺序弹出所有元素
func (pq *PriorityQueue) SafePopAll() []IPriorityElement {
	pq.locker.Lock()
	defer pq.locker.Unlock()

	return pq.PopAll()
}

// Len 获取优先级队列长度
func (pq *PriorityQueue) Len() int {
	return len(pq.q)
}

// NewPriorityQueue 构建优先级队列，size为队列容量
func NewPriorityQueue(cap int) *PriorityQueue {
	return &PriorityQueue{
		q: make([]IPriorityElement, 0, cap),
	}
}

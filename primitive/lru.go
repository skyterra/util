package primitive

import (
	"sync"
	"time"
)

/*
 * LRU Cache 使用双向链表管理缓存对象，最近被访问的对象会放到链表尾部，即将被淘汰的对象放到链表头部；
 * 可以通过 ttl 设置缓存对象对象存活时间，如果对象存活时间超过了 ttl，Get接口会返回nil
 */

type lruNode struct {
	timestamp int64
	key       interface{}
	value     interface{}
	pre       *lruNode
	next      *lruNode
}

// lru 缓存
type lru struct {
	cap   int
	ttl   int64
	m     map[interface{}]*lruNode
	head  *lruNode
	tail  *lruNode
	mutex sync.Mutex
}

// NewLRU 创建缓存，cap执行缓存容量，ttl执行缓存对象存活时间，-1关闭ttl
func NewLRU(cap int, ttlMS int64) *lru {
	return &lru{
		cap: cap,
		ttl: ttlMS,
		m:   make(map[interface{}]*lruNode),
	}
}

// GetLength 获取lru队列长度
func (lru *lru) GetLength() int {
	return len(lru.m)
}

// Get 获取对象
func (lru *lru) Get(key interface{}) (interface{}, bool) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	v, exist := lru.m[key]
	if !exist || (lru.ttl > 0 && time.Now().UnixNano()/1e6-v.timestamp > lru.ttl) {
		return nil, false
	}

	lru.adjust(key)
	return v.value, true
}

// Put 添加对象
func (lru *lru) Put(key interface{}, value interface{}) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	// 如果已经存在，则直接更新
	if _, exist := lru.m[key]; exist {
		lru.m[key].value = value
		lru.m[key].timestamp = time.Now().UnixNano() / 1e6

		lru.adjust(key)
		return
	}

	node := &lruNode{key: key, value: value, timestamp: time.Now().UnixNano() / 1e6}
	lru.m[key] = node

	// lru中第一个元素
	if lru.head == nil {
		lru.head = node
		lru.tail = node
		return
	}

	// 添加到链表尾部
	lru.tail.next = node
	node.pre = lru.tail
	lru.tail = lru.tail.next

	// 超出容量，删除头部元素
	if len(lru.m) > lru.cap {
		p := lru.head
		lru.head = lru.head.next
		lru.head.pre = nil

		p.next = nil
		delete(lru.m, p.key)
	}

	return
}

func (lru *lru) adjust(key interface{}) {
	node := lru.m[key]

	// 已经在尾部
	if node.next == nil {
		return
	}

	// 头部
	if node.pre == nil {
		lru.head = node.next
		lru.head.pre = nil
	} else {
		node.pre.next = node.next
		node.next.pre = node.pre
	}

	lru.tail.next = node
	node.pre = lru.tail

	lru.tail = lru.tail.next
	lru.tail.next = nil
}

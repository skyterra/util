package primitive

type Node struct {
	Key   int
	Value int

	Next *Node
	Pre  *Node
}

type LruCache struct {
	cap int
	m   map[int]*Node

	head *Node
	tail *Node
}

func (l *LruCache) Get(key int) int {
	v, exist := l.m[key]
	if !exist {
		return -1
	}

	l.adjust(key)
	return v.Value
}

func (l *LruCache) Put(key int, value int) {
	if _, exist := l.m[key]; exist {
		l.m[key].Value = value
		l.adjust(key)
		return
	}

	node := &Node{Key: key, Value: value}
	l.m[key] = node

	if l.head == nil {
		l.head = node
		l.tail = node
		return
	}

	l.tail.Next = node
	node.Pre = l.tail
	l.tail = l.tail.Next

	if len(l.m) > l.cap {
		p := l.head
		l.head = l.head.Next
		l.head.Pre = nil

		p.Next = nil
		delete(l.m, p.Key)
	}

	return
}

func (l *LruCache) adjust(key int) {
	node := l.m[key]

	// 已经在尾部
	if node.Next == nil {
		return
	}

	// 头部
	if node.Pre == nil {
		l.head = node.Next
		l.head.Pre = nil
	} else {
		node.Pre.Next = node.Next
		node.Next.Pre = node.Pre
	}

	l.tail.Next = node
	node.Pre = l.tail

	l.tail = l.tail.Next
	l.tail.Next = nil
}

func NewLruCache(cap int) LruCache {
	return LruCache{
		cap: cap,
		m:   make(map[int]*Node),
	}
}

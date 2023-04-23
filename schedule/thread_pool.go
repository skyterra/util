package schedule

import (
	"fmt"
	"sync"
	"util/primitive"
)

// 定义协程任务接口
type ITheadTask interface {
	Do()
	OnPanic(err error)
	Priority() int64
}

// 定义协程池
type ThreadPool struct {
	exit          bool
	maxGoroutines int
	queue         *primitive.BlockingQueue
	wg            sync.WaitGroup
}

// Push 向协程池任务队列中添加任务
func (p *ThreadPool) Push(t ITheadTask) {
	if !p.exit {
		p.queue.Push(t)
	}
}

// Start 启动协程池
func (p *ThreadPool) Start() {
	for i := 0; i < p.maxGoroutines; i++ {
		go func() {
			var t ITheadTask

			defer func() {
				if r := recover(); r != nil {
					safeOnPanic(t, fmt.Errorf("%v", r))
				}
			}()

			for !p.exit {
				t, _ = p.queue.Pop().(ITheadTask)
				if t != nil {
					t.Do()
				}
			}
		}()
	}
}

// Exit 退出协程池
func (p *ThreadPool) Exit() {
	p.exit = true
	p.queue.EndWait()
}

func safeOnPanic(t ITheadTask, err error) {
	defer func() {
		recover()
	}()

	if t != nil {
		t.OnPanic(err)
	}
}

// NewThreadPool 创建协程池，maxGoroutines指定工作协程数量，initSize指定任务队列初始大小，
// maxSize指定任务队列上限，达到上限后，Push会等待；maxSize为0表示无上限，即Push操作不会等待
func NewThreadPool(maxGoroutines int, initSize, maxSize int) *ThreadPool {
	return &ThreadPool{
		exit:          false,
		maxGoroutines: maxGoroutines,
		queue:         primitive.NewBlockingQueue(initSize, maxSize),
	}
}

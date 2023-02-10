package schedule

import (
	"container/heap"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 一个按照指定时间执行任务的调度器，用户可以添加ITimingTask类型的任务，并设置执行时间；
// 当时间达到时，调度器负责执行任务。
//
// 调度器的任务队列采用了小顶堆设计，保证了任务调度按照设定时间的先后顺序执行而不是加入调
// 度器的先后顺序

// 调度器接收的任务接口
type ITimingTask interface {
	RunAt() int64
	Run(schedule *TimingSchedule)
	OnError(err error)
}

// 使用小顶堆定义任务队列
type taskQueue []ITimingTask

func (q taskQueue) Len() int {
	return len(q)
}

func (q taskQueue) Less(i, j int) bool {
	return q[i].RunAt() < q[j].RunAt()
}

func (q taskQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *taskQueue) Push(x interface{}) {
	*q = append(*q, x.(ITimingTask))
}

func (q *taskQueue) Pop() interface{} {
	old := *q
	n := len(old)

	x := old[n-1]
	*q = old[0 : n-1]
	return x
}

// 调度器对象
type TimingSchedule struct {
	shutdown    bool
	workerCount int
	intervalS   int

	tasks taskQueue

	mutex sync.Mutex
	wg    sync.WaitGroup
}

// Push 获取调度器中未执行的任务数
func (s *TimingSchedule) Len() int {
	return len(s.tasks)
}

// Push 向调度器中添加任务
func (s *TimingSchedule) Push(t ITimingTask) {
	if s.shutdown {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	heap.Push(&s.tasks, t)
}

func (s *TimingSchedule) pop() ITimingTask {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now().Unix()
	if s.tasks.Len() > 0 && now > s.tasks[0].RunAt() {
		return heap.Pop(&s.tasks).(ITimingTask)
	}

	return nil
}

// 启动调度器
func (s *TimingSchedule) Start() {
	s.wg.Add(s.workerCount)

	for i := 0; i < s.workerCount; i++ {
		go func() {
			var curTask ITimingTask
			for !s.shutdown {
				if curTask = s.pop(); curTask != nil {
					// 为了防止调用task.OnError()发生panic，此处做了异常保护
					onError := func(r interface{}) {
						defer func() {
							recover()
						}()

						curTask.OnError(fmt.Errorf("%v", r))
					}

					f := func() {
						defer func() {
							if r := recover(); r != nil {
								onError(r)
							}
						}()

						curTask.Run(s)
					}

					f()
				} else {
					time.Sleep(time.Duration(s.intervalS) * time.Second)
				}
			}

			s.wg.Done()
		}()
	}
}

// Shutdown 关闭调度器
func (s *TimingSchedule) Shutdown() {
	s.shutdown = true
	s.wg.Wait()

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.tasks = nil
}

// NewTimingSchedule 创建基于具体时间点的调度器
func NewTimingSchedule(workerCount int, intervalS int, tasks ...ITimingTask) *TimingSchedule {
	s := &TimingSchedule{
		shutdown:    false,
		workerCount: workerCount,
		intervalS:   intervalS,
	}

	if len(tasks) > 0 {
		s.tasks = tasks
		heap.Init(&s.tasks)
	}

	return s
}

// NewTodayTime 通过value指定当天的具体时间，格式为 "hh:mm:ss"
//  e.g
//  NewToday("10:30:00")
func NewTodayTime(value string) (int64, error) {
	const layout = "23:59:59"

	if len(value) != len(layout) || value[2] != ':' || value[5] != ':' {
		return 0, errors.New("time format is wrong")
	}

	hour, err := strconv.Atoi(value[:2])
	if err != nil || hour < 0 || hour > 23 {
		return 0, errors.New("time format is wrong")
	}

	min, err := strconv.Atoi(value[3:5])
	if err != nil || min < 0 || min > 59 {
		return 0, errors.New("time format is wrong")
	}

	sec, err := strconv.Atoi(value[6:8])
	if err != nil || sec < 0 || sec > 59 {
		return 0, errors.New("time format is wrong")
	}

	now := time.Now().Format(time.RFC3339)
	value = strings.Replace(now, now[11:19], value, 1)

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}

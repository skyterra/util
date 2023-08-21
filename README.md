# util 工具库

## 控制 http 并发请求数
提供并发请求控制

```go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	util "github.com/skyterra/util/http"
)

func GinMiddleware(handle *util.MaxClientsHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		f := func(w http.ResponseWriter, r *http.Request) {
			ctx.Next()
		}

		handle.Middleware(f)(ctx.Writer, ctx.Request)
	}
}

func main() {
	maxClientsHandler := util.NewMaxClientsHandler(10, 300*1e3)
	engine := gin.New()

	engine.Use(GinMiddleware(maxClientsHandler))

	engine.Handle(http.MethodGet, "/echo", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, maxClientsHandler.Stats())
	})

	engine.Run(":9900")
}

// http response:
//
//  {
// 		"throttles": true,
// 		"request_incoming": 1,
// 		"request_in_queue": 0,
// 		"request_in_processing": 1,
// 		"request_done": 0,
// 		"request_wait_timeout": 0,
// 		"request_cancel": 0
//  }


```

## schedule
提供基于时间点任务调度器，用户可以对任务设置执行时间，当时间到达时，调度器会自动执行该任务

```go

package main

import (
	"fmt"
	"strconv"
	"time"
	"util/schedule"
)

type DemoTask struct {
	timestamp int64
}

func (r *DemoTask) RunAt() int64 {
	return r.timestamp
}

func (r *DemoTask) Run(s *schedule.TimingSchedule) {
	t := time.Unix(r.timestamp, 0)
	fmt.Printf("do demo task. runAt:%s now:%s\n", t.String(), time.Now())
}

func (r *DemoTask) OnError(err error) {

}

type ErrTask struct {
	timestamp int64
}

func (r *ErrTask) RunAt() int64 {
	return r.timestamp
}

func (r *ErrTask) Run(s *schedule.TimingSchedule) {
	panic("panic on run ErrTask")
}

func (r *ErrTask) OnError(err error) {
	fmt.Println(err)
	panic("panic on error")
}

func NewDemoTask(time string) schedule.ITimingTask {
	t, _ := schedule.NewTodayTime(time)
	return &DemoTask{timestamp: t}
}

func NewErrTask(time string) schedule.ITimingTask {
	t, _ := schedule.NewTodayTime(time)
	return &ErrTask{timestamp: t}
}

func GetNextSecondTime() string {
	now := time.Now().Format(time.RFC3339)
	hour, _ := strconv.Atoi(now[11:13])
	min, _ := strconv.Atoi(now[14:16])
	second, _ := strconv.Atoi(now[17:19])

	if second == 59 && min == 59 {
		hour++
		return fmt.Sprintf("%02d:%02d:%02d", hour, min, second)
	}

	if second == 59 {
		min++
		return fmt.Sprintf("%02d:%02d:%02d", hour, min, second)
	}

	second++
	return fmt.Sprintf("%02d:%02d:%02d", hour, min, second)
}

func main() {
	s := schedule.NewTimingSchedule(2, 1, NewErrTask(GetNextSecondTime()))
	for i := 0; i < 5; i++ {
		s.Push(NewDemoTask(GetNextSecondTime()))
	}

	s.Start()

	time.Sleep(3 * time.Second)
}


```

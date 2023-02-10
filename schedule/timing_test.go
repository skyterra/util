package schedule

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	. "github.com/onsi/gomega"

	. "github.com/onsi/ginkgo"
)

type DemoTask struct {
	timestamp int64
}

func (r *DemoTask) RunAt() int64 {
	return r.timestamp
}

func (r *DemoTask) Run(s *TimingSchedule) {
	t := time.Unix(r.timestamp, 0)
	fmt.Printf("do demo task. runat:%s now:%s\n", t.String(), time.Now())
	Expect(time.Now().Sub(t) >= 0).Should(BeTrue())
}

func (r *DemoTask) OnError(err error) {

}

type ErrTask struct {
	timestamp int64
}

func (r *ErrTask) RunAt() int64 {
	return r.timestamp
}

func (r *ErrTask) Run(s *TimingSchedule) {
	panic("panic on run task")
}

func (r *ErrTask) OnError(err error) {
	Expect(err).ShouldNot(Succeed())
	fmt.Println(err)
	panic("panic on error")
}

func NewDemoTask(time string) ITimingTask {
	t, _ := NewTodayTime(time)
	return &DemoTask{timestamp: t}
}

func NewErrTask(time string) ITimingTask {
	t, _ := NewTodayTime(time)
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

var _ = Describe("Timing", func() {
	Context("timing", func() {
		It("time parse", func() {
			v, err := NewTodayTime("18:53:00")
			Expect(err).Should(Succeed())

			t := time.Unix(v, 0)

			Expect(t.Hour() == 18).Should(BeTrue())
			Expect(t.Minute() == 53).Should(BeTrue())
			Expect(t.Second() == 0).Should(BeTrue())

			v, err = NewTodayTime("12：30：00")
			Expect(err).ShouldNot(Succeed())

			v, err = NewTodayTime("24:00:00")
			Expect(err).ShouldNot(Succeed())

			v, err = NewTodayTime("23:60:00")
			Expect(err).ShouldNot(Succeed())

			v, err = NewTodayTime("23:59:60")
			Expect(err).ShouldNot(Succeed())

		})

		It("push and pop with timestamp sort", func() {
			s := NewTimingSchedule(2, 1)
			for i := 0; i < 100; i++ {
				s.Push(&DemoTask{
					timestamp: rand.Int63n(10000),
				})
			}
			Expect(s.tasks.Len() == 100).Should(BeTrue())

			pre := s.pop()
			fmt.Printf("%d ", pre.RunAt())
			for i := 1; i < 100; i++ {
				cur := s.pop()
				Expect(pre.RunAt() < cur.RunAt())

				fmt.Printf("%d ", cur.RunAt())
			}
			fmt.Println()
		})

		It("concurrent push & pop", func() {
			s := NewTimingSchedule(2, 1)

			wg := sync.WaitGroup{}
			wg.Add(1)

			count := 3
			wg2 := sync.WaitGroup{}
			wg2.Add(count)

			for i := 0; i < count; i++ {
				go func(n int) {
					wg.Wait()
					for i := 0; i < 10; i++ {
						time := fmt.Sprintf("10:00:%02d", i+n*10)
						s.Push(NewDemoTask(time))
					}

					wg2.Done()
				}(i)
			}

			wg.Done()
			wg2.Wait()
			Expect(s.Len() == count*10)

			wg3 := sync.WaitGroup{}
			wg3.Add(count)

			for i := 0; i < count; i++ {
				go func(n int) {
					for i := 0; i < 10; i++ {
						t := s.pop()
						Expect(t != nil).Should(BeTrue())
					}

					wg3.Done()
				}(i)
			}

			wg3.Wait()
			Expect(s.Len() == 0).Should(BeTrue())
		})

		It("start & stop", func() {
			s := NewTimingSchedule(2, 1, NewErrTask(GetNextSecondTime()))
			for i := 0; i < 5; i++ {
				s.Push(NewDemoTask(GetNextSecondTime()))
			}

			s.Start()

			go func() {
				for i := 0; i < 100; i++ {
					s.Push(NewDemoTask(GetNextSecondTime()))
					time.Sleep(10 * time.Millisecond)
				}
			}()

			time.Sleep(3 * time.Second)
			s.Shutdown()
			s.Push(NewDemoTask(GetNextSecondTime()))
			Expect(s.Len() == 0).Should(BeTrue())
		})
	})

})

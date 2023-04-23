package primitive_test

import (
	"fmt"
	"sync"
	"time"
	"util/primitive"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lru", func() {
	Context("lru put", func() {
		It("should be succeed.", func() {
			lru := primitive.NewLRU(10, 200)
			for i := 0; i < 100; i++ {
				lru.Put(i, fmt.Sprintf("value:%d", i))
			}

			Expect(lru.GetLength() == 10).Should(BeTrue())
			for i := 90; i < 100; i++ {
				v, exist := lru.Get(i)
				Expect(exist).Should(BeTrue())
				Expect(v.(string) == fmt.Sprintf("value:%d", i)).Should(BeTrue())
			}
		})

		It("concurrent put", func() {
			lru := primitive.NewLRU(10, -1)

			wg := sync.WaitGroup{}
			wg.Add(1)

			wg2 := sync.WaitGroup{}
			wg2.Add(3)

			const count = 1e5

			go func() {
				wg.Wait()
				for i := 0; i < count; i++ {
					lru.Put(i, fmt.Sprintf("value:%d", i))
				}
				wg2.Done()
			}()

			go func() {
				wg.Wait()
				for i := 0; i < count; i++ {
					lru.Put(i, fmt.Sprintf("value:%d", i))
				}
				wg2.Done()
			}()

			go func() {
				wg.Wait()
				for i := 0; i < count; i++ {
					v, exist := lru.Get(i)
					if exist {
						fmt.Println(v)
					}
				}
				wg2.Done()
			}()

			wg.Done()

			wg2.Wait()
			Expect(lru.GetLength() == 10).Should(BeTrue())
		})

		It("ttl", func() {
			lru := primitive.NewLRU(10, 10)

			lru.Put(1, "hello")
			time.Sleep(11 * time.Millisecond)
			_, exist := lru.Get(1)
			Expect(exist).Should(BeFalse())

			lru.Put(10, "world")
			time.Sleep(8 * time.Millisecond)
			v, exist := lru.Get(10)
			Expect(exist && (v.(string) == "world")).Should(BeTrue())
		})

		It("disable ttl", func() {
			lru := primitive.NewLRU(10, -1)
			lru.Put(10, "world")
			time.Sleep(10 * time.Millisecond)
			v, exist := lru.Get(10)
			Expect(exist && (v.(string) == "world")).Should(BeTrue())
		})

		It("lru", func() {
			lru := primitive.NewLRU(3, -1)
			lru.Put(1, "a")
			lru.Put(2, "b")
			lru.Put(3, "c")

			_, exist := lru.Get(1)
			Expect(exist).Should(BeTrue())

			lru.Put(5, "e")
			lru.Put(6, "f")

			_, exist = lru.Get(2)
			Expect(exist).Should(BeFalse())

			_, exist = lru.Get(3)
			Expect(exist).Should(BeFalse())
		})
	})
})

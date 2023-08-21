package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MaxClients", func() {
	Context("NewMaxClientsHandler", func() {
		It("should be succeed", func() {
			opts := MaxClientsOpts{}
			opts.SetTimeoutResponse([]byte("this is timeout response"))
			opts.SetTimeoutStatusCode(http.StatusOK)

			h := NewMaxClientsHandler(10, 3000)
			Expect(h.pool != nil).Should(BeTrue())
			Expect(h.maxClients == 10).Should(BeTrue())
			Expect(h.deadlineMS == 3000).Should(BeTrue())

			data, _ := json.MarshalIndent(h.Stats(), "", " ")
			fmt.Println(string(data))

			h1 := NewMaxClientsHandler(10, 3000, opts)
			Expect(h1.waitTimeoutStatusCode == http.StatusOK).Should(BeTrue())
			fmt.Println(string(h1.waitTimeoutResponse))
		})
	})

	FContext("middleware", func() {
		It("should be succeed", func() {
			h := NewMaxClientsHandler(5, 10)

			ts := httptest.NewServer(h.Middleware(func(w http.ResponseWriter, r *http.Request) {
				n := rand.Intn(20)
				time.Sleep(time.Duration(n) * time.Millisecond)
				fmt.Fprint(w, "Hello, client")
			}))

			defer ts.Close()

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				for i := 0; i < 1000; i++ {
					data, _ := json.MarshalIndent(h.Stats(), "", " ")
					fmt.Println(string(data))
					time.Sleep(time.Millisecond)
				}
				wg.Done()
			}()

			for i := 0; i < 100; i++ {
				go func() {
					resp, _ := http.Get(ts.URL)
					data, _ := ioutil.ReadAll(resp.Body)
					fmt.Println(string(data))
				}()

				n := rand.Intn(3)
				time.Sleep(time.Duration(n) * time.Millisecond)
			}

			wg.Wait()
		})
	})
})

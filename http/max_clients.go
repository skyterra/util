package http

import (
	"net/http"
	"sync/atomic"
	"time"
)

const (
	defaultDeadlineMS        = 60 * 1e3 // 默认请求等待超时时间60s
	defaultTimeoutStatusCode = http.StatusServiceUnavailable
	defaultTimeoutResp       = "Deadline exceeded while waiting in incoming queue, please reduce your request rate"
)

type MaxClientsStatus struct {
	Throttles           bool   `json:"throttles"`             // 是否开启限流阀门
	RequestIncoming     uint64 `json:"request_incoming"`      // 收到请求数
	RequestInQueue      int32  `json:"request_in_queue"`      // 等待请求数
	RequestInProcessing int32  `json:"request_in_processing"` // 处理中的请求数
	RequestDone         uint64 `json:"request_done"`          // 处理完成的请求数
	RequestWaitTimeout  uint64 `json:"request_wait_timeout"`  // 等待超时请求数
	RequestCancel       uint64 `json:"request_cancel"`        // 客户端取消请求数
}

type MaxClientsOpts struct {
	waitTimeoutStatusCode *int    // 设置等待超时错误码
	waitTimeoutResponse   *[]byte // 设置等待超时response
}

type MaxClientsHandler struct {
	maxClients uint
	deadlineMS uint

	requestInQueue     int32  // 统计等待中的请求数
	requestIncoming    uint64 // 统计收到的请求数
	requestDone        uint64 // 统计处理完成请求数
	requestWaitTimeout uint64 // 统计等待超时请求数
	requestCancel      uint64 // 统计取消请求数

	waitTimeoutStatusCode int    // 请求等待超时返回错误码
	waitTimeoutResponse   []byte // 请求等待超时返回的response

	pool chan struct{}
}

func (opts *MaxClientsOpts) SetTimeoutStatusCode(code int) {
	opts.waitTimeoutStatusCode = &code
}

func (opts *MaxClientsOpts) SetTimeoutResponse(data []byte) {
	temp := make([]byte, len(data))
	copy(temp, data)

	opts.waitTimeoutResponse = &temp
}

// MaxClientsHandler 流控状态信息
func (mc *MaxClientsHandler) Stats() *MaxClientsStatus {
	stat := &MaxClientsStatus{
		RequestIncoming:     atomic.LoadUint64(&mc.requestIncoming),
		RequestInQueue:      atomic.LoadInt32(&mc.requestInQueue),
		RequestInProcessing: int32(len(mc.pool)),
		RequestDone:         atomic.LoadUint64(&mc.requestDone),
		RequestWaitTimeout:  atomic.LoadUint64(&mc.requestWaitTimeout),
		RequestCancel:       atomic.LoadUint64(&mc.requestCancel),
		Throttles:           mc.pool != nil,
	}

	return stat
}

// Middleware 控制最大并发请求数中间件
func (mc *MaxClientsHandler) Middleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&mc.requestIncoming, 1)

		// 未开启限流控制
		if mc.pool == nil {
			f.ServeHTTP(w, r)
			return
		}

		atomic.AddInt32(&mc.requestInQueue, 1)

		deadlineTimer := time.NewTimer(time.Duration(mc.deadlineMS) * time.Millisecond)
		defer deadlineTimer.Stop()

		select {
		case mc.pool <- struct{}{}: // 处理当前请求
			defer func() {
				<-mc.pool // 请求处理完成后记得出队
				atomic.AddUint64(&mc.requestDone, 1)
			}()

			atomic.AddInt32(&mc.requestInQueue, -1)
			f.ServeHTTP(w, r)

		case <-deadlineTimer.C: // 请求等待超时
			w.WriteHeader(mc.waitTimeoutStatusCode)
			w.Write(mc.waitTimeoutResponse)

			atomic.AddInt32(&mc.requestInQueue, -1)
			atomic.AddUint64(&mc.requestWaitTimeout, 1)

		case <-r.Context().Done(): // 客户端中断请求
			w.WriteHeader(499)

			atomic.AddInt32(&mc.requestInQueue, -1)
			atomic.AddUint64(&mc.requestCancel, 1)
		}
	}
}

// NewMaxClientsHandler 控制最大并发连接数；maxClients 指服务可以同时处理的最大请求数量，0表示没有限制；
// deadlineMS 指当并发请求数已达上限后，后续请求的最长等待时间(毫秒)，0表示使用默认值（60s）
func NewMaxClientsHandler(maxClients, deadlineMS uint, opts ...MaxClientsOpts) *MaxClientsHandler {
	handler := &MaxClientsHandler{
		maxClients:            maxClients,
		deadlineMS:            deadlineMS,
		waitTimeoutStatusCode: defaultTimeoutStatusCode,
		waitTimeoutResponse:   []byte(defaultTimeoutResp),
	}

	if maxClients > 0 {
		handler.pool = make(chan struct{}, maxClients)
	}

	if deadlineMS == 0 {
		handler.deadlineMS = defaultDeadlineMS
	}

	if len(opts) > 0 {
		opt := opts[0]

		if opt.waitTimeoutStatusCode != nil {
			handler.waitTimeoutStatusCode = *opt.waitTimeoutStatusCode
		}

		if opt.waitTimeoutResponse != nil {
			handler.waitTimeoutResponse = *opt.waitTimeoutResponse
		}
	}

	return handler
}

package requester

import (
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/edpryk/buter/internal/modules/stability"
)

/*
Need to created requests queue with concurrent restriction
Need 1 buff channel with cap like a max concurrent requests
And signal channel that will block provider till the
N-request done and may be wait group
*/
type RequestQueueConfig struct {
	MaxConcurrentRequests int
	Retries               int
	RetryDelay            int
}

type CustomResponse struct {
	Duration time.Duration
	Payloads []string
	http.Response
}

type ReuqestParameters struct {
	Method   string
	Url      string
	Header   map[string]string
	Body     io.Reader
	Payloads []string
}

type RequestQueue struct {
	receiveQueue   chan ReuqestParameters
	sendQueue      chan CustomResponse
	errQueue       chan error
	spawnedThreads *stability.Counter
	// spawnedRequets *stability.Counter

	wg    *sync.WaitGroup
	block chan int

	CloseSig chan int
	RequestQueueConfig
}

func (rq *RequestQueue) StartWorker() (consumer chan ReuqestParameters, provider chan CustomResponse, errQ chan error) {
	go func() {
		for requestParameter := range rq.receiveQueue {
			// Blocking more than N threads per period ?
			if (rq.spawnedThreads.Number()) >= rq.MaxConcurrentRequests {
				<-rq.block
			}

			requst := func() (any, error) {
				return Do(
					requestParameter.Method,
					requestParameter.Url,
					requestParameter.Header,
					requestParameter.Body,
				)
			}

			rq.wg.Add(1)
			rq.spawnedThreads.Increment()
			go func(params ReuqestParameters) {
				defer rq.wg.Done()
				defer func() {
					// Unblock receiving
					if rq.spawnedThreads.Number() <= 0 {
						rq.block <- 0
					}
				}()
				defer rq.spawnedThreads.Decrement()

				startTime := time.Now()
				res, err := stability.Retry(requst, rq.Retries, rq.RetryDelay)
				if err != nil {
					rq.errQueue <- err
				} else {
					rq.sendQueue <- CustomResponse{
						Response: res.(http.Response),
						Duration: time.Since(startTime),
						Payloads: params.Payloads,
					}
				}
			}(requestParameter)
		}

		rq.wg.Wait()
		close(rq.sendQueue)
	}()

	return rq.receiveQueue, rq.sendQueue, rq.errQueue
}

func NewRequestQueue(config RequestQueueConfig) *RequestQueue {
	return &RequestQueue{
		receiveQueue:   make(chan ReuqestParameters, config.MaxConcurrentRequests),
		sendQueue:      make(chan CustomResponse, config.MaxConcurrentRequests),
		errQueue:       make(chan error, config.MaxConcurrentRequests),
		spawnedThreads: stability.NewCounter(),
		// spawnedRequets: stability.NewCounter(),
		wg:    &sync.WaitGroup{},
		block: make(chan int),

		RequestQueueConfig: config,
	}
}

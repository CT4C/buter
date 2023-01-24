package requester

import (
	"fmt"
	"io"
	"net/http"
	"sync"

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

type ReuqestParameters struct {
	Method string
	Url    string
	Header map[string]string
	Body   io.Reader
}

type RequestQueue struct {
	receiveQueue   chan ReuqestParameters
	sendQueue      chan http.Response
	errQueue       chan error
	requestCounter *stability.Counter
	wg             *sync.WaitGroup
	block          chan int

	RequestQueueConfig
}

func (rq *RequestQueue) StartWorker() (consumer chan ReuqestParameters, provider chan http.Response, errQ chan error) {
	go func() {
		for requestParameter := range rq.receiveQueue {
			defer rq.requestCounter.Increment()
			fmt.Println(len(rq.receiveQueue))
			requst := func() (any, error) {
				return Do(
					requestParameter.Method,
					requestParameter.Url,
					requestParameter.Header,
					requestParameter.Body,
				)
			}

			rq.wg.Add(1)
			go func() {
				defer rq.wg.Done()
				// <-rq.block

				res, err := stability.Retry(requst, rq.Retries, rq.RetryDelay)
				if err != nil {
					rq.errQueue <- err
				} else {
					rq.sendQueue <- res.(http.Response)
				}
			}()

			// if (rq.requestCounter.Number() % rq.MaxConcurrentRequests) == 0 {
			// 	rq.block <- 0
			// }
		}

		rq.wg.Wait()
	}()

	return rq.receiveQueue, rq.sendQueue, rq.errQueue
}

func NewRequestQueue(config RequestQueueConfig) *RequestQueue {
	return &RequestQueue{
		receiveQueue:   make(chan ReuqestParameters, config.MaxConcurrentRequests),
		sendQueue:      make(chan http.Response, config.MaxConcurrentRequests),
		errQueue:       make(chan error, config.MaxConcurrentRequests),
		requestCounter: stability.NewCounter(),
		wg:             &sync.WaitGroup{},
		block:          make(chan int),

		RequestQueueConfig: config,
	}
}

package requester

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/edpryk/buter/lib/stability"
)

type LimitedQConfig struct {
	MaxThreads int
	Delay      int
	Retries    int
	ResponseQ  chan CustomResponse
	ErrQ       chan error
	Ctx        context.Context
}

type LimitedQueue struct {
	q  chan ReuqestParameters
	wg *sync.WaitGroup

	LimitedQConfig
}

func (lm *LimitedQueue) IsFull() bool {
	return len(lm.q) == lm.MaxThreads
}

func (lm *LimitedQueue) IsNotEmpty() bool {
	return len(lm.q) > 0
}

func (lm *LimitedQueue) Receive(rp ReuqestParameters) {
	lm.q <- rp
}

func (lm *LimitedQueue) reNew() {
	lm.q = make(chan ReuqestParameters, lm.MaxThreads)
}

func (lm *LimitedQueue) Proceed() {
	/*
		Need to add context handling
	*/
	lm.wg.Add(1)
	go func() {
		defer lm.wg.Done()

		ticker := time.NewTicker(time.Duration(lm.Delay) * time.Millisecond)
		for parameters := range lm.q {
			requstCaller := func() (any, error) {
				return Do(
					parameters.Method,
					parameters.Url,
					parameters.Header,
					parameters.Body,
				)
			}

			lm.wg.Add(1)
			go func(params ReuqestParameters) {
				defer lm.wg.Done()

				startTime := time.Now()
				res, err := stability.Retry(requstCaller, lm.Retries, lm.Delay)
				if err != nil {
					lm.ErrQ <- err
				} else {
					lm.ResponseQ <- CustomResponse{
						Response: res.(http.Response),
						Duration: time.Since(startTime),
						Payloads: params.Payloads,
					}
				}
			}(parameters)

			<-ticker.C
		}
	}()

	close(lm.q)
	lm.wg.Wait()
	lm.reNew()
}

func (lm *LimitedQueue) ProceedIFFull() {
	if lm.IsFull() {
		lm.Proceed()
	}
}

func (lm *LimitedQueue) ProceedIFNotFull() {
	if lm.IsNotEmpty() && !lm.IsFull() {
		lm.Proceed()
	}
}

func NewLimitedQ(config LimitedQConfig) LimitedQueue {
	return LimitedQueue{
		LimitedQConfig: config,

		q:  make(chan ReuqestParameters, config.MaxThreads),
		wg: &sync.WaitGroup{},
	}
}

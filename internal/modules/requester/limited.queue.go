package requester

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/edpryk/buter/lib/stability"
)

type LimitedQConfig struct {
	MaxThreads  int
	Delay       int
	RetrayDelay int
	Retries     int
	ResponseQ   chan CustomResponse
	ErrQ        chan error
	Ctx         context.Context
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
				reader := strings.NewReader(parameters.Body.String())

				// TODO: move to seprated func
				defaultHeaders := make(map[string]string)

				if parameters.Method == http.MethodPost {
					defaultHeaders["Content-Length"] = fmt.Sprintf("%d", len(parameters.Body.String()))
				}

				for key := range parameters.Header {
					defaultHeaders[key] = parameters.Header[key]
				}

				return Do(
					parameters.Method,
					parameters.Url,
					defaultHeaders,
					reader,
				)
			}

			lm.wg.Add(1)
			go func(params ReuqestParameters) {
				defer lm.wg.Done()
				startTime := time.Now()
				res, err := stability.Retry(requstCaller, lm.Retries, lm.RetrayDelay)

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

		ticker.Stop()
	}()
	// close(lm.q)
	lm.wg.Wait()
	lm.reNew()
}

func (lm *LimitedQueue) ProceedIFFull() {
	if lm.IsFull() {
		close(lm.q)
		lm.Proceed()
	}
}

func (lm *LimitedQueue) ProceedIFNotFull() {
	if lm.IsNotEmpty() && !lm.IsFull() {
		close(lm.q)
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

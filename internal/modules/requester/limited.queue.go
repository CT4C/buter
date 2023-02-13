package requester

import (
	"context"
	"log"
	"sync"
	"time"
)

type LimitedQConfig struct {
	Ctx        context.Context
	ErrQ       chan error
	Delay      int
	Retries    int
	RetryDelay int
	MaxThreads int
	ResponseQ  chan CustomResponse
}

type LimitedQueue struct {
	requestsQ chan RequestParameters
	wg        *sync.WaitGroup

	LimitedQConfig
}

func (lq *LimitedQueue) IsFull() bool {
	return len(lq.requestsQ) == lq.MaxThreads
}

func (lq *LimitedQueue) IsNotEmpty() bool {
	return len(lq.requestsQ) > 0
}

func (lq *LimitedQueue) Receive(rp RequestParameters) {
	lq.requestsQ <- rp
}

func (lq *LimitedQueue) reNew() {
	lq.requestsQ = make(chan RequestParameters, lq.MaxThreads)
}

func (lq *LimitedQueue) CLose() {}

func (lq *LimitedQueue) Proceed() {
	lq.wg.Add(1)
	go func() {
		defer lq.wg.Done()

		ticker := time.NewTicker(time.Duration(lq.Delay) * time.Millisecond)
		for parameters := range lq.requestsQ {
			resCh, errCh := AsyncRequestWithRetry(parameters, lq.Retries, lq.RetryDelay)

			lq.wg.Add(1)
			go func() {
				defer lq.wg.Done()

				select {
				case res := <-resCh:
					lq.ResponseQ <- res
					return
				case err := <-errCh:
					lq.ErrQ <- err
					return
				case <-lq.Ctx.Done():
					log.Printf("Request Canceled\n")
					return
				}
			}()

			select {
			case <-lq.Ctx.Done():
				log.Println("LimitedQ Canceled")
				return
			case <-ticker.C:
			}
		}

		ticker.Stop()
	}()
	lq.wg.Wait()
}

func (lq *LimitedQueue) ProceedIFFull() {
	if lq.IsFull() {
		close(lq.requestsQ)
		lq.Proceed()
		lq.reNew()
	}
}

func (lq *LimitedQueue) ProceedIFNotFull() {
	if lq.IsNotEmpty() && !lq.IsFull() {
		close(lq.requestsQ)
		lq.Proceed()
		lq.reNew()
	}
}

func NewLimitedQ(config LimitedQConfig) LimitedQueue {
	return LimitedQueue{
		LimitedQConfig: config,

		requestsQ: make(chan RequestParameters, config.MaxThreads),
		wg:        &sync.WaitGroup{},
	}
}

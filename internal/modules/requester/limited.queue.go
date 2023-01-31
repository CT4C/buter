package requester

import (
	"context"
	"log"
	"sync"
	"time"
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
	requetsQ chan ReuqestParameters
	wg       *sync.WaitGroup

	LimitedQConfig
}

func (lq *LimitedQueue) IsFull() bool {
	return len(lq.requetsQ) == lq.MaxThreads
}

func (lq *LimitedQueue) IsNotEmpty() bool {
	return len(lq.requetsQ) > 0
}

func (lq *LimitedQueue) Receive(rp ReuqestParameters) {
	lq.requetsQ <- rp
}

func (lq *LimitedQueue) reNew() {
	lq.requetsQ = make(chan ReuqestParameters, lq.MaxThreads)
}

func (lq *LimitedQueue) Proceed() {
	lq.wg.Add(1)
	go func() {
		defer lq.wg.Done()

		ticker := time.NewTicker(time.Duration(lq.Delay) * time.Millisecond)
		for parameters := range lq.requetsQ {
			resCh, errCh := AsyncRequestWitnRetry(parameters, lq.Retries, lq.RetrayDelay)

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
	lq.reNew()
}

func (lq *LimitedQueue) ProceedIFFull() {
	if lq.IsFull() {
		close(lq.requetsQ)
		lq.Proceed()
	}
}

func (lq *LimitedQueue) ProceedIFNotFull() {
	if lq.IsNotEmpty() && !lq.IsFull() {
		close(lq.requetsQ)
		lq.Proceed()
	}
}

func NewLimitedQ(config LimitedQConfig) LimitedQueue {
	return LimitedQueue{
		LimitedQConfig: config,

		requetsQ: make(chan ReuqestParameters, config.MaxThreads),
		wg:       &sync.WaitGroup{},
	}
}

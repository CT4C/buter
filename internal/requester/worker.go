package requester

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Stringer interface {
	String() string
}

type QueueWorkerConfig struct {
	Ctx                   context.Context
	Delay                 int
	Retries               int
	RetryDelay            int
	MaxConcurrentRequests int
}

type CustomResponse struct {
	Duration time.Duration
	Payloads []string
	http.Response
}

type RequestParameters struct {
	Url      string
	Body     Stringer
	Method   string
	Header   map[string]string
	Payloads []string
}

type QueueWorker struct {
	errQueue  chan error
	requestQ  chan RequestParameters
	responseQ chan CustomResponse

	QueueWorkerConfig
}

func (rq *QueueWorker) Run() (reqConsumer chan RequestParameters, resProvider chan CustomResponse, errQ chan error) {
	go func() {
		limitedQ := NewLimitedQ(LimitedQConfig{
			MaxThreads: rq.MaxConcurrentRequests,
			Delay:      rq.Delay,
			Retries:    rq.Retries,
			RetryDelay: rq.RetryDelay,
			ResponseQ:  rq.responseQ,
			ErrQ:       rq.errQueue,
			Ctx:        rq.Ctx,
		})

		allowRun := true
		for allowRun == true {
			select {
			case requestParameters, ok := <-rq.requestQ:
				limitedQ.ProceedIFFull()

				if !ok {
					allowRun = false
					limitedQ.ProceedIFNotFull()
					break
				}

				limitedQ.Receive(requestParameters)

			case <-rq.Ctx.Done():
				log.Println("Worked Canceled")
				allowRun = false
				break
			}
		}

		close(rq.responseQ)
	}()

	return rq.requestQ, rq.responseQ, rq.errQueue
}

func NewRequestQueue(config QueueWorkerConfig) *QueueWorker {
	return &QueueWorker{
		errQueue:  make(chan error, config.MaxConcurrentRequests),
		requestQ:  make(chan RequestParameters, config.MaxConcurrentRequests),
		responseQ: make(chan CustomResponse, config.MaxConcurrentRequests),

		QueueWorkerConfig: config,
	}
}

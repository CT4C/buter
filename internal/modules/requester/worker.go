package requester

import (
	"context"
	"io"
	"net/http"
	"time"
)

type QueueWorkerConfig struct {
	MaxConcurrentRequests int
	Retries               int
	Delay                 int
	Ctx                   context.Context
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

type QueueWorker struct {
	receiveQueue chan ReuqestParameters
	sendQueue    chan CustomResponse
	errQueue     chan error

	QueueWorkerConfig
}

func (rq *QueueWorker) Run() (reqConsumer chan ReuqestParameters, resProvider chan CustomResponse, errQ chan error) {
	go func() {
		limitedQ := NewLimitedQ(LimitedQConfig{
			MaxThreads: rq.MaxConcurrentRequests,
			Delay:      rq.Delay,
			Retries:    rq.Retries,
			ResponseQ:  rq.sendQueue,
			ErrQ:       rq.errQueue,
		})

		for allowRun := true; allowRun == true; {
			select {
			case requestParameters, ok := <-rq.receiveQueue:
				if !ok {
					allowRun = false
					limitedQ.ProceedIFNotFull()
					break
				}

				limitedQ.ProceedIFFull()
				limitedQ.Receive(requestParameters)

			case <-rq.Ctx.Done():
				allowRun = false
				break
			}
		}

		close(rq.sendQueue)
	}()

	return rq.receiveQueue, rq.sendQueue, rq.errQueue
}

func NewRequestQueue(config QueueWorkerConfig) *QueueWorker {
	return &QueueWorker{
		receiveQueue: make(chan ReuqestParameters, config.MaxConcurrentRequests),
		sendQueue:    make(chan CustomResponse, config.MaxConcurrentRequests),
		errQueue:     make(chan error, config.MaxConcurrentRequests),

		QueueWorkerConfig: config,
	}
}

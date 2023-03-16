package payload

import (
	"sync"

	"github.com/edpryk/buter/internal/helpers/prepare"
	"github.com/edpryk/buter/pkg/requester"
)

type ConsumerPayload struct {
	requestQueue chan requester.RequestParameters
	errQueue     chan error
	httpMethod   string
	wg           *sync.WaitGroup
}

func NewPayloadConsumer(requestQueue chan requester.RequestParameters, errQueue chan error, httpMethod string) ConsumerPayload {
	consumer := ConsumerPayload{
		httpMethod:   httpMethod,
		requestQueue: requestQueue,
		errQueue:     errQueue,
		wg:           &sync.WaitGroup{},
	}

	return consumer
}

func (consumer ConsumerPayload) Consume(updatedPayloadValue string, payloads []string, err error) {
	if err != nil {
		// TODO
		consumer.errQueue <- err
		return
	}

	parsedAttackValue, err := prepare.ParseAttackValue(updatedPayloadValue)
	if err != nil {
		consumer.errQueue <- err
		return
	}

	consumer.requestQueue <- requester.RequestParameters{
		Url:      parsedAttackValue.Url,
		Body:     parsedAttackValue.Body,
		Method:   consumer.httpMethod,
		Header:   parsedAttackValue.Headers,
		Payloads: payloads,
	}
}

func (consumer ConsumerPayload) Close() {
	close(consumer.requestQueue)
	close(consumer.errQueue)
}

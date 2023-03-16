package buter

import (
	"context"
)

type Dos struct {
	ctx         context.Context
	proceeded   int
	maxRequests int
	attackValue string
}

func NewDos(ctx context.Context, attackValue string, maxRequests int) Dos {
	d := Dos{
		ctx:         ctx,
		attackValue: attackValue,
		maxRequests: maxRequests,
	}

	return d
}

func (d Dos) Proceeded() int {
	return d.proceeded
}

func (d Dos) ProducePayload(payloadConsumer PayloadConsumer) chan int {
	end := make(chan int, 0)
	go func() {
		for i := 0; i < d.maxRequests; i++ {
			payloadConsumer.Consume(d.attackValue, []string{}, nil)
			// TODO: ctx
		}
		end <- 0
	}()

	return end
}

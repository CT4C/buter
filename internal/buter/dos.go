package buter

import (
	"context"
)

type DosConfig struct {
	Ctx         context.Context
	Url         string
	Body        string
	Headers     map[string]string
	MaxRequests int
}

type Dos struct {
	proceeded int
	DosConfig
}

func NewDos(config DosConfig) Dos {
	d := Dos{
		DosConfig: config,
	}

	return d
}

func (d Dos) Proceeded() int {
	return d.proceeded
}

func (d Dos) ProducePayload(payloadConsumer chan CraftedPayload) chan int {
	end := make(chan int, 0)
	go func() {
		for i := 0; i < d.MaxRequests; i++ {
			payloadConsumer <- CraftedPayload{
				Url:      d.Url,
				Body:     d.Body,
				Headers:  d.Headers,
				Payloads: nil,
			}
		}

		// close(payloadConsumer)
		end <- 0
	}()

	return end
}

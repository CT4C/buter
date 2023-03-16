package buter

import (
	"context"
)

type Sniper struct {
	ctx               context.Context
	payloadNode       *PayloadNode
	attackValue       string
	producedPayloads  int
	workingPayloadSet []string
}

func (s *Sniper) Proceeded() int {
	return s.producedPayloads
}

func (s *Sniper) ProducePayload(payloadConsumer chan CraftedPayload) chan int {
	endChan := make(chan int)
	defer func() {
		endChan <- 0
	}()
	defer close(payloadConsumer)

	s.producedPayloads += buildPayload(s.attackValue, s.payloadNode, s.workingPayloadSet, payloadConsumer)

	return endChan
}

func NewSniper(ctx context.Context, attackValue string, payloadNode *PayloadNode) *Sniper {
	return &Sniper{
		ctx:               ctx,
		attackValue:       attackValue,
		payloadNode:       payloadNode,
		producedPayloads:  0,
		workingPayloadSet: make([]string, 1),
	}
}

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

func (s *Sniper) ProducePayload(payloadConsumer PayloadConsumer) chan int {
	endChan := make(chan int)
	defer func() {
		endChan <- 0
	}()

	workingPayloadSet := make([]string, 1)

	onUpdate := func(updatedTargetString string, payloadInserted string, payloadNumber int) {
		workingPayloadSet[payloadNumber] = payloadInserted

		workingPayloadSetCopy := make([]string, len(workingPayloadSet))
		copy(workingPayloadSetCopy, workingPayloadSet)

		payloadConsumer.Consume(updatedTargetString, workingPayloadSetCopy, nil)
	}

	s.producedPayloads += buildPayloadList(s.attackValue, s.payloadNode, onUpdate)

	payloadConsumer.Close()

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

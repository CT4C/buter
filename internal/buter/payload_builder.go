package buter

import (
	"github.com/edpryk/buter/internal/helpers/prepare"
)

func insertPayload(target string, payload string, positions [2]int) string {
	return target[:positions[0]] + payload + target[positions[1]:]
}

func buildPayload(target string, payloadNode *PayloadNode, workingPayloadsSet []string, payloadConsumer chan CraftedPayload) (produced int) {
	for _, payload := range payloadNode.PayloadList {
		workingPayloadsSet[payloadNode.Number] = payload

		target = insertPayload(target, payload, payloadNode.Points)
		payloadNode.CurrentPayloadIdx += 1
		payloadNode.Points[1] = payloadNode.Points[0] + len(payload)
		/*
			Because in another situation chanel has a pointer to
			the workingPayloadSet slice, and last values will changed
			when a client will read from consumer
		*/
		workingPayloadCopy := make([]string, len(workingPayloadsSet))
		copy(workingPayloadCopy, workingPayloadsSet)

		/*
			1. Send to channel
			2. Increment proceeded payload
		*/
		parsedAttackValue, _ := prepare.ParseAttackValue(target)

		payloadConsumer <- CraftedPayload{
			Url:      parsedAttackValue.Url,
			Headers:  parsedAttackValue.Headers,
			Payloads: workingPayloadCopy,
			Body:     parsedAttackValue.Body,
		}
		produced++
	}

	return
}

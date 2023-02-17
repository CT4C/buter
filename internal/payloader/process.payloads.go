package payloader

import (
	"github.com/edpryk/buter/internal/helpers/prepare"
)

func updateValue(value string, payload string, positions [2]int) string {
	return value[:positions[0]] + payload + value[positions[1]:]
}

func processPayloads(value string, payloadNode *PayloadNode, workingPayloadsSet []string, payloadConsumer chan CraftedPayload) (produced int) {
	for _, payload := range payloadNode.PayloadList {
		workingPayloadsSet[payloadNode.Number] = payload

		value = updateValue(value, payload, payloadNode.Points)
		// fmt.Println(value)
		payloadNode.CurrentPayloadIdx += 1
		payloadNode.Points[1] = payloadNode.Points[0] + len(payload)
		/*
			Because in another situation chanel has a pointer to
			the workingPayloadSet slice, and last values will change
			when a client will read from consumer
		*/
		workingPayloadCopy := make([]string, len(workingPayloadsSet))
		copy(workingPayloadCopy, workingPayloadsSet)

		/*
			1. Send to channel
			2. Increment proceeded payloader
		*/
		parsedAttackValue, _ := prepare.ParseAttackValue(value)
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

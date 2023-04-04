package buter

type onUpdate func(updatedTargetString string, payloadInserted string, payloadNumber int)

func insertPayload(targetString string, payload string, positions [2]int) string {
	return targetString[:positions[0]] + payload + targetString[positions[1]:]
}

func buildPayloadList(targetString string, payloadNode *PayloadNode, onUpdate onUpdate) (produced int) {
	for _, payload := range payloadNode.PayloadList {

		targetString = insertPayload(targetString, payload, payloadNode.Points)
		payloadNode.CurrentPayloadIdx += 1
		payloadNode.Points[1] = payloadNode.Points[0] + len(payload)
		/*
			Because in another situation chanel has a pointer to
			the workingPayloadSet slice, and last values will changed
			when a client will read from consumer
		*/
		onUpdate(targetString, payload, payloadNode.Number)
		produced++
	}

	return
}

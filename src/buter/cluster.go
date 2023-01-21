package buter

import (
	"context"
)

func Cluster(ctx context.Context, url string, provider UrlProvider, payloadNode *PayloadNode) {
	defer close(provider)

	var (
		updatedText = url
		// payloadNode        *PayloadNode = &(*payloadNode)
	)

	for payloadNode != nil && !(proceededPayloads == NeedToProceedPayloads) {

		if payloadNode.NextNode == nil {
			// Proceed last level payload in Set
			for _, payload := range payloadNode.PayloadList {
				updatedText = updatedText[:payloadNode.Points[0]] + payload + updatedText[payloadNode.Points[1]:]
				payloadNode.CurrentPayloadIdx += 1
				payloadNode.Points[1] = payloadNode.Points[0] + len(payload)
				/*
					- Send to channel
					- Increment proceeded payload
				*/
				provider <- updatedText
				proceededPayloads += 1
			}

			payloadNode.CurrentPayloadIdx = 0
			payloadNode = payloadNode.PreviousNode
			payloadNode.CurrentPayloadIdx += 1
			/*
				- Increment previous paylaod index
				- set next payloadNode to previous one
			*/
			// break
		} else {
			/*
				IF current payload index == payload list length (IS END)
				1. set next payloadNode to previous one
				2. reset current payload index
				3. Incremetn Previous payload index ?
			*/

			isEndOfCurrentPayloadProcessing := payloadNode.CurrentPayloadIdx == len(payloadNode.PayloadList)

			if isEndOfCurrentPayloadProcessing {
				payloadNode.CurrentPayloadIdx = 0
				payloadNode = payloadNode.PreviousNode
				// It's time to nexet payload in list
				payloadNode.CurrentPayloadIdx += 1
			}

			// Proceed top level payload
			(*payloadNode).WorkingPayload = payloadNode.PayloadList[payloadNode.CurrentPayloadIdx]

			updatedText = updatedText[:payloadNode.Points[0]] + payloadNode.WorkingPayload + updatedText[payloadNode.Points[1]:]

			(*payloadNode).Points[1] = payloadNode.Points[0] + len(payloadNode.WorkingPayload)

			positions := rePayloadPosition.FindAllStringSubmatchIndex(updatedText, -1)
			if len(positions) > 0 {
				payloadNode.NextNode.Points[0] = positions[0][0]
				payloadNode.NextNode.Points[1] = positions[0][1]
			}

			currentNodeCopy := *payloadNode
			payloadNode = payloadNode.NextNode
			payloadNode.PreviousNode = &currentNodeCopy
		}
	}
}

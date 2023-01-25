package payloader

type PayloadNode struct {
	Points            [2]int
	Number            int
	PayloadList       []string
	CurrentPayloadIdx int
	NextNode          *PayloadNode
	PreviousNode      *PayloadNode
	WorkingPayload    string
}

/*
Transform [][]string to Linked List
*/
func transformPayload(text string, payloadSet [][]string) (totalPayloads int, entryNode *PayloadNode, err error) {
	matchedPositions := rePayloadPosition.FindAllStringSubmatchIndex(text, -1)
	matchedPatterns := rePayloadPosition.FindAllString(text, -1)
	positionsAmount := len(matchedPositions)
	payloadsAmount := len(payloadSet)
	totalPayloads = 1

	// Validate payloads and position amount
	if positionsAmount != payloadsAmount {
		if positionsAmount < payloadsAmount {
			err = errLessPositionsThanPayloads
			return
		}
		if payloadsAmount < positionsAmount {
			err = errLessPayloadsThanPositions
			return
		}

		return
	}

	var previousNode *PayloadNode
	for number, payloads := range payloadSet {
		totalPayloads *= len(payloads)

		newNode := &PayloadNode{
			Number:            number,
			PayloadList:       payloads,
			PreviousNode:      previousNode,
			CurrentPayloadIdx: 0,
			WorkingPayload:    matchedPatterns[number],
		}
		newNode.Points = [2]int{matchedPositions[number][0], matchedPositions[number][1]}
		// Init Next Node
		if previousNode != nil && previousNode.NextNode == nil {
			previousNode.NextNode = newNode
			previousNode = newNode
		}
		// Init entry node
		if entryNode == nil {
			entryNode = newNode
			previousNode = newNode
			continue
		}
	}

	return totalPayloads, entryNode, nil
}

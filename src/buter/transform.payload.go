package buter

type PayloadNode struct {
	Points            [2]int
	Number            int
	PayloadList       []string
	CurrentPayloadIdx int
	NextNode          *PayloadNode
	PreviousNode      *PayloadNode
	WorkingPayload    string
}

var (
	NeedToProceedPayloads = 1
	proceededPayloads     = 0
)

func transformPayload(url string, payloadSet [][]string) (PayloadNode, error) {
	matchedPositions := rePayloadPosition.FindAllStringSubmatchIndex(url, -1)
	matchedPatterns := rePayloadPosition.FindAllString(url, -1)

	positionsAmount := len(matchedPositions)
	payloadsAmount := len(payloadSet)

	// Validate payloads and position amount
	if positionsAmount != payloadsAmount {
		p := PayloadNode{}
		if positionsAmount < payloadsAmount {
			return p, errLessPositionsThanPayloads
		}
		if payloadsAmount < positionsAmount {
			return p, errLessPayloadsThanPositions
		}

		return p, errPayloadErr
	}

	var entryNode, previousNode *PayloadNode
	for number, payloads := range payloadSet {
		NeedToProceedPayloads *= len(payloads)
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

	return *entryNode, nil
}

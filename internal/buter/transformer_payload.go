package buter

import (
	"log"
	"os"
	"strings"

	"github.com/edpryk/buter/internal/helpers/prepare"
	"github.com/edpryk/buter/lib/convert"
)

type PayloadNode struct {
	Points            [2]int
	Number            int
	NextNode          *PayloadNode
	PayloadList       []string
	PreviousNode      *PayloadNode
	WorkingPayload    string
	CurrentPayloadIdx int
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

func transformHttpInputToString[V comparable](url string, headers map[string]V, body string) string {
	rawHeaders, err := convert.ToString(headers)
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}

	// TODO: payload positions
	attackValues := []string{url, rawHeaders, body}
	return strings.Join(attackValues, prepare.AttackValueSeparator)
}

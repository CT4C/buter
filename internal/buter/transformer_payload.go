package buter

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/edpryk/buter/lib/convert"
)

type PayloadNode struct {
	PayloadSpan       [2]int
	Number            int
	NextNode          *PayloadNode
	PayloadList       []string
	PreviousNode      *PayloadNode
	WorkingPayload    string
	CurrentPayloadIdx int
}

type HttpRequestProps struct {
	Url     string
	Headers map[string]string
	Body    string
}

var (
	httpRequestPropSeparator = "^"
	rePartsOfAttackValue     = regexp.MustCompile(fmt.Sprintf("([^%s]+)", httpRequestPropSeparator))
)

/*
Transform [][]string to Linked List
*/
func transformPayloadPayloadListToLinked(text string, payloadSet [][]string) (totalPayloads int, entryNode *PayloadNode, err error) {
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
		newNode.PayloadSpan = [2]int{matchedPositions[number][0], matchedPositions[number][1]}
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

func transformHttpRequestPropsToString[V comparable](url string, headers map[string]V, body string) string {
	rawHeaders, err := convert.ToString(headers)
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}

	// TODO: payload positions
	attackValues := []string{url, rawHeaders, body}
	return strings.Join(attackValues, httpRequestPropSeparator)
}

func TransformAttackValueToHttpRequestProps(value string) (HttpRequestProps, error) {
	attackValue := HttpRequestProps{
		Headers: make(map[string]string),
		Body:    "",
		Url:     "",
	}

	matchedValues := rePartsOfAttackValue.FindAllString(value, -1)
	if matchedValues == nil {
		return attackValue, errors.New("Invalid attack value " + value)
	}

	if len(matchedValues) > 0 {
		attackValue.Url = matchedValues[0]
	}

	if len(matchedValues) > 1 {
		rawHeaders := matchedValues[1]

		if err := json.Unmarshal([]byte(rawHeaders), &attackValue.Headers); err != nil && len(rawHeaders) > 2 {
			return attackValue, err
		}
	}

	if len(matchedValues) > 2 {
		attackValue.Body = matchedValues[2]
	}

	return attackValue, nil
}

package main

import (
	"fmt"
	"regexp"
)

var (
	rePayloadPosition = regexp.MustCompile("(%[^%]+%)")
)

type PayloadNode struct {
	Points            [2]int
	WorkingValue      string
	Number            int
	PayloadList       []string
	CurrentPayloadIdx int
	NextNode          *PayloadNode
	PreviousNode      *PayloadNode
}

func main() {
	text := "?param1=%x%&param2=%x%&param3=%x%"
	payload1 := []string{"1", "2", "3"}
	payload2 := []string{"a", "b", "c"}
	payload3 := []string{"L", "M", "N"}

	payloadsSet := [][]string{payload1, payload2, payload3}

	matchedPositions := rePayloadPosition.FindAllStringSubmatchIndex(text, len(payloadsSet))

	// Validate match length == payloads length

	var entryNode, previousNode *PayloadNode
	for number, payload := range payloadsSet {
		newNode := &PayloadNode{
			Number:            number,
			PayloadList:       payload,
			PreviousNode:      previousNode,
			CurrentPayloadIdx: 0,
		}
		newNode.Points = [2]int{matchedPositions[number][0], matchedPositions[number][1]}
		newNode.WorkingValue = text[newNode.Points[0]:newNode.Points[1]]

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

	// Traverse
	node := entryNode
	updatedText := text
	for node != nil {
		currentPayload := node.PayloadList[node.CurrentPayloadIdx]

		if node.NextNode == nil {
			// Proceed last payload in Set
			for _, payload := range node.PayloadList {
				updatedText = updatedText[:node.Points[0]] + payload + updatedText[node.Points[1]:]
				node.CurrentPayloadIdx += 1
				node.Points[1] = node.Points[0] + len(payload)

				fmt.Println(updatedText)
			}

			/*
				Increment previous paylaod index
				set next node to previous one
			*/
			break
		} else {
			// Proceed top level payload
			updatedText = updatedText[:node.Points[0]] + currentPayload + updatedText[node.Points[1]:]

			node.Points[1] = node.Points[0] + len(currentPayload)

			positions := rePayloadPosition.FindAllStringSubmatchIndex(updatedText, 1)
			if len(positions) > 0 {
				node.NextNode.Points[0] = positions[0][0]
				node.NextNode.Points[1] = positions[0][1]
			}
			node.WorkingValue = currentPayload

			// node.CurrentPayloadIdx += 1
		}

		/*
			IF current payload index == payload list length (IS END)
			1. set next node to previous one
			2. reset current payload index
		*/
		node = node.NextNode
	}
}

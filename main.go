package main

import (
	"fmt"
	"regexp"
)

var (
	rePayloadPosition = regexp.MustCompile("(%[^%]+%)")
)

var (
	NeedToProceedPayloads = 1
	proceededPayloads     = 0
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
	text := "?param1=%x%&param2=%y%&param3=%z%"
	payload1 := []string{"1", "2"}
	payload2 := []string{"a", "b"}
	payload3 := []string{"L", "M"}

	payloadsSet := [][]string{payload1, payload2, payload3}

	matchedPositions := rePayloadPosition.FindAllStringSubmatchIndex(text, len(payloadsSet))

	// Validate match length == payloads length

	// Handling text substitute like Cluster Bobm in Burp Suite
	var entryNode, previousNode *PayloadNode
	for number, payloads := range payloadsSet {
		NeedToProceedPayloads *= len(payloads)

		newNode := &PayloadNode{
			Number:            number,
			PayloadList:       payloads,
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
	var node *PayloadNode = entryNode
	updatedText := text
	for node != nil && !(proceededPayloads == NeedToProceedPayloads) {

		if node.NextNode == nil {
			// Proceed last level payload in Set
			for _, payload := range node.PayloadList {
				updatedText = updatedText[:node.Points[0]] + payload + updatedText[node.Points[1]:]
				node.CurrentPayloadIdx += 1
				node.Points[1] = node.Points[0] + len(payload)

				/*
					Send to channel
					Increment proceeded payload
				*/
				fmt.Println(updatedText)
				proceededPayloads += 1
			}

			node.CurrentPayloadIdx = 0
			node = node.PreviousNode
			node.CurrentPayloadIdx += 1
			/*
				Increment previous paylaod index
				set next node to previous one
			*/
			// break
		} else {
			/*
				IF current payload index == payload list length (IS END)
				1. set next node to previous one
				2. reset current payload index
				3. Incremetn Previous payload index ?
			*/

			isEndOfCurrentPayloadProcessing := node.CurrentPayloadIdx == len(node.PayloadList)
			// isLastNode := node.PreviousNode == nil

			// Return back to previous level
			if isEndOfCurrentPayloadProcessing {
				node.CurrentPayloadIdx = 0 //
				node = node.PreviousNode
				// It's time to nexet payload in list
				node.CurrentPayloadIdx += 1
			}

			// isPreviousEmpty := node.PreviousNode == nil
			// isFullSetProceeded := isPreviousEmpty && isEndOfCurrentPayloadProcessing
			// IS FUll Payload Set has proceeded
			// if isFullSetProceeded {
			// 	fmt.Println("COMPLETED")
			// 	break
			// }

			// Proceed top level payload
			currentPayload := node.PayloadList[node.CurrentPayloadIdx]

			updatedText = updatedText[:node.Points[0]] + currentPayload + updatedText[node.Points[1]:]

			node.Points[1] = node.Points[0] + len(currentPayload)

			positions := rePayloadPosition.FindAllStringSubmatchIndex(updatedText, 1)
			if len(positions) > 0 {
				node.NextNode.Points[0] = positions[0][0]
				node.NextNode.Points[1] = positions[0][1]
			}
			node.WorkingValue = currentPayload

			// node.CurrentPayloadIdx += 1
			node = node.NextNode
		}
	}
}

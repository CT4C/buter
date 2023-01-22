package buter

import (
	"context"
)

type Cluster struct {
	AttackValue       string
	Ctx               context.Context
	EntryNode         *PayloadNode
	TotalPayloads     int
	proceededPayloads int
	errChanel         chan error
}

func (c *Cluster) ProduceUrls(urlConsumer chan string) chan error {
	/*
		TODO: Miss when one payload grater the another one
	*/
	defer close(urlConsumer)
	defer close(c.errChanel)

	if c.TotalPayloads <= 0 {
		c.errChanel <- errInvalidTotalPayload
		return c.errChanel
	}

	var (
		updatedText = c.AttackValue
	)

	for c.EntryNode != nil && !(c.proceededPayloads == c.TotalPayloads) {

		if c.EntryNode.NextNode == nil {
			// Proceed last level payload in Set
			for _, payload := range c.EntryNode.PayloadList {
				updatedText = updatedText[:c.EntryNode.Points[0]] + payload + updatedText[c.EntryNode.Points[1]:]
				c.EntryNode.CurrentPayloadIdx += 1
				c.EntryNode.Points[1] = c.EntryNode.Points[0] + len(payload)
				/*
					- Send to channel
					- Increment proceeded payload
				*/
				urlConsumer <- updatedText
				c.proceededPayloads += 1
			}

			/*
				- Increment previous paylaod index
				- set next c.EntryNode to previous one
			*/
			c.EntryNode.CurrentPayloadIdx = 0
			c.EntryNode = c.EntryNode.PreviousNode
			c.EntryNode.CurrentPayloadIdx += 1
		} else {
			// ### TOP level paylaod processing ###

			/*
				IF current payload index == payload list length (IS END)
				1. set next c.EntryNode to previous one
				2. reset current payload index
				3. Incremetn Previous payload index ?
			*/

			isEndOfCurrentPayloadProcessing := c.EntryNode.CurrentPayloadIdx == len(c.EntryNode.PayloadList)

			if isEndOfCurrentPayloadProcessing {
				/*
					Reset current payload index before beign go to
					the next in set
				*/
				c.EntryNode.CurrentPayloadIdx = 0
				/*
					Set previous payload to the current one
				*/
				c.EntryNode = c.EntryNode.PreviousNode
				/*
					Increment working payload index
				*/
				c.EntryNode.CurrentPayloadIdx += 1
			}

			nextPayload := c.EntryNode.PayloadList[c.EntryNode.CurrentPayloadIdx]
			currentPayload := c.EntryNode.WorkingPayload

			/*
				Points correction - when one payload
				length greater then anoher one
			*/
			if len(c.EntryNode.WorkingPayload) < len(nextPayload) {
				nextPointStart := c.EntryNode.NextNode.Points[0]
				payloadShift := (len(nextPayload) - len(currentPayload))
				c.EntryNode.NextNode.Points[0] = nextPointStart + payloadShift
				c.EntryNode.NextNode.Points[1] = nextPointStart + len(nextPayload)
			}

			c.EntryNode.WorkingPayload = nextPayload

			updatedText = updatedText[:c.EntryNode.Points[0]] + c.EntryNode.WorkingPayload + updatedText[c.EntryNode.Points[1]:]

			c.EntryNode.Points[1] = c.EntryNode.Points[0] + len(c.EntryNode.WorkingPayload)

			positions := rePayloadPosition.FindAllStringSubmatchIndex(updatedText, -1)
			if len(positions) > 0 {
				c.EntryNode.NextNode.Points[0] = positions[0][0]
				c.EntryNode.NextNode.Points[1] = positions[0][1]
			}

			currentNodeCopy := *c.EntryNode
			c.EntryNode = c.EntryNode.NextNode
			c.EntryNode.PreviousNode = &currentNodeCopy
		}
	}

	// c.errChanel <- nil
	return c.errChanel
}

func (c Cluster) Proceeded() int {
	return c.proceededPayloads
}

func NewCluster(ctx context.Context, attackValue string, entryNode *PayloadNode, totalPyaloads int) *Cluster {
	return &Cluster{
		Ctx:           ctx,
		AttackValue:   attackValue,
		EntryNode:     entryNode,
		TotalPayloads: totalPyaloads,
		/*
			Using unbuff channles in synchronous code causes deadlock,
			because runtime is blocked in the place where you send the
			value to a chnnel, til someone else read the value, but in
			synchronous code, no one else can't read in the same time
		*/
		errChanel: make(chan error, 1),
	}
}

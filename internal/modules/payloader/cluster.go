package payloader

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
	PositionsAmount   int
}

func (c *Cluster) ProduceUrls(urlConsumer chan CraftedPayload) chan error {
	/*
		TODO: Miss when one payloader grater the another one
	*/
	defer close(urlConsumer)
	defer close(c.errChanel)

	if c.TotalPayloads <= 0 {
		c.errChanel <- errInvalidTotalPayload
		return c.errChanel
	}

	var (
		updatedText        = c.AttackValue
		workingPayloadsSet = make([]string, c.PositionsAmount)
	)

	for c.EntryNode != nil && !(c.proceededPayloads == c.TotalPayloads) {

		if c.EntryNode.NextNode == nil {
			// ### Last Level payloader Processing ###
			for _, payload := range c.EntryNode.PayloadList {
				workingPayloadsSet[c.EntryNode.Number] = payload

				updatedText = updatedText[:c.EntryNode.Points[0]] + payload + updatedText[c.EntryNode.Points[1]:]

				c.EntryNode.CurrentPayloadIdx += 1
				c.EntryNode.Points[1] = c.EntryNode.Points[0] + len(payload)
				/*
					1. Send to channel
					2. Increment proceeded payloader
				*/
				urlConsumer <- CraftedPayload{
					Value:    updatedText,
					Payloads: workingPayloadsSet,
				}
				c.proceededPayloads += 1
			}

			/*
				1. Increment previous paylaod index
				2. set next c.EntryNode to previous one
			*/
			c.EntryNode.CurrentPayloadIdx = 0
			c.EntryNode = c.EntryNode.PreviousNode
			c.EntryNode.CurrentPayloadIdx += 1
			workingPayloadsSet = make([]string, c.PositionsAmount)
		} else {
			// ### TOP level paylaod processing ###

			/*
				IF current payloader index == payloader list length (IS END)
				1. set next c.EntryNode to previous one
				2. reset current payloader index
				3. Incremetn Previous payloader index ?
			*/

			isEndOfCurrentPayloadProcessing := c.EntryNode.CurrentPayloadIdx == len(c.EntryNode.PayloadList)

			if isEndOfCurrentPayloadProcessing {
				/*
					Reset current payloader index before beign go to
					the next in set
				*/
				c.EntryNode.CurrentPayloadIdx = 0
				/*
					Set previous payloader to the current one
				*/
				c.EntryNode = c.EntryNode.PreviousNode
				/*
					Increment working payloader index
				*/
				c.EntryNode.CurrentPayloadIdx += 1
			}

			currentPayload := c.EntryNode.WorkingPayload
			nextPayload := c.EntryNode.PayloadList[c.EntryNode.CurrentPayloadIdx]

			/*
				Points correction - when one payloader
				length greater then another one
			*/
			if len(c.EntryNode.WorkingPayload) < len(nextPayload) {
				nextPointStart := c.EntryNode.NextNode.Points[0]
				payloadShift := (len(nextPayload) - len(currentPayload))
				c.EntryNode.NextNode.Points[0] = nextPointStart + payloadShift
				c.EntryNode.NextNode.Points[1] = nextPointStart + len(nextPayload)
			}

			c.EntryNode.WorkingPayload = nextPayload

			updatedText = updatedText[:c.EntryNode.Points[0]] + c.EntryNode.WorkingPayload + updatedText[c.EntryNode.Points[1]:]

			/*
				Current Points correction
			*/
			c.EntryNode.Points[1] = c.EntryNode.Points[0] + len(c.EntryNode.WorkingPayload)
			/*
				Defined payloader correction, check if it exists yet
				found and update points, if it doesn't exists that's
				mean that all defined pattern already in substitute
				process wihtin payloads from lists
			*/
			positions := rePayloadPosition.FindAllStringSubmatchIndex(updatedText, -1)
			if len(positions) > 0 {
				c.EntryNode.NextNode.Points[0] = positions[0][0]
				c.EntryNode.NextNode.Points[1] = positions[0][1]
			}

			/*
				Added working payload to working payload set
			*/
			workingPayloadsSet[c.EntryNode.Number] = c.EntryNode.WorkingPayload
			/*
				Copy created becuase it was lose fo pointer to the
				EntryNode
			*/
			currentNodeCopy := *c.EntryNode
			c.EntryNode = c.EntryNode.NextNode
			c.EntryNode.PreviousNode = &currentNodeCopy
		}
	}

	return c.errChanel
}

func (c Cluster) Proceeded() int {
	return c.proceededPayloads
}

func NewCluster(ctx context.Context, attackValue string, entryNode *PayloadNode, totalPyaloads int, positionsAmount int) *Cluster {
	return &Cluster{
		Ctx:             ctx,
		AttackValue:     attackValue,
		EntryNode:       entryNode,
		TotalPayloads:   totalPyaloads,
		PositionsAmount: positionsAmount,
		/*
			Using unbuff channles in synchronous code causes deadlock,
			because runtime is blocked in the place where you send the
			value to a chnnel, til someone else read the value, but in
			synchronous code, no one else can't read in the same time
		*/
		errChanel: make(chan error, 1),
	}
}

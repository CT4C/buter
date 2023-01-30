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
	workingPayloadSet []string
	endSig            chan int
}

func (c *Cluster) ProducePayload(payloadConsumer chan CraftedPayload) chan int {
	defer func() {
		c.endSig <- 0
	}()
	/*
		TODO: Miss when one payloader grater the another one
	*/
	defer close(c.errChanel)

	if c.TotalPayloads <= 0 {
		c.errChanel <- errInvalidTotalPayload
		return c.endSig
	}

	var (
		updatedAttackValue = c.AttackValue
	)

	for c.EntryNode != nil && !(c.proceededPayloads == c.TotalPayloads) {

		if c.EntryNode.NextNode == nil {
			producedPayloads := processPayloads(updatedAttackValue, c.EntryNode, c.workingPayloadSet, payloadConsumer)

			/*
				1. Increment previous paylaod index
				2. set next c.EntryNode to previous one
			*/
			/*
				TODO: Need to add functionality to the linked list
				like as BackPrevioudNode/StepBack/Return/Forwared/Next
				for movind back/forward btw nodes
			*/
			c.proceededPayloads += producedPayloads
			c.EntryNode.CurrentPayloadIdx = 0
			c.EntryNode = c.EntryNode.PreviousNode
			c.EntryNode.CurrentPayloadIdx += 1
			c.workingPayloadSet = make([]string, c.PositionsAmount)
		} else {
			// ### TOP level paylaod processing ###
			c.workingPayloadSet[c.EntryNode.Number] = c.EntryNode.PayloadList[c.EntryNode.CurrentPayloadIdx]

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

			updatedAttackValue = updateValue(updatedAttackValue, c.EntryNode.WorkingPayload, c.EntryNode.Points)

			/*
				Current Points correction

				Need to add method to the node to proceed this case
			*/
			c.EntryNode.Points[1] = c.EntryNode.Points[0] + len(c.EntryNode.WorkingPayload)
			/*
				Defined payloader correction, check if it exists yet
				found and update points, if it doesn't exists that's
				mean that all defined pattern already in substitute
				process wihtin payloads from lists
			*/
			positions := rePayloadPosition.FindAllStringSubmatchIndex(updatedAttackValue, -1)
			if len(positions) > 0 {
				c.EntryNode.NextNode.Points[0] = positions[0][0]
				c.EntryNode.NextNode.Points[1] = positions[0][1]
			}

			/*
				Added working payload to working payload set
			*/
			/*
				Copy created becuase it was lose fo pointer to the
				EntryNode
			*/
			currentNodeCopy := *c.EntryNode
			c.EntryNode = c.EntryNode.NextNode
			c.EntryNode.PreviousNode = &currentNodeCopy
		}
	}

	return c.endSig
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
		errChanel:         make(chan error, 1),
		workingPayloadSet: make([]string, positionsAmount),
		endSig:            make(chan int, 1),
	}
}

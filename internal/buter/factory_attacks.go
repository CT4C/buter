package buter

import (
	"context"
	"log"
	"os"

	"github.com/edpryk/buter/cli"
)

type attackConfig struct {
	Ctx         context.Context
	Consumer    PayloadConsumer
	PayloadNode *PayloadNode

	AttackType            string
	TargetTextRaw         string
	ItemProducePlan       int
	TotalPayloadPositions int
}

type attackFactory struct {
	attackConfig

	isClosed          chan int
	errorsReport      chan error
	producedItems     int
	workingPayloadSet []string
}

func newAttackFactory(config attackConfig) *attackFactory {
	factory := &attackFactory{
		attackConfig: config,

		isClosed:          make(chan int, 1),
		errorsReport:      make(chan error, 1),
		workingPayloadSet: make([]string, config.TotalPayloadPositions),
	}
	return factory
}

func (factory attackFactory) Launch() chan int {
	go func() {
		defer func() {
			factory.Consumer.Close()
			factory.isClosed <- 0
		}()

		switch factory.AttackType {
		case cli.ClusterAttack:
			factory.clusterWorker()
			return
		case cli.SniperAttack:
			factory.sniperWorker()
			return
		case cli.DOSAttack:
			factory.dosWorker()
			return
		case cli.PitchForkAttack:
			factory.pitchforkWorker()
			return
		default:
			log.Println(errAttackNotSupported)
			os.Exit(1)
		}
	}()

	return factory.isClosed
}

func (factory *attackFactory) onPayloadUpdated(updatedTargetString string, payloadInserted string, payloadNumber int) {
	factory.workingPayloadSet[payloadNumber] = payloadInserted

	workingPayloadSetCopy := make([]string, len(factory.workingPayloadSet))
	copy(workingPayloadSetCopy, factory.workingPayloadSet)

	factory.Consumer.Consume(updatedTargetString, workingPayloadSetCopy, nil)
}

func (factory attackFactory) dosWorker() chan int {
	for i := 0; i < factory.ItemProducePlan; i++ {
		factory.Consumer.Consume(factory.TargetTextRaw, []string{}, nil)
		// TODO: ctx
	}

	return factory.isClosed
}

func (factory *attackFactory) sniperWorker() {
	factory.producedItems += buildPayloadList(
		factory.TargetTextRaw,
		factory.PayloadNode,
		factory.onPayloadUpdated,
	)
}

func (factory *attackFactory) pitchforkWorker() {
	updatedValue := factory.TargetTextRaw
	for factory.producedItems < factory.ItemProducePlan {
		// walk through each node once
		for i := 0; i < factory.TotalPayloadPositions; i++ {
			node := factory.PayloadNode
			if node.CurrentPayloadIdx == len(node.PayloadList) {
				/*
					Save point for situation when one list shorted
					than another one
				*/
				node.CurrentPayloadIdx = 0
			}
			updatedValue = insertPayload(updatedValue, node.PayloadList[node.CurrentPayloadIdx], node.PayloadSpan)
			node.WorkingPayload = node.PayloadList[node.CurrentPayloadIdx]
			node.CurrentPayloadIdx++

			// Update payload working set for advance is in reporter
			factory.workingPayloadSet[node.Number] = node.WorkingPayload
			// Switch node back
			if node.NextNode == nil {
				factory.onPayloadUpdated(updatedValue, node.WorkingPayload, node.Number)
				factory.PayloadNode = node.PreviousNode
			} else {
				/*
					Update payload positions
				*/
				shift := node.NextNode.PayloadSpan[0] - node.PayloadSpan[1]
				if shift < 0 {
					shift = shift * -1
				}
				/*
					Update Current Node Span
				*/
				node.PayloadSpan[1] = node.PayloadSpan[0] + len(node.WorkingPayload)
				/*
					Update Next Node Span
				*/
				nextNodeSpanStart := node.PayloadSpan[1] + shift
				nextNodeSpanEnd := nextNodeSpanStart + len(node.NextNode.WorkingPayload)
				factory.PayloadNode = node.NextNode
				factory.PayloadNode.PayloadSpan[0] = nextNodeSpanStart
				factory.PayloadNode.PayloadSpan[1] = nextNodeSpanEnd
			}

		}
		factory.producedItems += 1
	}
}

func (factory *attackFactory) clusterWorker() {
	/*
		TODO: Miss when one payload grater the another one
	*/

	for factory.PayloadNode != nil && !(factory.producedItems == factory.ItemProducePlan) {

		if factory.PayloadNode.NextNode == nil {
			producedPayloads := buildPayloadList(
				factory.TargetTextRaw,
				factory.PayloadNode,
				factory.onPayloadUpdated,
			)
			/*
				1. Increment previous payload index
				2. set next factory.PayloadNode to previous one
			*/
			/*
				TODO: Need to add functionality to the linked list
				like as BackPreviousNode/StepBack/Return/Forwarded/Next
				for moving back/forward btw nodes
			*/
			factory.producedItems += producedPayloads
			factory.PayloadNode.CurrentPayloadIdx = 0
			// factory.PayloadNode = factory.PayloadNode.PreviousNode
			factory.PayloadNode.Prev().CurrentPayloadIdx += 1
			// factory.PayloadNode.CurrentPayloadIdx += 1
			factory.workingPayloadSet = make([]string, factory.TotalPayloadPositions)
		} else {
			// ### TOP level payload processing ###
			factory.workingPayloadSet[factory.PayloadNode.Number] = factory.PayloadNode.PayloadList[factory.PayloadNode.CurrentPayloadIdx]

			/*
				IF current payload index == payload list length (IS END)
				1. set next factory.PayloadNode to previous one
				2. reset current payload index
				3. Increment Previous payload index ?
			*/

			isEndOfCurrentPayloadProcessing := factory.PayloadNode.CurrentPayloadIdx == len(factory.PayloadNode.PayloadList)

			if isEndOfCurrentPayloadProcessing {
				/*
					Reset current payload index before being go to
					the next in set
				*/
				factory.PayloadNode.CurrentPayloadIdx = 0
				/*
					Set previous payload to the current one
				*/
				factory.PayloadNode = factory.PayloadNode.PreviousNode
				/*
					Increment working payload index
				*/
				factory.PayloadNode.CurrentPayloadIdx += 1
			}

			currentPayload := factory.PayloadNode.WorkingPayload
			nextPayload := factory.PayloadNode.PayloadList[factory.PayloadNode.CurrentPayloadIdx]

			/*
				PayloadSpan correction - when one payload
				length greater then another one
			*/
			if len(factory.PayloadNode.WorkingPayload) < len(nextPayload) {
				nextPointStart := factory.PayloadNode.NextNode.PayloadSpan[0]
				payloadShift := (len(nextPayload) - len(currentPayload))
				factory.PayloadNode.NextNode.PayloadSpan[0] = nextPointStart + payloadShift
				factory.PayloadNode.NextNode.PayloadSpan[1] = nextPointStart + len(nextPayload)
			}

			factory.PayloadNode.WorkingPayload = nextPayload

			factory.TargetTextRaw = insertPayload(
				factory.TargetTextRaw,
				factory.PayloadNode.WorkingPayload,
				factory.PayloadNode.PayloadSpan,
			)

			/*
				Current PayloadSpan correction

				Need to add method to the node to proceed this case
			*/
			factory.PayloadNode.PayloadSpan[1] = factory.PayloadNode.PayloadSpan[0] + len(factory.PayloadNode.WorkingPayload)
			/*
				Defined payload correction, check if it exists yet
				found and update points, if it doesn't exists that'factory
				mean that all defined pattern already in substitute
				process within payloads from lists
			*/
			positions := rePayloadPosition.FindAllStringSubmatchIndex(factory.TargetTextRaw, -1)
			if len(positions) > 0 {
				factory.PayloadNode.NextNode.PayloadSpan[0] = positions[0][0]
				factory.PayloadNode.NextNode.PayloadSpan[1] = positions[0][1]
			}

			/*
				Added working payload to working payload set
			*/
			/*
				Copy created because it was lose of ptr to the
				PayloadNode
			*/
			currentNodeCopy := *factory.PayloadNode
			factory.PayloadNode = factory.PayloadNode.NextNode
			factory.PayloadNode.PreviousNode = &currentNodeCopy
		}
	}
}

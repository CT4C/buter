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
	RawPayload            string
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
		default:
			log.Println(errAttackNotSupported)
			os.Exit(1)
		}
	}()

	return factory.isClosed
}

func (factory attackFactory) onPayloadUpdated(updatedTargetString string, payloadInserted string, payloadNumber int) {
	factory.workingPayloadSet[payloadNumber] = payloadInserted
	/*
		as factory.workingPayloadSet is a pointer need to copy this before pass
		to the consumer to prevent inconsistent state
	*/
	workingPayloadSetCopy := make([]string, len(factory.workingPayloadSet))
	copy(workingPayloadSetCopy, factory.workingPayloadSet)

	factory.Consumer.Consume(updatedTargetString, workingPayloadSetCopy, nil)
}

func (factory attackFactory) dosWorker() chan int {
	for i := 0; i < factory.ItemProducePlan; i++ {
		factory.Consumer.Consume(factory.RawPayload, []string{}, nil)
		// TODO: ctx
	}

	return factory.isClosed
}

func (factory *attackFactory) sniperWorker() {
	factory.producedItems += buildPayloadList(
		factory.RawPayload,
		factory.PayloadNode,
		factory.onPayloadUpdated,
	)
}

func (factory *attackFactory) clusterWorker() {
	/*
		TODO: Miss when one payload grater the another one
	*/

	for factory.PayloadNode != nil && !(factory.producedItems == factory.ItemProducePlan) {

		if factory.PayloadNode.NextNode == nil {
			producedPayloads := buildPayloadList(
				factory.RawPayload,
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
			factory.PayloadNode = factory.PayloadNode.PreviousNode
			factory.PayloadNode.CurrentPayloadIdx += 1
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
				Points correction - when one payload
				length greater then another one
			*/
			if len(factory.PayloadNode.WorkingPayload) < len(nextPayload) {
				nextPointStart := factory.PayloadNode.NextNode.Points[0]
				payloadShift := (len(nextPayload) - len(currentPayload))
				factory.PayloadNode.NextNode.Points[0] = nextPointStart + payloadShift
				factory.PayloadNode.NextNode.Points[1] = nextPointStart + len(nextPayload)
			}

			factory.PayloadNode.WorkingPayload = nextPayload

			factory.RawPayload = insertPayload(
				factory.RawPayload,
				factory.PayloadNode.WorkingPayload,
				factory.PayloadNode.Points,
			)

			/*
				Current Points correction

				Need to add method to the node to proceed this case
			*/
			factory.PayloadNode.Points[1] = factory.PayloadNode.Points[0] + len(factory.PayloadNode.WorkingPayload)
			/*
				Defined payload correction, check if it exists yet
				found and update points, if it doesn't exists that'factory
				mean that all defined pattern already in substitute
				process within payloads from lists
			*/
			positions := rePayloadPosition.FindAllStringSubmatchIndex(factory.RawPayload, -1)
			if len(positions) > 0 {
				factory.PayloadNode.NextNode.Points[0] = positions[0][0]
				factory.PayloadNode.NextNode.Points[1] = positions[0][1]
			}

			/*
				Added working payload to working payload set
			*/
			/*
				Copy created because it was lose fo pointer to the
				PayloadNode
			*/
			currentNodeCopy := *factory.PayloadNode
			factory.PayloadNode = factory.PayloadNode.NextNode
			factory.PayloadNode.PreviousNode = &currentNodeCopy
		}
	}
}

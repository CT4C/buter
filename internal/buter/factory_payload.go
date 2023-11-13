package buter

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/edpryk/buter/cli"
)

type Attacker interface {
	Launch() chan int
}

type PayloadConsumer interface {
	Consume(updatedAttackValue string, payloads []string, err error)
	Close()
}

type CraftedPayload struct {
	Url      string
	Body     string
	Headers  map[string]string
	Payloads []string
}

type Config struct {
	HttpRequestProps string
	Url              string
	Ctx              context.Context
	Body             string
	Headers          map[string]string
	PayloadSet       [][]string
	AttackType       string
	QueueLength      int
	MaxRequests      int
	TotalPayloads    int
}

type PayloadFactory struct {
	Config

	startTime time.Time
}

func (factory *PayloadFactory) Launch(consumer PayloadConsumer) {
	go func() {
		/*
			TODO: Add name/number to payloads
		*/
		attackValue := transformHttpRequestPropsToString(factory.Url, factory.Headers, factory.Body)

		_, entryPayloadNode, err := convertPayloadListToLinked(attackValue, factory.PayloadSet)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		totalPayloads := 1

		if factory.AttackType == cli.ClusterAttack {
			for _, list := range factory.PayloadSet {
				totalPayloads *= len(list)
			}
		}

		if factory.AttackType == cli.SniperAttack {
			totalPayloads = len(factory.PayloadSet[0])
		}

		if factory.AttackType == cli.PitchForkAttack {
			for _, list := range factory.PayloadSet {
				l := len(list)

				totalPayloads = l
				if l > totalPayloads {
					totalPayloads = l
				}
			}
		}

		if factory.AttackType == cli.DOSAttack {
			totalPayloads = factory.MaxRequests
		}

		attackFactory := newAttackFactory(attackConfig{
			Ctx:                   factory.Ctx,
			Consumer:              consumer,
			AttackType:            factory.AttackType,
			TargetTextRaw:         attackValue,
			PayloadNode:           entryPayloadNode,
			ItemProducePlan:       totalPayloads,
			TotalPayloadPositions: len(factory.PayloadSet),
		})

		select {
		case <-attackFactory.Launch():
			return
		case <-factory.Ctx.Done():
			log.Println("PayloadFactory Canceled")
			return
		}
	}()
}

func NewFactory(config Config) *PayloadFactory {
	if config.QueueLength == 0 {
		config.QueueLength = 1
	}

	factory := &PayloadFactory{
		Config:    config,
		startTime: time.Now(),
	}

	return factory
}

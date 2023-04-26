package buter

import (
	"context"
	"log"
	"os"
	"time"
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

		totalPayloads, entryPayloadNode, err := convertPayloadListToLinked(attackValue, factory.PayloadSet)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		attackFactory := newAttackFactory(attackConfig{
			Ctx:                   factory.Ctx,
			Consumer:              consumer,
			AttackType:            factory.AttackType,
			RawPayload:            attackValue,
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

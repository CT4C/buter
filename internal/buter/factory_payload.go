package buter

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/edpryk/buter/cli"
)

type Attacker interface {
	Proceeded() int
	ProducePayload(payloadConsumer PayloadConsumer) chan int
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
	AttackValue   string
	Url           string
	Ctx           context.Context
	Body          string
	Headers       map[string]string
	PayloadSet    [][]string
	AttackType    string
	QueueLength   int
	MaxRequests   int
	TotalPayloads int
}

type ProcessStatus struct {
	Message string
	Err     bool
}

type PayloadFactory struct {
	Config

	attackValue string
	payloadNode *PayloadNode
	attacker    Attacker
	startTime   time.Time
}

func (b *PayloadFactory) Launch(consumer PayloadConsumer) {
	go func() {
		/*
			TODO: Add name/number to payloads
		*/
		b.attackValue = transformHttpInputToString(b.Url, b.Headers, b.Body)

		totalPayloads, entryNode, err := transformPayload(b.attackValue, b.PayloadSet)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		b.TotalPayloads = totalPayloads
		b.payloadNode = entryNode

		if err := b.chooseAttackWorker(b.AttackType); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// attacker close the channel

		select {
		case <-b.attacker.ProducePayload(consumer):
			return
		case <-b.Ctx.Done():
			fmt.Println("PayloadFactory Canceled")
			return
		}
	}()
}

func (b *PayloadFactory) chooseAttackWorker(attackType string) error {
	switch attackType {
	case cli.ClusterAttack:
		b.attacker = NewCluster(b.Ctx, b.attackValue, b.payloadNode, b.TotalPayloads, len(b.PayloadSet))
		return nil
	case cli.SniperAttack:
		b.attacker = NewSniper(b.Ctx, b.attackValue, b.payloadNode)
		return nil
	case cli.DOSAttack:
		b.attacker = NewDos(b.Ctx, b.attackValue, b.MaxRequests)
		return nil
	default:
		return errAttackNotSupported
	}
}

func NewFactory(config Config) *PayloadFactory {
	if config.QueueLength == 0 {
		config.QueueLength = 1
	}

	b := &PayloadFactory{
		Config:    config,
		startTime: time.Now(),
	}

	return b
}

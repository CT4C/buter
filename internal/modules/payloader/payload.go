package payloader

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/edpryk/buter/cli"
)

type Attacker interface {
	ProducePayload(urlConsumer chan CraftedPayload) chan error
	Proceeded() int
}

type CraftedPayload struct {
	Url      string
	Payloads []string
	Headers  map[string]string
}

type Config struct {
	Ctx           context.Context
	PayloadSet    [][]string
	AttackType    string
	AttackValue   string
	TotalPayloads int
	QueueLength   int
}

type ProcessStatus struct {
	Message string
	Err     bool
}

type Buter struct {
	Config

	payloadEntryNode *PayloadNode
	attacker         Attacker
	startTime        time.Time
	payloadProvider  chan CraftedPayload
	errQ             chan error
}

func (b *Buter) PrepareAttack() (payloadPrivider chan CraftedPayload, err chan error) {
	go func() {
		defer close(b.errQ)

		totalPayloads, entryNode, err := transformPayload(b.AttackValue, b.PayloadSet)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		b.TotalPayloads = totalPayloads
		b.payloadEntryNode = entryNode

		if err := b.setAttacker(b.AttackType); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		go func() {
			defer close(b.payloadProvider)

			select {
			case err := <-b.attacker.ProducePayload(b.payloadProvider):
				if err != nil {
					fmt.Println(err)
					b.errQ <- err
				}

				return
			case <-b.Ctx.Done():
				return
			}
		}()
	}()

	return b.payloadProvider, b.errQ
}

func (b *Buter) setAttacker(attackType string) error {
	switch attackType {
	case cli.ClusterAttack:
		b.attacker = NewCluster(b.Ctx, b.AttackValue, b.payloadEntryNode, b.TotalPayloads, len(b.PayloadSet))
		return nil
	default:
		return errAttackNotSupported
	}
}

func New(config Config) *Buter {
	if config.QueueLength == 0 {
		config.QueueLength = 1
	}

	b := &Buter{
		Config:          config,
		startTime:       time.Now(),
		payloadProvider: make(chan CraftedPayload, config.QueueLength),
		errQ:            make(chan error, 1),
	}

	return b
}

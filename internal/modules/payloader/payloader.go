package payloader

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/edpryk/buter/cli"
)

type Attacker interface {
	ProducePayload(payloadConsumer chan CraftedPayload) chan int
	Proceeded() int
}

type CraftedPayload struct {
	Url      string
	Payloads []string
	Headers  map[string]string
	Body     map[string]string
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

	payloadNode     *PayloadNode
	attacker        Attacker
	startTime       time.Time
	payloadProvider chan CraftedPayload
	errQ            chan error
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
		b.payloadNode = entryNode

		if err := b.setAttacker(b.AttackType); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		go func() {
			defer close(b.payloadProvider)

			select {
			case <-b.attacker.ProducePayload(b.payloadProvider):
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
		b.attacker = NewCluster(b.Ctx, b.AttackValue, b.payloadNode, b.TotalPayloads, len(b.PayloadSet))
		return nil
	case cli.SniperAttack:
		b.attacker = NewSniper(b.Ctx, b.AttackValue, b.payloadNode)
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

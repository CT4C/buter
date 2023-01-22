package payloader

import (
	"context"
	"fmt"
	"time"

	"github.com/edpryk/buter/internal/docs"
)

type Attacker interface {
	ProduceUrls(urlConsumer chan CraftedPayload) chan error
	Proceeded() int
}

type CraftedPayload struct {
	Value    string
	Payloads []string
}

type Config struct {
	Ctx             context.Context
	PayloadConsumer chan CraftedPayload
	PayloadSet      [][]string
	AttackType      string
	AttackValue     string
	TotalPayloads   int
	StatusChan      chan ProcessStatus
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
}

func (b *Buter) PrepareAttack() error {
	totalPayloads, entryNode, err := transformPayload(b.AttackValue, b.PayloadSet)
	if err != nil {
		return err
	}
	b.sendStatus(fmt.Sprintf("[+] Prepared %d payloads\n", totalPayloads), false)

	b.TotalPayloads = totalPayloads
	b.payloadEntryNode = entryNode

	if err := b.setAttacker(b.AttackType); err != nil {
		return err
	}

	go func() {
		defer close(b.StatusChan)

		select {
		case err := <-b.attacker.ProduceUrls(b.PayloadConsumer):
			if err == nil {
				b.complete()
			} else {
				b.sendStatus(err.Error(), true)
			}
		case <-b.Ctx.Done():
			b.terminate()
		}
	}()

	return nil
}

func (b Buter) sendStatus(message string, err bool) {
	b.StatusChan <- ProcessStatus{
		Message: message,
		Err:     err,
	}
}

func (b Buter) terminate() {
	b.sendStatus(fmt.Sprintf("[*] Process timeout, proceeded %d payloads\n", b.attacker.Proceeded()), true)
}

func (b Buter) complete() {
	b.sendStatus(fmt.Sprintf("[+] Completed %d paloads in %s\n", b.attacker.Proceeded(), time.Since(b.startTime)), false)
}

func (b *Buter) setAttacker(attackType string) error {
	switch attackType {
	case docs.ClusterAttack:
		b.attacker = NewCluster(b.Ctx, b.AttackValue, b.payloadEntryNode, b.TotalPayloads, len(b.PayloadSet))
		return nil
	default:
		return errAttackNotSupported
	}
}

func New(config Config) *Buter {
	b := &Buter{
		Config:    config,
		startTime: time.Now(),
	}

	return b
}

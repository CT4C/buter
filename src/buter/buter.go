package buter

import (
	"context"
	"fmt"
	"time"

	"github.com/edpryk/buter/src/docs"
)

type Attacker interface {
	ProduceUrls(urlConsumer chan string) chan error
	Proceeded() int
}

type Config struct {
	Ctx           context.Context
	UrlConsumer   chan string
	PayloadSet    [][]string
	AttackType    string
	Url           string
	TotalPayloads int
	StatusChan    chan ProcessStatus
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

func (b *Buter) PrepareAttackUrls() error {
	totalPayloads, entryNode, err := transformPayload(b.Url, b.PayloadSet)
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
		case err := <-b.attacker.ProduceUrls(b.UrlConsumer):
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
		b.attacker = NewCluster(b.Ctx, b.Url, b.payloadEntryNode, b.TotalPayloads)
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

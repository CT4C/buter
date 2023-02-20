package buter

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/edpryk/buter/cli"
	"github.com/edpryk/buter/internal/helpers/prepare"
	"github.com/edpryk/buter/lib/convert"
)

type MapStringer interface {
	Map() map[string]any
	String() string
}

type Attacker interface {
	Proceeded() int
	ProducePayload(payloadConsumer chan CraftedPayload) chan int
}

type CraftedPayload struct {
	Url      string
	Body     map[string]string
	Headers  map[string]string
	Payloads []string
}

type Config struct {
	Url           string
	Ctx           context.Context
	Body          map[string]string
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

type Buter struct {
	Config

	attackValue     string
	payloadNode     *PayloadNode
	attacker        Attacker
	startTime       time.Time
	payloadProvider chan CraftedPayload
	errQ            chan error
}

func (b *Buter) RunPrepareAttack() (payloadProvider chan CraftedPayload, err chan error) {
	go func() {
		defer close(b.errQ)

		/*
			TODO: Add name to payloads
		*/

		rawHeaders, err := convert.ToString(b.Headers)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}
		rawBody, err := convert.ToString(b.Body)
		if err != nil {
			log.Println(err)
			os.Exit(0)
		}

		// Order depends on flags ordering
		attackValues := []string{b.Url, rawHeaders, rawBody}
		attackValueString := strings.Join(attackValues, prepare.AttackValueSeparator)

		totalPayloads, entryNode, err := transformPayload(attackValueString, b.PayloadSet)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		b.TotalPayloads = totalPayloads
		b.payloadNode = entryNode
		b.attackValue = attackValueString

		if err := b.setAttacker(b.AttackType); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer close(b.payloadProvider)

		select {
		case <-b.attacker.ProducePayload(b.payloadProvider):
			return
		case <-b.Ctx.Done():
			fmt.Println("Buter Canceled")
			return
		}
	}()

	return b.payloadProvider, b.errQ
}

func (b *Buter) setAttacker(attackType string) error {
	switch attackType {
	case cli.ClusterAttack:
		b.attacker = NewCluster(b.Ctx, b.attackValue, b.payloadNode, b.TotalPayloads, len(b.PayloadSet))
		return nil
	case cli.SniperAttack:
		b.attacker = NewSniper(b.Ctx, b.attackValue, b.payloadNode)
		return nil
	case cli.DOSAttack:
		b.attacker = NewDos(DosConfig{
			Url:         b.Url,
			Ctx:         b.Ctx,
			Body:        b.Body,
			MaxRequests: b.MaxRequests,
			Headers:     b.Headers,
		})
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

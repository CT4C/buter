package dispatcher

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/edpryk/buter/internal/helpers/prepare"
	"github.com/edpryk/buter/internal/helpers/transform"
	"github.com/edpryk/buter/internal/modules/payloader"
	"github.com/edpryk/buter/internal/modules/reporter"
	"github.com/edpryk/buter/internal/modules/requester"
	"github.com/edpryk/buter/lib/transport"
)

func attackWithPayload(ctx context.Context, config AttackConfig) {
	wg := &sync.WaitGroup{}

	totalPayloads, payloadSet, err := prepare.PreparePayloads(config.PayloadFiles)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	attackValues := []string{config.Url, config.Headers.String(), config.Body.String()}

	attackValueString := strings.Join(attackValues, prepare.AttackValueSeparator)

	Payloader := payloader.New(payloader.Config{
		AttackValue:   attackValueString,
		AttackType:    config.AttackType,
		PayloadSet:    payloadSet,
		TotalPayloads: totalPayloads,
		Ctx:           ctx,
	})

	payloadProvider, _ := Payloader.PrepareAttack()

	requestWorker := requester.NewRequestQueue(requester.QueueWorkerConfig{
		MaxConcurrentRequests: config.MaxConcurrent,
		Ctx:                   ctx,
		Delay:                 config.Delay,
		Retries:               config.Retries,
	})

	requestConsumer, responseProvider, _ := requestWorker.Run()
	reporter := reporter.New()

	wg.Add(1)
	go func() {
		defer wg.Done()

		transport.MutableTransit(
			payloadProvider,
			requestConsumer,
			func(srcValue payloader.CraftedPayload) requester.RequestParameters {
				return requester.RequestParameters{
					Url:      srcValue.Url,
					Method:   config.Method,
					Header:   srcValue.Headers,
					Payloads: srcValue.Payloads,
					Body:     transform.NewMapStringer(srcValue.Body),
				}
			},
			time.Duration(0),
		)

		close(requestConsumer)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		reporter.StartWorker(responseProvider, nil)
	}()

	wg.Wait()
	config.AttackCompletedSig <- 0
}

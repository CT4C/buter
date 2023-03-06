package runner

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/edpryk/buter/cli"
	"github.com/edpryk/buter/internal/buter"
	"github.com/edpryk/buter/internal/helpers/prepare"
	"github.com/edpryk/buter/internal/reporter"
	"github.com/edpryk/buter/lib/transport"
	"github.com/edpryk/buter/pkg/requester"
)

type AttackConfig struct {
	AttackCompletedSig chan int

	cli.UserConfig
}

var headers = map[string]string{
	"User-Agent": "unknown",
}

func RunAttack(ctx context.Context, config AttackConfig) {
	wg := &sync.WaitGroup{}
	requestWorker := requester.NewRequestQueue(requester.QueueWorkerConfig{
		MaxConcurrentRequests: config.MaxConcurrent,
		Ctx:                   ctx,
		Delay:                 config.Delay,
		Retries:               config.Retries,
	})
	consumer, provider, _ := requestWorker.Run()
	reporter := reporter.New()

	totalPayloads, payloadSet, err := prepare.PreparePayloads(config.PayloadFiles)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// merge headers
	for key := range config.Headers {
		headers[key] = config.Headers[key]
	}

	Buter := buter.New(buter.Config{
		Ctx:           ctx,
		Url:           config.Url,
		Body:          config.Body,
		Headers:       headers,
		PayloadSet:    payloadSet,
		AttackType:    config.AttackType,
		MaxRequests:   config.DosRequest,
		TotalPayloads: totalPayloads,
	})

	payloadProvider, _ := Buter.RunPrepareAttack()

	wg.Add(1)
	go func() {
		defer wg.Done()

		transport.MutableTransit(
			payloadProvider,
			consumer,
			func(srcValue buter.CraftedPayload) requester.RequestParameters {
				return requester.RequestParameters{
					Url:      srcValue.Url,
					Method:   config.Method,
					Header:   srcValue.Headers,
					Payloads: srcValue.Payloads,
					Body:     srcValue.Body,
				}
			},
			time.Duration(0),
		)

		close(consumer)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		reporter.StartWorker(provider, config.Filters, config.Stop, config.AttackCompletedSig)
	}()

	wg.Wait()
	config.AttackCompletedSig <- 0
}

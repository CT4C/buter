package runner

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/edpryk/buter/cli"
	"github.com/edpryk/buter/internal/buter"
	"github.com/edpryk/buter/internal/connectors/payload"
	"github.com/edpryk/buter/internal/helpers/prepare"
	"github.com/edpryk/buter/internal/reporter"
	"github.com/edpryk/buter/pkg/requester"
)

type AttackConfig struct {
	AttackCompletedSig chan int

	cli.UserConfig
}

func RunAttack(ctx context.Context, config AttackConfig) {
	// merge headers
	for key := range config.Headers {
		headers[key] = config.Headers[key]
	}
	config.Headers = headers

	cli.PrintConfig(config.UserConfig)

	wg := &sync.WaitGroup{}
	errorQueue := make(chan error, 1)

	requestWorker := requester.NewRequestQueue(requester.QueueWorkerConfig{
		MaxConcurrentRequests: config.MaxConcurrent,
		Ctx:                   ctx,
		Delay:                 config.Delay,
		Retries:               config.Retries,
	})
	requestConsumer, responseProvider, _ := requestWorker.Run()
	reporter := reporter.New()

	totalPayloads, payloadSet, err := prepare.PreparePayloads(config.PayloadFiles)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	payloadFactory := buter.NewFactory(buter.Config{
		Ctx:           ctx,
		Url:           config.Url,
		Body:          config.Body,
		Headers:       config.Headers,
		PayloadSet:    payloadSet,
		AttackType:    config.AttackType,
		MaxRequests:   config.DosRequest,
		TotalPayloads: totalPayloads,
	})

	payloadConsumer := payload.NewPayloadConsumer(requestConsumer, errorQueue, config.Method)

	payloadFactory.Launch(payloadConsumer)

	wg.Add(1)
	go func() {
		defer wg.Done()
		reporter.StartWorker(responseProvider, config.Filters, config.Stop, config.AttackCompletedSig)
	}()

	wg.Wait()
	config.AttackCompletedSig <- 0
}

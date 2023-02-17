package dispatcher

import (
	"context"
	"sync"

	"github.com/edpryk/buter/cli"
	"github.com/edpryk/buter/internal/reporter"
	"github.com/edpryk/buter/internal/requester"
)

func dosAttack(ctx context.Context, config AttackConfig) {
	wg := &sync.WaitGroup{}

	requestWorker := requester.NewRequestQueue(requester.QueueWorkerConfig{
		Ctx:                   ctx,
		Delay:                 config.Delay,
		Retries:               config.Retries,
		RetryDelay:            config.RetryDelay,
		MaxConcurrentRequests: config.MaxConcurrent,
	})
	reporter := reporter.New()
	consumer, provider, _ := requestWorker.Run()

	wg.Add(1)
	go func(ctx context.Context, config cli.UserConfig) {
		defer wg.Done()

		for i := 0; i < config.DosRequest; i++ {
			consumer <- requester.RequestParameters{
				Url:      config.Url,
				Body:     config.Body,
				Method:   config.Method,
				Header:   config.Headers,
				Payloads: nil,
			}
		}
		close(consumer)
	}(ctx, config.UserConfig)

	wg.Add(1)
	go func() {
		defer wg.Done()
		reporter.StartWorker(provider, config.Filters, config.Stop, config.AttackCompletedSig)
	}()

	wg.Wait()
	config.AttackCompletedSig <- 0
}
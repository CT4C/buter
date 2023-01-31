package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/edpryk/buter/cli"
	"github.com/edpryk/buter/internal/helpers/prepare"
	"github.com/edpryk/buter/internal/helpers/transform"
	"github.com/edpryk/buter/internal/modules/payloader"
	"github.com/edpryk/buter/internal/modules/reporter"
	"github.com/edpryk/buter/internal/modules/requester"
	"github.com/edpryk/buter/lib/transport"
)

var (
	config             payloader.Config
	configs            cli.Input
	payloadSet         [][]string
	totalPayloads      int
	rootContext        context.Context
	cancelRootContetxt context.CancelFunc

	err error

	wg     = &sync.WaitGroup{}
	mut    = &sync.Mutex{}
	sigEnd = make(chan os.Signal)
)

func main() {
	signal.Notify(sigEnd, syscall.SIGINT)
	log.SetFlags(2)

	/*
		Need to test target connection before start
	*/

	configs = cli.ParseFlags()
	totalPayloads, payloadSet, err = prepare.PreparePayloads(configs.PayloadFiles)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if configs.Timeout > 0 {
		rootContext, cancelRootContetxt = context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	} else {
		rootContext, cancelRootContetxt = context.WithCancel(context.Background())
	}
	defer cancelRootContetxt()

	// TODO: move to separated func or method

	attackValues := []string{configs.Url, configs.Headers.String(), configs.Body.String()}

	attackValueString := strings.Join(attackValues, prepare.AttackValueSeparator)

	Payloader := payloader.New(payloader.Config{
		AttackValue:   attackValueString,
		AttackType:    configs.AttackType,
		PayloadSet:    payloadSet,
		TotalPayloads: totalPayloads,
		Ctx:           rootContext,
	})

	payloadProvider, _ := Payloader.PrepareAttack()
	// if err != <-errQ {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	requestQueueWorker := requester.NewRequestQueue(requester.QueueWorkerConfig{
		MaxConcurrentRequests: configs.MaxConcurrent,
		Ctx:                   rootContext,
		Delay:                 configs.Delay,
		Retries:               configs.Retries,
	})

	requestConsumer, responseProvider, _ := requestQueueWorker.Run()
	reporter := reporter.New()
	wg.Add(1)
	go func() {
		defer wg.Done()

		transport.MutableTransit(
			payloadProvider,
			requestConsumer,
			func(srcValue payloader.CraftedPayload) requester.ReuqestParameters {
				return requester.ReuqestParameters{
					Url:      srcValue.Url,
					Method:   configs.Method,
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

	go wg.Wait()

	select {
	case <-sigEnd:
		log.Println("Closed by Interruption")
		cancelRootContetxt()
		time.Sleep(2 * time.Second)
		// 	os.Exit(0)
		// case <-rootContext.Done():
		// 	log.Println("Canceled")
		// 	os.Exit(0)
	}
}

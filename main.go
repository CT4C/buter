package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/edpryk/buter/cli"
	"github.com/edpryk/buter/internal/helpers/prepare"
	"github.com/edpryk/buter/internal/modules/payloader"
	"github.com/edpryk/buter/internal/modules/reporter"
	"github.com/edpryk/buter/internal/modules/requester"
	"github.com/edpryk/buter/lib/transport"
)

var (
	config        payloader.Config
	userInput     cli.Input
	payloadSet    [][]string
	totalPayloads int
	wg            = &sync.WaitGroup{}
	mut           = &sync.Mutex{}

	err error
)

func main() {
	log.SetFlags(2)

	/*
		Need to test target connection before start
	*/

	userInput = cli.ParseFlags()
	totalPayloads, payloadSet, err = prepare.PreparePayloads(userInput.PayloadFiles)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootContext, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()

	attackValue := userInput.Url + prepare.AttackValueSeparator + userInput.Headers.String()

	Payloader := payloader.New(payloader.Config{
		AttackValue:   attackValue,
		AttackType:    userInput.AttackType,
		PayloadSet:    payloadSet,
		TotalPayloads: totalPayloads,
		Ctx:           rootContext,
	})

	payloadProvider, _ := Payloader.PrepareAttack()
	// if err != <-errQ {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	queueWorker := requester.NewRequestQueue(requester.QueueWorkerConfig{
		MaxConcurrentRequests: userInput.MaxConcurrent,
		Ctx:                   rootContext,
		Delay:                 userInput.Delay,
		Retries:               userInput.Retries,
	})

	requestConsumer, responseProvider, _ := queueWorker.Run()
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
					Method:   userInput.Method,
					Header:   srcValue.Headers,
					Payloads: srcValue.Payloads,
					Body:     nil,
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
}

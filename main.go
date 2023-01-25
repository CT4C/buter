package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/edpryk/buter/internal/docs"
	"github.com/edpryk/buter/internal/helpers/prepare"
	"github.com/edpryk/buter/internal/modules/payloader"
	"github.com/edpryk/buter/internal/modules/reporter"
	"github.com/edpryk/buter/internal/modules/requester"
	"github.com/edpryk/buter/lib/transport"
)

var (
	config        payloader.Config
	userInput     docs.Input
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

	userInput = docs.ParseFlags()
	totalPayloads, payloadSet, err = prepare.PreparePayloads(userInput.PayloadFiles)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootContext, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()

	paylaodConsumer := make(chan payloader.CraftedPayload, userInput.MaxConcurrent)
	statuses := make(chan payloader.ProcessStatus, 1)

	attackValue := userInput.Url + prepare.AttackValueSeparator + userInput.Headers

	Payloader := payloader.New(payloader.Config{
		AttackValue:     attackValue,
		AttackType:      userInput.AttackType,
		PayloadSet:      payloadSet,
		TotalPayloads:   totalPayloads,
		Ctx:             rootContext,
		PayloadConsumer: paylaodConsumer,
		StatusChan:      statuses,
	})

	// Return payload provider instead of passing payloadConsumer
	err = Payloader.PrepareAttack()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	queueWorker := requester.NewRequestQueue(requester.QueueWorkerConfig{
		MaxConcurrentRequests: userInput.MaxConcurrent,
		Ctx:                   rootContext,
		Delay:                 userInput.Delay,
		Retries:               3,
	})

	requestConsumer, responseProvider, _ := queueWorker.Run()
	reporter := reporter.New()

	wg.Add(1)
	go func() {
		defer wg.Done()

		transport.MutableTransit(
			paylaodConsumer,
			requestConsumer,
			func(original payloader.CraftedPayload) requester.ReuqestParameters {
				return requester.ReuqestParameters{
					Url:      original.Url,
					Method:   userInput.Method,
					Header:   original.Headers,
					Payloads: original.Payloads,
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

	// for status := range statuses {
	// 	log.Println(status.Message)
	// 	if status.Err {
	// 		os.Exit(1)
	// 	}
	// }

	wg.Wait()
}

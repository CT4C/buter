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
	"github.com/edpryk/buter/internal/modules/requester"
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

	config = payloader.Config{
		AttackValue:     attackValue,
		AttackType:      userInput.AttackType,
		PayloadSet:      payloadSet,
		TotalPayloads:   totalPayloads,
		Ctx:             rootContext,
		PayloadConsumer: paylaodConsumer,
		StatusChan:      statuses,
	}

	Payloader := payloader.New(config)

	err = Payloader.PrepareAttack()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ticker := time.NewTicker(time.Duration(userInput.Delay * int(time.Millisecond)))

	requetsQueue := requester.NewRequestQueue(requester.RequestQueueConfig{
		MaxConcurrentRequests: userInput.MaxConcurrent,
	})

	consumer, provider, _ := requetsQueue.StartWorker()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for craftedPayload := range paylaodConsumer {
			consumer <- requester.ReuqestParameters{
				Method: userInput.Method,
				Url:    craftedPayload.Url,
				Header: craftedPayload.Headers,
				Body:   nil,
			}

			<-ticker.C
		}

		close(consumer)
	}()

	for res := range provider {
		fmt.Println(res.StatusCode, res.Request.URL)
	}

	for status := range statuses {
		log.Println(status.Message)
		if status.Err {
			os.Exit(1)
		}
	}

	wg.Wait()
	ticker.Stop()
}

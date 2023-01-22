package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/edpryk/buter/internal/docs"
	"github.com/edpryk/buter/internal/helpers/definer"
	"github.com/edpryk/buter/internal/helpers/prepare"
	"github.com/edpryk/buter/internal/modules/payloader"
	"github.com/edpryk/buter/internal/modules/requester"
)

var (
	config        payloader.Config
	userInput     docs.Input
	payloadSet    [][]string
	totalPayloads int
	wg            sync.WaitGroup

	err error
)

func main() {
	userInput = docs.ParseFlags()
	totalPayloads, payloadSet, err = prepare.PreparePayloads(userInput.PayloadFiles)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	method := "get"

	rootContext, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()

	paylaodConsumer := make(chan payloader.CraftedPayload, userInput.ThreadsInTime)
	statuses := make(chan payloader.ProcessStatus, 1)

	attackValue := userInput.Url + definer.AttackValueSeparator + userInput.Headers

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

	go func() {
		for craftedPayload := range paylaodConsumer {
			wg.Add(1)

			attackValue, err := definer.ParseAttackValues(craftedPayload.Value)
			if err != nil {
				statuses <- payloader.ProcessStatus{
					Err:     true,
					Message: err.Error(),
				}
				continue
			}

			/*
				Throttle requests
			*/
			time.Sleep(time.Duration(userInput.Delay * int(time.Millisecond)))
			go func(m, u string, h map[string]string, payload payloader.CraftedPayload) {

				defer wg.Done()
				reqStartTime := time.Now()

				res, err := requester.Do(m, u, h, nil)
				if err != nil {
					fmt.Println(err)
					return
				}
				/*
					Need to detect body length
				*/
				report := "Time: %-12s Status: %-5d Length: %-5d "

				for number, payloadValue := range payload.Payloads {
					report += fmt.Sprintf("P_%d: %-5s ", number+1, payloadValue)
				}

				// Thist must be sent to Reporter channel, another entity reponsibile for printing report
				fmt.Printf(report+"\n", time.Since(reqStartTime), res.StatusCode, res.ContentLength)
			}(method, attackValue.Url, attackValue.Headers, craftedPayload)
		}
	}()

	for status := range statuses {
		log.Println(status.Message)
		if status.Err {
			os.Exit(1)
		}
	}

	wg.Wait()
}

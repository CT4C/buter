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
	wg            = &sync.WaitGroup{}
	mut           = &sync.Mutex{}

	err error
)

func doRequest(m, u string, h map[string]string, payload payloader.CraftedPayload) {
	defer wg.Done()
	// fmt.Println(payload)

	reqStartTime := time.Now()
	res, err := requester.Do(m, u, h, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	/*
		Need to detect body length
	*/
	report := "Time: %-12s Status: %-5d Length: %-5d"
	// URL: %-5s Headers: %s"

	payloadsText := ""
	for number, payloadValue := range payload.Payloads {
		payloadsText += fmt.Sprintf("%-3sP_%d: %-5s ", "", number+1, payloadValue)
	}

	report = payloadsText + report

	// Thist must be sent to Reporter channel, another entity reponsibile for printing report
	log.Printf(report+"\n", time.Since(reqStartTime), res.StatusCode, res.ContentLength)
}

func main() {
	log.SetFlags(2)
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

	ticker := time.NewTicker(time.Duration(userInput.Delay * int(time.Millisecond)))
	go func() {
		for craftedPayload := range paylaodConsumer {
			// fmt.Println(555, craftedPayload)

			attackValue, err := definer.ParseAttackValues(craftedPayload.Value)
			if err != nil {
				statuses <- payloader.ProcessStatus{
					Err:     true,
					Message: err.Error(),
				}
				continue
			}

			wg.Add(1)
			go doRequest(method, attackValue.Url, attackValue.Headers, craftedPayload)
			<-ticker.C
		}
	}()

	for status := range statuses {
		log.Println(status.Message)
		if status.Err {
			os.Exit(1)
		}
	}

	wg.Wait()
	ticker.Stop()
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/edpryk/buter/src/buter"
	"github.com/edpryk/buter/src/docs"
	"github.com/edpryk/buter/src/prepare"
)

var (
	config        buter.Config
	userInput     docs.Input
	payloadSet    [][]string
	totalPayloads int

	err error
)

func main() {

	userInput = docs.ParseFlags()
	totalPayloads, payloadSet, err = prepare.PreparePayloads(userInput.PayloadFiles)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootContext, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()

	paylaodConsumer := make(chan string, userInput.ThreadsInTime)
	statuses := make(chan buter.ProcessStatus, 1)

	attackValue := userInput.Url + "," + userInput.Headers

	config = buter.Config{
		AttackValue:     attackValue,
		AttackType:      userInput.AttackType,
		PayloadSet:      payloadSet,
		TotalPayloads:   totalPayloads,
		Ctx:             rootContext,
		PayloadConsumer: paylaodConsumer,
		StatusChan:      statuses,
	}

	Butter := buter.New(config)

	err = Butter.PrepareAttackValue()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go func() {
		// var wg sync.WaitGroup

		for payload := range paylaodConsumer {
			// wg.Add(1)
			// method := "GET"

			fmt.Println(payload)
			// go func(m, u string) {
			// 	requester.Do(m, u)
			// }(method, url)
		}
	}()

	for status := range statuses {
		log.Println(status.Message)
		if status.Err {
			os.Exit(1)
		}
	}
}

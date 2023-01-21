package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/edpryk/buter/src/buter"
	"github.com/edpryk/buter/src/docs"
	"github.com/edpryk/buter/src/prepare"
)

var (
	userInput   docs.Input
	variants    int
	payloadSet  [][]string
	config      buter.Config
	urlProvider buter.UrlProvider

	err error
)

func main() {

	userInput = docs.ParseFlags()
	variants, payloadSet, err = prepare.PreparePayloads(userInput.PayloadFiles)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootContext, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()

	config = buter.Config{
		Url:        userInput.Url,
		AttackType: userInput.AttackType,
		PayloadSet: payloadSet,
		Variants:   variants,
		Ctx:        rootContext,
	}

	urlProvider, err = buter.Run(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for url := range urlProvider {
		fmt.Println(url)
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/edpryk/buter/src/buter"
	"github.com/edpryk/buter/src/docs"
	"github.com/edpryk/buter/src/prepare"
)

func main() {
	userInput := docs.ParseFlags()
	variants, payloadSet, err := prepare.PreparePayloads(userInput.PayloadFiles)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config := ButerConfig{
		url:        userInput.Url,
		attackType: userInput.AttackType,
		payloadSet: payloadSet,
		variants:   variants,
	}

	buter.Run(config)
}

type ButerConfig struct {
	attackType string
	url        string
	payloadSet [][]string
	variants   int
	consumer   chan string
}

func (bc ButerConfig) Attack() string {
	return bc.attackType
}

func (bc ButerConfig) Url() string {
	return bc.url
}

func (bc ButerConfig) PayloadSet() [][]string {
	return ([][]string)(bc.payloadSet)
}

func (bc ButerConfig) Variants() int {
	return bc.variants
}

func (bc ButerConfig) Consume(url string) {

}

package docs

import (
	"flag"
	"fmt"
	"os"
)

type Input struct {
	Url        string
	AttackType string
	PayloadFiles
}

/*
Parse flgas validate it
and eather return it or
pritn an error and usage
and exis process
*/
func ParseFlags() Input {
	input := Input{}

	flag.Var(&input.PayloadFiles, payloadFlag, payloadUsage)
	flag.StringVar(&input.Url, urlFlag, defaultUrl, urlUsage)
	flag.StringVar(&input.AttackType, attackTypeFlag, defaultAttackType, attackTypeUsage)

	flag.Parse()

	if err := validateInput(input); err != nil {
		fmt.Println(err.Error())
		flag.Usage()
		os.Exit(1)
	}

	return input
}

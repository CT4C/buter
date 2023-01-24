package docs

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Input struct {
	Url           string
	AttackType    string
	Headers       string
	MaxConcurrent int
	Delay         int
	Method        string

	PayloadFiles
}

/*
Parse flags validate it
and either return it or
print an error and usage
and exits process
*/
func ParseFlags() Input {
	input := Input{}

	flag.Var(&input.PayloadFiles, payloadFlag, payloadUsage)
	flag.StringVar(&input.Url, urlFlag, defaultUrl, urlUsage)
	flag.StringVar(&input.Method, methodFlag, defaultMethod, methodUsage)
	flag.StringVar(&input.AttackType, attackTypeFlag, defaultAttackType, attackTypeUsage)
	flag.IntVar(&input.MaxConcurrent, threadsFlag, defaultThreads, threadsUsage)
	flag.StringVar(&input.Headers, headersFlag, defaultHeaders, headersUsage)
	flag.IntVar(&input.Delay, delayFlag, defaultDealy, delayUsage)

	flag.Parse()

	if len(input.Headers) > 0 {
		d := make(map[string]string)

		if err := json.Unmarshal([]byte(input.Headers), &d); err != nil {
			fmt.Println("Can't parse headers", err)
			os.Exit(1)
		}
		b, _ := json.Marshal(d)

		input.Headers = string(b)
	}

	if err := validateInput(input); err != nil {
		fmt.Println(err.Error())
		flag.Usage()
		os.Exit(1)
	}

	if input.Delay <= 0 {
		input.Delay = 1
	}

	return input
}

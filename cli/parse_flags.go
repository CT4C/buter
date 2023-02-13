package cli

import (
	"flag"
	"fmt"
	"os"
)

type UserConfig struct {
	Url           string
	AttackType    string
	MaxConcurrent int
	Delay         int
	Method        string
	Retries       int
	RetryDelay    int
	Timeout       int
	DosRequest    int

	Headers
	PayloadFiles
	Body *Body
}

/*
Parse flags validate it
and either return it or
print an error and usage
and exits process
*/
func ParseFlags() UserConfig {
	UserConfig := &UserConfig{
		Body: &Body{},
	}

	flag.Var(&UserConfig.PayloadFiles, payloadFlag, payloadUsage)
	flag.StringVar(&UserConfig.Url, urlFlag, defaultUrl, urlUsage)
	flag.StringVar(&UserConfig.Method, methodFlag, defaultMethod, methodUsage)
	flag.StringVar(&UserConfig.AttackType, attackTypeFlag, defaultAttackType, attackTypeUsage)
	flag.IntVar(&UserConfig.MaxConcurrent, threadsFlag, defaultThreads, threadsUsage)
	flag.Var(&UserConfig.Headers, headersFlag, headersUsage)
	flag.IntVar(&UserConfig.Delay, delayFlag, defaultDelay, delayUsage)
	flag.IntVar(&UserConfig.RetryDelay, retriesDelayFlag, defaultRetryDelay, retriyDelayUsage)
	flag.IntVar(&UserConfig.Retries, retriesAmountFlag, defaultRetriesAmount, retriesAmountUsage)
	flag.Var(UserConfig.Body, bodyFlag, bodyUsage)
	flag.IntVar(&UserConfig.Timeout, timeoutFlag, defaultTimeout, timeoutUsage)
	flag.IntVar(&UserConfig.DosRequest, dosRequestsFlag, defaultDosRequests, timeoutUsage)

	flag.Parse()

	if err := validateInput(UserConfig); err != nil {
		fmt.Println(err.Error())
		flag.Usage()
		os.Exit(1)
	}

	if UserConfig.Delay <= 0 {
		UserConfig.Delay = 1
	}

	return *UserConfig
}

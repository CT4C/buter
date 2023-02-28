package cli

import (
	"flag"
	"fmt"
	"os"
)

/*
	TODO: Add raw request support.
*/
type UserConfig struct {
	Url           string  `json:"url"`
	AttackType    string  `json:"attackType"`
	MaxConcurrent int     `json:"maxConcurrent"`
	Delay         int     `json:"delay"`
	Method        string  `json:"method"`
	Retries       int     `json:"retries"`
	RetryDelay    int     `json:"retryDelay"`
	Timeout       int     `json:"timeout"`
	DosRequest    int     `json:"dosMaxRequests"`
	Stop          Stopper `json:"stopper"`

	ConfigFile     string `json:"-"`
	ConfigTemplate int    `json:"-"`
	IsBatch        bool   `json:"-"`

	Filters      `json:"filters"`
	Headers      `json:"headers"`
	PayloadFiles `json:"payloadFiles"`
	Body         Body `json:"body"`
}

/*
Parse flags validate it
and either return it or
print an error and usage
and exits process
*/
func ParseFlags() []UserConfig {
	config := &UserConfig{
		Body: Body{},
		Filters: Filters{
			"length": make([]string, 0),
			"status": make([]string, 0),
		},
		Stop: Stopper{
			"status": make([]string, 0),
		},
	}

	flag.Var(&config.PayloadFiles, payloadFlag, payloadUsage)
	flag.StringVar(&config.Url, urlFlag, defaultUrl, urlUsage)
	flag.StringVar(&config.Method, methodFlag, defaultMethod, methodUsage)
	flag.StringVar(&config.AttackType, attackTypeFlag, defaultAttackType, attackTypeUsage)
	flag.StringVar(&config.ConfigFile, configFlag, defaultConfig, configUsage)
	flag.IntVar(&config.ConfigTemplate, configTemplateFlag, defaultConfigTemplate, configTemplateUsage)
	flag.IntVar(&config.MaxConcurrent, threadsFlag, defaultThreads, threadsUsage)
	flag.Var(&config.Headers, headersFlag, headersUsage)
	flag.IntVar(&config.Delay, delayFlag, defaultDelay, delayUsage)
	flag.IntVar(&config.RetryDelay, retriesDelayFlag, defaultRetryDelay, retryDelayUsage)
	flag.IntVar(&config.Retries, retriesAmountFlag, defaultRetriesAmount, retriesAmountUsage)
	flag.Var(&config.Body, bodyFlag, bodyUsage)
	flag.IntVar(&config.Timeout, timeoutFlag, defaultTimeout, timeoutUsage)
	flag.IntVar(&config.DosRequest, dosRequestsFlag, defaultDosRequests, timeoutUsage)
	flag.Var(&config.Filters, filterOutFlag, filterOutUsage)
	flag.Var(&config.Stop, stopFlag, stopUsage)

	flag.Parse()

	batchConfig := make([]UserConfig, 0)

	if config.ConfigTemplate > 0 {
		printConfigTemplate()
		os.Exit(0)
	}

	if err := validateInput(config); err != nil {
		fmt.Printf("%-10s %s\n", "Error:", err.Error())
		printUsage()
		os.Exit(0)
	}

	if config.ConfigFile != "" {
		configs, err := loadJSONConfig(config.ConfigFile)
		if err != nil {
			panic(err)
		}

		batchConfig = configs
	} else {
		batchConfig = append(batchConfig, *config)
	}

	if config.Delay <= 0 {
		config.Delay = 1
	}

	return batchConfig
}

package cli

import (
	"encoding/json"
	"fmt"
)

func printConfigTemplate() {
	config := []UserConfig{
		{
			Url:           defaultUrl,
			Stop:          Stopper{},
			Body:          Body{},
			Delay:         defaultDelay,
			Method:        defaultMethod,
			Filters:       Filters{},
			Headers:       Headers{},
			Retries:       defaultRetriesAmount,
			Timeout:       defaultTimeout,
			DosRequest:    defaultDosRequests,
			RetryDelay:    defaultRetryDelay,
			AttackType:    defaultAttackType,
			PayloadFiles:  PayloadFiles{},
			MaxConcurrent: defaultThreads,
		},
	}

	for _, f := range knownFilters {
		config[0].Filters[f] = []string{}
	}

	config[0].Stop.Set("status:200")
	config[0].Filters.Set("status:500,403,401")

	b, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}

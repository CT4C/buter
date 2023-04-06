package cli

import "fmt"

func PrintConfig(config UserConfig) {
	if len(config.Headers) > 0 {
		fmt.Println("-----HEADERS-----")
		for key, value := range config.Headers {
			fmt.Printf("%-15s: %s\n", key, value)
		}
	}

	if len(config.Body) > 0 {
		fmt.Println("-----BODY--------")
		fmt.Println(config.Body)
	}
	fmt.Println("-----------------")
}

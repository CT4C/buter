package cli

import "fmt"

func PrintConfig(config UserConfig) {
	fmt.Println("-----HEADERS-----")
	for key, value := range config.Headers {
		fmt.Printf("%-15s: %s\n", key, value)
	}
	fmt.Println("-----BODY--------")
	fmt.Println(config.Body)
	fmt.Println("-----------------")
}

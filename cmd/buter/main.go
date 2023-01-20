package main

import (
	"fmt"

	"github.com/edpryk/buter/src/docs"
)

func main() {
	userInput := docs.ParseFlags()
	fmt.Println(userInput)
	// config := struct{}{}

	// buter.Run(config)
}

package cli

import "fmt"

const author = "https://github.com/0x640x720x650x640x64"
const email = "https://github.com/0x640x720x650x640x64"
const version = "0.0.1"

func PrintInfo() {
	fmt.Println()
	fmt.Printf("%-10s %s\n", "Version:", version)
	fmt.Printf("%-10s %s\n", "Contact:", email)
	fmt.Printf("%-10s %s \n", "Author:", author)
	fmt.Println()
}

package cli

import "fmt"

const author = "Dmytro Prykhodko"
const email = "edpryk@gmail.com"
const version = "0.0.1"

func PrintInfo() {
	fmt.Println()
	fmt.Printf("%-10s %s\n", "Version:", version)
	fmt.Printf("%-10s %s\n", "Contact:", email)
	fmt.Printf("%-10s %s \n", "Author:", author)
	fmt.Println()
}

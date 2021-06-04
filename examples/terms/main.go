package main

import (
	"fmt"
	"os"

	"github.com/avocatl/admiral/pkg/prompter"
)

func main() {
	err := prompter.Confirm("Agree to the terms and conditions", nil)

	if err != nil {
		fmt.Println("You need to agree to the terms and conditions to continue.")
		os.Exit(1)
	}

	fmt.Println("Thanks for accepting our terms and conditions!")
}

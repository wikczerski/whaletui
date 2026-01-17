package cmd

import (
	"fmt"
	"strings"
)

// UserInteraction handles user input and interaction for command-line operations
type UserInteraction struct{}

func (ui UserInteraction) AskYesNo(prompt string) bool {
	fmt.Printf("%s (y/n): ", prompt)
	var response string
	_, _ = fmt.Scanln(&response)

	normalizedResponse := strings.ToLower(strings.TrimSpace(response))
	return normalizedResponse == "y" || normalizedResponse == "yes"
}

func (ui UserInteraction) WaitForEnter() {
	fmt.Println("Press Enter to exit...")
	_, _ = fmt.Scanln()
}

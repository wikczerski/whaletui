package cmd

import (
	"fmt"
	"strings"
)

type UserInteraction struct{}

func (ui UserInteraction) askYesNo(prompt string) bool {
	fmt.Printf("%s (y/n): ", prompt)
	var response string
	_, _ = fmt.Scanln(&response)

	normalizedResponse := strings.ToLower(strings.TrimSpace(response))
	return normalizedResponse == "y" || normalizedResponse == "yes"
}

func (ui UserInteraction) waitForEnter() {
	fmt.Println("Press Enter to exit...")
	_, _ = fmt.Scanln()
}

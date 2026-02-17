package security

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// PromptHidden displays a prompt and reads user input from the terminal
// without echoing the characters back to the screen. This is typically
// used for sensitive values such as passwords or master keys.
//
// The function writes the prompt, reads a single line of hidden input
// from stdin, and returns it as a string. A newline is printed after
// the input for clean terminal formatting. If reading from the terminal
// fails, an error is returned.
func PromptHidden(prompt string) (string, error) {
	fmt.Print(prompt)

	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // newline after input
	if err != nil {
		return "", err
	}

	return string(bytePassword), nil
}

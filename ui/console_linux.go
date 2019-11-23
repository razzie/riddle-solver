//+build linux
package ui

import (
	"fmt"
)

// SetConsoleTitle sets the console title
func SetConsoleTitle(title string) error {
	fmt.Printf("\033]0;%s\007", title)
	return nil
}

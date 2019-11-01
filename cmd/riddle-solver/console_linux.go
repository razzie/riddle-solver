//+build linux
package main

import (
	"fmt"
)

// SetConsoleTitle sets the console title
func SetConsoleTitle(title string) error {
	fmt.Printf("\033]0;Title goes here\007")
	return nil
}

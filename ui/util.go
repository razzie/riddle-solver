package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

func colorize(text string) string {
	primaryColor := tview.Styles.PrimaryTextColor.Hex()
	secondaryColor := tview.Styles.SecondaryTextColor.Hex()
	return fmt.Sprintf("[#%06x]%s[#%06x]", secondaryColor, text, primaryColor)
}

func colorizeItem(item string) string {
	parts := strings.SplitN(item, ":", 2)
	return fmt.Sprintf("%s:%s", parts[0], colorize(parts[1]))
}

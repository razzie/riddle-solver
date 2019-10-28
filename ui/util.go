package ui

import (
	"fmt"

	"github.com/razzie/riddle-solver/riddle"
	"github.com/rivo/tview"
)

func colorize(text string) string {
	primaryColor := tview.Styles.PrimaryTextColor.Hex()
	secondaryColor := tview.Styles.SecondaryTextColor.Hex()
	return fmt.Sprintf("[#%06x]%s[#%06x]", secondaryColor, text, primaryColor)
}

func colorizeItem(item riddle.Item) string {
	itemType, value := item.Split()
	return fmt.Sprintf("%s:%s", itemType, colorize(value))
}

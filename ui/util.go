package ui

import (
	"fmt"

	"github.com/razzie/riddle-solver/riddle"
	"github.com/rivo/tview"
)

func colorize(text string) string {
	color := tview.Styles.SecondaryTextColor.Hex()
	return fmt.Sprintf("[#%06x]%s[-]", color, text)
}

func colorizeItem(item riddle.Item) string {
	itemType, value := item.Split()
	return fmt.Sprintf("[-]%s:%s", itemType, colorize(value))
}

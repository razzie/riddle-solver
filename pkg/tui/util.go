package tui

import (
	"fmt"
	"strings"

	"github.com/razzie/riddle-solver/pkg/riddle"
)

func colorize(text string) string {
	normal := currentTheme.TextColor
	highlight := currentTheme.HighlightTextColor
	return fmt.Sprintf("[%s]%s[%s]", highlight, text, normal)
}

func colorizeItem(item riddle.Item) string {
	itemType, value := item.Split()
	return fmt.Sprintf("%s:%s", itemType, colorize(value))
}

func colorizeItems(itemType string, items []string) string {
	normal := currentTheme.TextColor
	highlight := currentTheme.HighlightTextColor
	list := strings.Join(items, fmt.Sprintf("[%s], [%s]", normal, highlight))
	return fmt.Sprintf("%s:[%s]%s[%s]", itemType, highlight, list, normal)
}

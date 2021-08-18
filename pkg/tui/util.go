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

func getAutocompleteFunc(words []string) func(string) []string {
	return func(currentText string) (results []string) {
		if len(currentText) == 0 {
			return
		}

		for _, word := range words {
			if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) {
				results = append(results, word)
			}
		}

		return
	}
}

func getAutocompleteItemsFunc(items []riddle.Item) func(string) []string {
	return func(currentText string) (results []string) {
		if len(currentText) == 0 {
			return
		}

		for _, item := range items {
			if strings.HasPrefix(strings.ToLower(string(item)), strings.ToLower(currentText)) {
				results = append(results, string(item))
			}
		}

		return
	}
}

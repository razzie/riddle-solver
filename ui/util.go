package ui

import (
	"fmt"
	"strings"

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

		if len(results) <= 1 {
			results = nil
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

		if len(results) <= 1 {
			results = nil
		}

		return
	}
}

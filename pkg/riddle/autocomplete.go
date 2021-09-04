package riddle

import (
	"strings"
)

func GetAutocompleteFunc(words []string) func(string) []string {
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

func GetAutocompleteItemsFunc(items []Item) func(string) []string {
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

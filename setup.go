package main

import (
	"fmt"
)

// Setup is a map of riddle item types and their possible values
type Setup map[string][]string

func hasDuplicates(values []string) bool {
	count := len(values)
	for i := 0; i < count-1; i++ {
		val := values[i]
		for j := i + 1; j < count; j++ {
			if values[j] == val {
				return true
			}
		}
	}

	return false
}

// GetItems returns all items in a slice in itemType:value format
func (setup Setup) GetItems() []string {
	var items []string
	for itemType, values := range setup {
		for _, val := range values {
			items = append(items, fmt.Sprintf("%s:%s", itemType, val))
		}
	}
	return items
}

// Check returns an error if the provided Setup is invalid
func (setup Setup) Check() error {
	valueCount := 0
	itemTypes := make([]string, 0, len(setup))

	for itemType, values := range setup {
		itemTypes = append(itemTypes, itemType)

		if hasDuplicates(values) {
			return fmt.Errorf("item type '%s' has duplicate values", itemType)
		}

		if valueCount == 0 {
			valueCount = len(values)
		} else if len(values) != valueCount {
			return fmt.Errorf("all item types should have an equal number of values")
		}
	}

	if hasDuplicates(itemTypes) {
		return fmt.Errorf("duplicate item types")
	}

	return nil
}

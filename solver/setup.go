package solver

import (
	"fmt"
)

// Setup is a map of riddle item types and their possible values
type Setup map[string][]string

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

// GetItemTypes returns all item types in a slice
func (setup Setup) GetItemTypes() []string {
	itemTypes := make([]string, 0, len(setup))
	for itemType := range setup {
		itemTypes = append(itemTypes, itemType)
	}
	return itemTypes
}

// Check returns an error if the provided Setup is invalid
func (setup Setup) Check() error {
	valueCount := 0
	itemTypes := make([]string, 0, len(setup))

	for itemType, values := range setup {
		itemTypes = append(itemTypes, itemType)

		if hasDuplicates(values) {
			return fmt.Errorf("Item type '%s' has duplicate values", itemType)
		}

		if valueCount == 0 {
			valueCount = len(values)
		} else if len(values) != valueCount {
			return fmt.Errorf("All item types should have an equal number of values")
		}
	}

	if hasDuplicates(itemTypes) {
		return fmt.Errorf("Duplicate item types")
	}

	return nil
}

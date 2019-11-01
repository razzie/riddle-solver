package riddle

import (
	"fmt"
	"strings"
)

// Item represents a riddle item in ItemType:Item format (like color:red)
type Item string

// NewItem returns a new Item
func NewItem(itemType, value string) Item {
	return Item(fmt.Sprintf("%s:%s", itemType, value))
}

// Split returns the item type and item as two separate strings
func (item Item) Split() (string, string) {
	parts := strings.SplitN(string(item), ":", 2)
	return parts[0], parts[1]
}

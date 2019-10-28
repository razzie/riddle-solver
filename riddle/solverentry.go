package riddle

// SolverEntry is a unit of riddle items that belong together
type SolverEntry map[string][]string

// NewSolverEntry returns a new SolverEntry
func NewSolverEntry(setup Setup) SolverEntry {
	entry := make(SolverEntry)
	for itemType, values := range setup {
		entry[itemType] = copySlice(values)
	}
	return entry
}

// Set sets an item as the only possible value for that type
func (entry SolverEntry) Set(item Item) bool {
	itemType, value := item.Split()
	values, _ := entry[itemType]
	if len(values) > 1 || values[0] != value {
		entry[itemType] = []string{value}
		return true
	}
	return false
}

// Unset removes an item from the possible values for that type
func (entry SolverEntry) Unset(item Item) bool {
	itemType, value := item.Split()
	values, _ := entry[itemType]
	for i, val := range values {
		if val == value {
			entry[itemType] = append(values[:i], values[i+1:]...)
			return true
		}
	}
	return false
}

// Contains returns whether the entry contains the specified item
func (entry SolverEntry) Contains(item Item) bool {
	itemType, value := item.Split()
	values, _ := entry[itemType]
	return contains(values, value)
}

// OnlyContains returns whether the entry only contains the specified item of that item type
func (entry SolverEntry) OnlyContains(item Item) bool {
	itemType, value := item.Split()
	values, _ := entry[itemType]
	return len(values) == 1 && values[0] == value
}

// GetExcludedItems returns the excluded items compared to the setup
func (entry SolverEntry) GetExcludedItems(setup Setup) (excluded []Item) {
	items := setup.GetItems()
	for _, item := range items {
		if !entry.Contains(item) {
			excluded = append(excluded, item)
		}
	}
	return
}

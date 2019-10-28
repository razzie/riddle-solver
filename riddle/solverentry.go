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
func (entry SolverEntry) Set(item Item) {
	itemType, value := item.Split()
	entry[itemType] = []string{value}
}

// Unset removes an item from the possible values for that type
func (entry SolverEntry) Unset(item Item) {
	itemType, value := item.Split()
	values, _ := entry[itemType]
	for i, val := range values {
		if val == value {
			entry[itemType] = append(values[:i], values[i+1:]...)
			return
		}
	}
}

// Contains returns whether the entry contains the specified item
func (entry SolverEntry) Contains(item Item) bool {
	itemType, value := item.Split()
	values, _ := entry[itemType]
	return contains(values, value)
}

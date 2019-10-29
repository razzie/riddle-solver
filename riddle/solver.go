package riddle

// Solver can solve a riddle based on the setup and rules
type Solver struct {
	Entries []SolverEntry
}

// NewSolver returns a new Solver
func NewSolver(setup Setup) *Solver {
	entryCount := setup.GetItemCountPerType()
	entries := make([]SolverEntry, 0, entryCount)
	var primaryItemType string
	var primaryItemTypeValues []string

	for itemType, values := range setup {
		primaryItemType = itemType
		primaryItemTypeValues = values
		break
	}

	for i := 0; i < entryCount; i++ {
		entry := NewSolverEntry(setup)
		entry[primaryItemType] = []string{primaryItemTypeValues[i]}
		entries = append(entries, entry)
	}

	return &Solver{Entries: entries}
}

// ApplyRules applies the provided rules to reduce the item variations as much as possible
func (solver *Solver) ApplyRules(rules []Rule) {
	simpleRules, conditionalRules := SplitRules(rules)

	// looping until no rules make any change
	for changed := true; changed; changed = false {
		for i, entry := range solver.Entries {
			// applying simple rules
			for _, rule := range simpleRules {
				changed = rule.ApplySimple(entry) || changed
			}
			// running all variations of entryA and entryB
			for j := i; j < len(solver.Entries); j++ {
				// applying conditional rules
				for _, rule := range conditionalRules {
					changed = rule.ApplyConditional(entry, solver.Entries[j]) || changed
				}
			}
		}
	}
}

// FindEntriesWithItem returns the entries that contain the specified item
func (solver *Solver) FindEntriesWithItem(item Item) (entries []SolverEntry) {
	for _, entry := range solver.Entries {
		if entry.Contains(item) {
			entries = append(entries, entry)
		}
	}
	return
}

func mergeEntries(entries []SolverEntry) map[string][]string {
	result := make(map[string][]string)
	for _, entry := range entries {
		for itemType, values := range entry {
			resultValues, _ := result[itemType]
			result[itemType] = uniqueAppend(resultValues, values)
		}
	}
	return result
}

// FindAssociatedItems returns a map of items that possible belong to the provided item
func (solver *Solver) FindAssociatedItems(item Item) map[string][]string {
	entries := solver.FindEntriesWithItem(item)
	result := mergeEntries(entries)
	itemType, _ := item.Split()
	delete(result, itemType)
	return result
}

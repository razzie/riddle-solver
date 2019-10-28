package riddle

// Solver can solve a riddle based on the setup and rules
type Solver struct {
	Entries []SolverEntry
}

// NewSolver returns a new Solver
func NewSolver(setup Setup) *Solver {
	entryCount := setup.GetItemCountPerType()
	entries := make([]SolverEntry, 0, entryCount)
	for i := 0; i < entryCount; i++ {
		entries = append(entries, NewSolverEntry(setup))
	}

	return &Solver{Entries: entries}
}

// ApplyRules applies the provided rules to reduce the item variations as much as possible
func (solver *Solver) ApplyRules(rules []Rule) {

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

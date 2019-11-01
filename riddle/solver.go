package riddle

import "fmt"

// Solver can solve a riddle based on the setup and rules
type Solver struct {
	Entries []SolverEntry
	setup   Setup
	rules   []Rule
}

// NewSolver returns a new Solver
func NewSolver(setup Setup, rules []Rule) *Solver {
	entryCount := setup.GetItemCountPerType()
	entries := make([]SolverEntry, 0, entryCount)

	for i := 0; i < entryCount; i++ {
		entry := NewSolverEntry(setup)
		entries = append(entries, entry)
	}

	return &Solver{
		Entries: entries,
		setup:   setup,
		rules:   rules,
	}
}

// GuessPrimaryItemType tries to guess the primary item type as it's important for conditional rules
func (solver *Solver) GuessPrimaryItemType() string {
	for _, rule := range solver.rules {
		if len(rule.ConditionItemType) > 0 {
			return rule.ConditionItemType
		}
	}

	if len(solver.rules) > 0 {
		itemType, _ := solver.rules[0].ItemA.Split()
		return itemType
	}

	return ""
}

// Solve tries to solve the riddle by applying all the rule to all entries over and over
// I like to call it brute-force deduction
func (solver *Solver) Solve(primaryItemType string) (steps int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	if len(solver.rules) == 0 {
		return 0, fmt.Errorf("There are no rules")
	}

	primaryItemTypeValues, _ := solver.setup[primaryItemType]
	for i, entry := range solver.Entries {
		entry[primaryItemType] = []string{primaryItemTypeValues[i]}
	}

	simpleRules, conditionalRules := SplitRules(solver.rules)

	// looping until no rules make any change
	for {
		changed := false
		steps++

		for i, entry := range solver.Entries {
			// applying simple rules
			for _, rule := range simpleRules {
				changed = rule.ApplySimple(entry) || changed
			}
			// applying conditional rules on  all variations of entryA and entryB
			for _, rule := range conditionalRules {
				for j := i; j < len(solver.Entries); j++ {
					changed = rule.ApplyConditional(entry, solver.Entries[j]) || changed
				}
			}
		}

		if !changed {
			break
		}
	}

	return
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

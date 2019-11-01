package riddle

import (
	"fmt"

	"github.com/antonmedv/expr"
)

// Rule is used by the riddle solver algorithm
type Rule struct {
	ItemA             Item
	ItemB             Item
	Relation          Relation
	Condition         string `json:",omitempty"`
	ConditionItemType string `json:",omitempty"`
	// TODO: IsReversible or ReverseCondition
}

// Check returns an error if the provided Rule is invalid
func (rule *Rule) Check(setup Setup) error {
	if len(rule.ItemA) == 0 {
		return fmt.Errorf("Item A is missing")
	}
	if len(rule.ItemB) == 0 {
		return fmt.Errorf("Item B is missing")
	}
	if rule.ItemA == rule.ItemB {
		return fmt.Errorf("Item A and B cannot be the same")
	}

	if len(rule.Condition) > 0 && len(rule.ConditionItemType) == 0 {
		return fmt.Errorf("Condition item type missing")
	}

	itemTypeA, _ := rule.ItemA.Split()
	itemTypeB, _ := rule.ItemB.Split()
	if itemTypeA == itemTypeB {
		return fmt.Errorf("Item A and B cannot have the same type")
	}

	if !setup.Contains(rule.ItemA) {
		return fmt.Errorf("Item A is invalid")
	}
	if !setup.Contains(rule.ItemB) {
		return fmt.Errorf("Item B is invalid")
	}

	if len(rule.ConditionItemType) > 0 {
		itemTypes := setup.GetItemTypes()
		if !contains(itemTypes, rule.ConditionItemType) {
			return fmt.Errorf("Condition item type is invalid")
		}
	}

	return nil
}

// HasCondition returns whether a rule has a condition
func (rule *Rule) HasCondition() bool {
	return len(rule.Condition) > 0
}

// ApplySimple tries to apply a simple (non-conditional) rule to a SolverEntry
func (rule *Rule) ApplySimple(entry SolverEntry) bool {
	switch rule.Relation {
	case RelAssociated:
		if entry.OnlyContains(rule.ItemA) {
			return entry.Set(rule.ItemB)
		}
		if entry.OnlyContains(rule.ItemB) {
			return entry.Set(rule.ItemA)
		}
		if !entry.Contains(rule.ItemA) {
			return entry.Unset(rule.ItemB)
		}
		if !entry.Contains(rule.ItemB) {
			return entry.Unset(rule.ItemA)
		}

	case RelDisassociated:
		if entry.OnlyContains(rule.ItemA) {
			return entry.Unset(rule.ItemB)
		}
		if entry.OnlyContains(rule.ItemB) {
			return entry.Unset(rule.ItemA)
		}
	}

	return false
}

// ApplyConditional tries to apply a conditional rule to a SolverEntry
func (rule *Rule) ApplyConditional(entryA SolverEntry, others []SolverEntry) bool {
	A := entryA.GetValue(rule.ConditionItemType)
	if A == nil {
		return false
	}

	var matched []SolverEntry
	var unmatched []SolverEntry

	for _, entryB := range others {
		B := entryB.GetValue(rule.ConditionItemType)
		if B == nil {
			continue
		}

		if rule.isConditionMatching(A, B) {
			matched = append(matched, entryB)
		} else {
			unmatched = append(unmatched, entryB)
		}
	}

	switch rule.Relation {
	case RelAssociated:
		if entryA.OnlyContains(rule.ItemA) {
			return unsetMany(unmatched, rule.ItemB)
		}

	case RelDisassociated:
		if entryA.OnlyContains(rule.ItemA) {
			return unsetMany(matched, rule.ItemB)
		}
	}

	return false
}

func (rule *Rule) isConditionMatching(A, B interface{}) bool {
	environment := map[string]interface{}{
		"A": A,
		"B": B,
	}

	output, err := expr.Eval(rule.Condition, environment)
	if err != nil {
		panic(err)
	}

	result, ok := output.(bool)
	if !ok {
		panic(fmt.Errorf("'%s' output is not bool", rule.Condition))
	}

	return result
}

func unsetMany(entries []SolverEntry, item Item) bool {
	changed := false
	for _, entry := range entries {
		changed = entry.Unset(item) || changed
	}
	return changed
}

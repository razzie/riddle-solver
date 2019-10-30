package riddle

import (
	"fmt"
	"strconv"

	"github.com/antonmedv/expr"
)

// Rule is used by the riddle solver algorithm
type Rule struct {
	ItemA             Item
	ItemB             Item
	Relation          Relation
	Condition         string `json:",omitempty"`
	ConditionItemType string `json:",omitempty"`
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
func (rule *Rule) ApplyConditional(entryA, entryB SolverEntry) bool {
	switch rule.Relation {
	case RelAssociated:
		if entryA.OnlyContains(rule.ItemA) && rule.testCondition(entryA, entryB) {
			return entryB.Set(rule.ItemB)
		}
		if entryB.OnlyContains(rule.ItemB) && rule.testCondition(entryB, entryA) {
			return entryA.Set(rule.ItemA)
		}
		if !entryA.Contains(rule.ItemA) && rule.testCondition(entryA, entryB) {
			return entryB.Unset(rule.ItemB)
		}
		if !entryB.Contains(rule.ItemB) && rule.testCondition(entryB, entryA) {
			return entryA.Unset(rule.ItemA)
		}

	case RelDisassociated:
		// ???
	}

	return false
}

func (rule *Rule) testCondition(entryA, entryB SolverEntry) bool {
	valuesA, _ := entryA[rule.ConditionItemType]
	valuesB, _ := entryB[rule.ConditionItemType]

	if len(valuesA) != 1 || len(valuesB) != 1 {
		return false
	}

	var A interface{}
	if val, err := strconv.Atoi(valuesA[0]); err != nil {
		A = valuesA[0]
	} else {
		A = val
	}

	var B interface{}
	if val, err := strconv.Atoi(valuesB[0]); err != nil {
		B = valuesA[0]
	} else {
		B = val
	}

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

// SplitRules splits a slice of rules to slices of simple and conditional rules
func SplitRules(rules []Rule) (simple []Rule, conditional []Rule) {
	for _, rule := range rules {
		if len(rule.Condition) > 0 {
			conditional = append(conditional, rule)
		} else {
			simple = append(simple, rule)
		}
	}
	return
}

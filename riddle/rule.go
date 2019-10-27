package riddle

import (
	"fmt"
)

// Rule is used by the riddle solver algorithm
type Rule struct {
	ItemA             string
	ItemB             string
	Relation          Relation
	Condition         string
	ConditionItemType string
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

	items := setup.GetItems()
	if !contains(items, rule.ItemA) {
		return fmt.Errorf("Item A is invalid")
	}
	if !contains(items, rule.ItemB) {
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

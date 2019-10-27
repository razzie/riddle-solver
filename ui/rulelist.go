package ui

import (
	"fmt"

	"github.com/razzie/riddle-solver/solver"
	"github.com/rivo/tview"
)

// RuleList is a UI element that contains the list of rules
type RuleList struct {
	*tview.List
	rules    []*solver.Rule
	editFunc func(*solver.Rule)
	saveFunc func([]solver.Rule)
	modal    ModalHandler
}

// NewRuleList returns a new RuleList
func NewRuleList(modal ModalHandler) *RuleList {
	return &RuleList{
		List:  tview.NewList(),
		modal: modal}
}

// HandleSetup filters the list based on the new setup
func (l *RuleList) HandleSetup(setup solver.Setup) {

}

// SaveRule adds a new rule to the list or updates an existing one
func (l *RuleList) SaveRule(rule *solver.Rule) {
	for i, r := range l.rules {
		if r == rule {
			l.RemoveItem(i)
			l.addRule(rule, i)
			return
		}
	}

	l.addRule(rule, -1)
}

func (l *RuleList) addRule(rule *solver.Rule, index int) {
	text := fmt.Sprintf("%s - %s - %s", rule.ItemA, rule.ItemB, rule.Relation.String())
	selected := func() {
		if l.editFunc != nil {
			l.editFunc(rule)
		}
	}

	if len(rule.Condition) > 0 {
		secondaryText := fmt.Sprintf("[ %s ] %s", rule.ConditionItemType, rule.Condition)
		l.InsertItem(index, text, secondaryText, 0, selected)
	} else {
		l.InsertItem(index, text, "", 0, selected)
	}
}

// SetEditFunc sets a function that gets called on the selected rule
func (l *RuleList) SetEditFunc(editFunc func(*solver.Rule)) {
	l.editFunc = editFunc
}

// SetSaveFunc sets a function that gets the list of all rules upon an update
func (l *RuleList) SetSaveFunc(saveFunc func([]solver.Rule)) {
	l.saveFunc = saveFunc
}

// Reset resets the list
func (l *RuleList) Reset() {

}

package ui

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/razzie/riddle-solver/riddle"
	"github.com/rivo/tview"
)

// RuleList is a UI element that contains the list of rules
type RuleList struct {
	*tview.List
	rules    []*riddle.Rule
	editFunc func(*riddle.Rule)
	saveFunc func([]riddle.Rule)
	modal    ModalHandler
}

// NewRuleList returns a new RuleList
func NewRuleList(modal ModalHandler) *RuleList {
	l := &RuleList{
		List:  tview.NewList().ShowSecondaryText(false),
		modal: modal,
	}

	l.SetInputCapture(l.handleInput)

	return l
}

// HandleSetup filters the list based on the new setup
func (l *RuleList) HandleSetup(setup riddle.Setup) {
	toBeRemoved := make([]int, 0, len(l.rules))
	removeCount := 0

	for i, rule := range l.rules {
		if err := rule.Check(setup); err != nil {
			toBeRemoved = append(toBeRemoved, i)
		}
	}

	for _, index := range toBeRemoved {
		l.removeRule(index-removeCount, false)
		removeCount++
	}

	if removeCount > 0 {
		l.Save()
	}
}

// GetRules returns the rules list in a slice
func (l *RuleList) GetRules() []riddle.Rule {
	rules := make([]riddle.Rule, 0, len(l.rules))
	for _, rule := range l.rules {
		rules = append(rules, *rule)
	}
	return rules
}

// SetRules resets the list and inserts the provided rules
func (l *RuleList) SetRules(rules []riddle.Rule) {
	l.Reset()
	for _, rule := range rules {
		heapRule := &riddle.Rule{}
		*heapRule = rule
		l.addRule(heapRule, false)
	}
	l.Save()
}

// SaveRule adds a new rule to the list or updates an existing one
func (l *RuleList) SaveRule(rule *riddle.Rule) {
	l.addRule(rule, true)
}

func (l *RuleList) addRule(rule *riddle.Rule, save bool) {
	text := fmt.Sprintf("%s - %s - %s",
		colorizeItem(rule.ItemA),
		colorizeItem(rule.ItemB),
		rule.Relation.String())

	if rule.HasCondition() {
		text += fmt.Sprintf(" if A and B is %s and %s",
			colorize(rule.ConditionItemType),
			colorize(rule.Condition))
		if rule.IsReversible {
			text += " [reversible[]"
		}
	}

	selected := func() {
		if l.editFunc != nil {
			l.editFunc(rule)
		}
	}

	index, found := l.findRule(rule)
	if found {
		l.RemoveItem(index)
		l.InsertItem(index, text, "", 0, selected)
	} else {
		l.rules = append(l.rules, rule)
		l.InsertItem(-1, text, "", 0, selected)
	}

	if save {
		l.Save()
	}
}

func (l *RuleList) findRule(rule *riddle.Rule) (index int, found bool) {
	for i, r := range l.rules {
		if r == rule {
			return i, true
		}
	}
	return -1, false
}

func (l *RuleList) removeRule(index int, save bool) {
	if len(l.rules) == 0 {
		return
	}

	if index < 0 {
		index = len(l.rules) + index
	}
	if index >= len(l.rules) {
		index = len(l.rules) - 1
	}
	if index < 0 {
		index = 0
	}

	l.rules = append(l.rules[:index], l.rules[index+1:]...)
	l.RemoveItem(index)

	if save {
		l.Save()
	}
}

// Save calls the save function/callback with the list of rules
func (l *RuleList) Save() {
	if l.saveFunc != nil {
		l.saveFunc(l.GetRules())
	}
}

func (l *RuleList) handleInput(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	if key >= tcell.KeyDelete {
		remove := func() {
			l.removeRule(l.GetCurrentItem(), true)
		}
		l.modal.ModalYesNo("Do you really want to remove this rule?", remove)
		return nil
	}

	return event
}

// SetEditFunc sets a function that gets called on the selected rule
func (l *RuleList) SetEditFunc(editFunc func(*riddle.Rule)) {
	l.editFunc = editFunc
}

// SetSaveFunc sets a function that gets the list of all rules upon an update
func (l *RuleList) SetSaveFunc(saveFunc func([]riddle.Rule)) {
	l.saveFunc = saveFunc
}

// Reset resets the list
func (l *RuleList) Reset() {
	l.rules = nil
	l.Clear()
}

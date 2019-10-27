package ui

import (
	"fmt"

	"github.com/gdamore/tcell"
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
	l := &RuleList{
		List:  tview.NewList().ShowSecondaryText(false),
		modal: modal}

	l.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return l.handleInput(event)
	})

	return l
}

// HandleSetup filters the list based on the new setup
func (l *RuleList) HandleSetup(setup solver.Setup) {
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
		l.save()
	}
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

	l.rules = append(l.rules, rule)
	l.addRule(rule, -1)
}

func (l *RuleList) addRule(rule *solver.Rule, index int) {
	text := fmt.Sprintf("%s - %s - %s", rule.ItemA, rule.ItemB, rule.Relation.String())
	if len(rule.Condition) > 0 {
		text += fmt.Sprintf(" if [red][ %s ] %s", rule.ConditionItemType, rule.Condition)
	}

	selected := func() {
		if l.editFunc != nil {
			l.editFunc(rule)
		}
	}

	l.InsertItem(index, text, "", 0, selected)
	l.save()
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
		l.save()
	}
}

func (l *RuleList) save() {
	if l.saveFunc == nil {
		return
	}

	rules := make([]solver.Rule, 0, len(l.rules))
	for _, rule := range l.rules {
		rules = append(rules, *rule)
	}

	l.saveFunc(rules)
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
func (l *RuleList) SetEditFunc(editFunc func(*solver.Rule)) {
	l.editFunc = editFunc
}

// SetSaveFunc sets a function that gets the list of all rules upon an update
func (l *RuleList) SetSaveFunc(saveFunc func([]solver.Rule)) {
	l.saveFunc = saveFunc
}

// Reset resets the list
func (l *RuleList) Reset() {
	l.rules = nil
	l.Clear()
}

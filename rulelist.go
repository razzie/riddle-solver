package main

import (
	"github.com/rivo/tview"
)

// RuleList is a UI element that contains the list of rules
type RuleList struct {
	*tview.List
	rules    []*Rule
	saveFunc func([]Rule)
}

// NewRuleList returns a new RuleList
func NewRuleList() *RuleList {
	return &RuleList{List: tview.NewList()}
}

// HandleSetup filters the list based on the new setup
func (l *RuleList) HandleSetup(setup Setup) {

}

// SaveRule adds a new rule to the list or updates an existing one
func (l *RuleList) SaveRule(rule *Rule) {

}

// SetSaveFunc sets a function that gets the list of all rules upon an update
func (l *RuleList) SetSaveFunc(saveFunc func([]Rule)) {
	l.saveFunc = saveFunc
}

// Reset resets the list
func (l *RuleList) Reset() {

}

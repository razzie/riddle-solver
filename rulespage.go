package main

import (
	"github.com/rivo/tview"
)

// RulesPage is a form where the user can input riddle rules/hints
type RulesPage struct {
	*tview.Grid
	form     *RuleForm
	list     *RuleList
	saveFunc func([]Rule)
}

// NewRulesPage returns a new RulesPage
func NewRulesPage() *RulesPage {
	form := NewRuleForm()
	list := NewRuleList()
	grid := tview.NewGrid().
		SetRows(13, 0).
		SetColumns(0).
		AddItem(form, 0, 0, 1, 1, 0, 0, true).
		AddItem(list, 1, 0, 1, 1, 0, 0, false)

	return &RulesPage{
		Grid: grid,
		form: form,
		list: list}
}

// Reset resets the form and removes all rules from the list
func (p *RulesPage) Reset() {

}

// HandleSetup updates the form and list of rules based on the new setup
func (p *RulesPage) HandleSetup(setup Setup) {

}

// SetSaveFunc sets a function that handles the rules upon an update
func (p *RulesPage) SetSaveFunc(saveFunc func([]Rule)) {
	p.saveFunc = saveFunc
}

package ui

import (
	"github.com/razzie/riddle-solver/solver"
	"github.com/rivo/tview"
)

// RulesPage is a form where the user can input riddle rules/hints
type RulesPage struct {
	*tview.Grid
	form *RuleForm
	list *RuleList
}

// NewRulesPage returns a new RulesPage
func NewRulesPage(modal ModalHandler) *RulesPage {
	form := NewRuleForm(modal)
	list := NewRuleList(modal)
	grid := tview.NewGrid().
		SetRows(13, 0).
		SetColumns(0).
		AddItem(form, 0, 0, 1, 1, 0, 0, true).
		AddItem(list, 1, 0, 1, 1, 0, 0, false)

	form.SetSaveFunc(func(rule *solver.Rule) { list.SaveRule(rule) })
	list.SetEditFunc(func(rule *solver.Rule) { form.EditRule(rule) })

	return &RulesPage{
		Grid: grid,
		form: form,
		list: list}
}

// Reset resets the form and removes all rules from the list
func (p *RulesPage) Reset() {
	p.form.Reset()
	p.list.Reset()
}

// HandleSetup updates the form and list of rules based on the new setup
func (p *RulesPage) HandleSetup(setup solver.Setup) {
	p.form.HandleSetup(setup)
	p.list.HandleSetup(setup)
}

// SetSaveFunc sets a function that handles the rules upon an update
func (p *RulesPage) SetSaveFunc(saveFunc func([]solver.Rule)) {
	p.list.SetSaveFunc(saveFunc)
}

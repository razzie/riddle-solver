package ui

import (
	"github.com/razzie/riddle-solver/solver"
	"github.com/rivo/tview"
)

// RootElement is the root UI element
type RootElement struct {
	*PageHandler
	SetupForm *SetupForm
	RuleForm  *RuleForm
	RuleList  *RuleList
}

// NewRootElement returns a new RootElement
func NewRootElement() *RootElement {
	root := NewPageHandler()

	setup := NewSetupForm(root)
	root.AddPage("Setup", setup, nil)

	addRule := NewRuleForm(root)
	root.AddPage("Add rule", addRule, nil)

	rules := NewRuleList(root)
	root.AddPage("Rules", tview.NewFrame(rules), nil)

	results := NewResultsTree()
	root.AddPage("Results", tview.NewFrame(results), func() { results.Update() })

	setup.SetSaveFunc(func(setup solver.Setup) {
		addRule.HandleSetup(setup)
		rules.HandleSetup(setup)
		results.HandleSetup(setup)
		root.SwitchToPage(1)
	})
	addRule.SetSaveFunc(func(rule *solver.Rule) {
		rules.SaveRule(rule)
		root.ModalMessage("Saved")
	})
	rules.SetEditFunc(func(rule *solver.Rule) {
		addRule.EditRule(rule)
		root.SwitchToPage(1)
	})
	rules.SetSaveFunc(func(rules []solver.Rule) {
		results.HandleRules(rules)
	})

	return &RootElement{
		PageHandler: root,
		SetupForm:   setup,
		RuleForm:    addRule,
		RuleList:    rules}
}

package ui

import (
	"github.com/razzie/riddle-solver/riddle"
)

// RootElement is the root UI element
type RootElement struct {
	*PageHandler
	SetupForm *SetupPage
	RuleForm  *AddRulePage
	RuleList  *RulesPage
}

// NewRootElement returns a new RootElement
func NewRootElement(debug bool) *RootElement {
	if currentTheme == nil {
		LightTheme.Apply()
	}

	root := NewPageHandler()

	setup := NewSetupForm(root)
	root.AddPage(setup)

	addRule := NewRuleForm(root)
	root.AddPage(addRule)

	rules := NewRuleList(root)
	root.AddPage(rules)

	results := NewResultsTree(root)
	root.AddPage(results)

	solverdebug := NewSolverDebugTree(root)
	if debug {
		root.AddPage(solverdebug)
	}

	setup.SetSaveFunc(func(setup riddle.Setup) {
		addRule.HandleSetup(setup)
		rules.HandleSetup(setup)
		results.HandleSetup(setup)
		solverdebug.HandleSetup(setup)
		if setup != nil {
			root.SwitchToPage(1)
		}
	})
	addRule.SetSaveFunc(func(rule *riddle.Rule) {
		rules.SaveRule(rule)
		root.ModalMessage("Saved")
	})
	rules.SetEditFunc(func(rule *riddle.Rule) {
		addRule.EditRule(rule)
		root.SwitchToPage(1)
	})
	rules.SetSaveFunc(func(rules []riddle.Rule) {
		results.HandleRules(rules)
		solverdebug.HandleRules(rules)
	})

	return &RootElement{
		PageHandler: root,
		SetupForm:   setup,
		RuleForm:    addRule,
		RuleList:    rules,
	}
}

// GetRiddle returns the current riddle
func (root *RootElement) GetRiddle() (*riddle.Riddle, error) {
	rules := root.RuleList.GetRules()
	setup, err := root.SetupForm.GetSetup()
	if err != nil {
		return nil, err
	}

	return &riddle.Riddle{
		Setup: setup,
		Rules: rules,
	}, nil
}

// SetRiddle sets the current riddle
func (root *RootElement) SetRiddle(r *riddle.Riddle) {
	root.SetupForm.SetSetup(r.Setup)
	root.RuleList.SetRules(r.Rules)
	root.SwitchToPage(2)
}

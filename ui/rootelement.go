package ui

import (
	"github.com/razzie/riddle-solver/riddle"
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
func NewRootElement(debug bool) *RootElement {
	if currentTheme == nil {
		LightTheme.Apply()
	}

	root := NewPageHandler()

	setup := NewSetupForm(root)
	root.AddPage("Setup", setup, nil)

	addRule := NewRuleForm(root)
	root.AddPage("Add rule", addRule, nil)

	rules := NewRuleList(root)
	root.AddPage("Rules", tview.NewFrame(rules), nil)

	results := NewResultsTree()
	root.AddPage("Results", tview.NewFrame(results), func() { results.Update() })

	solverdebug := NewSolverDebugTree()
	if debug {
		root.AddPage("Debug", tview.NewFrame(solverdebug), func() { solverdebug.Update() })
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

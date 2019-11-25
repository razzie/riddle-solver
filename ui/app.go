package ui

import (
	"github.com/razzie/riddle-solver/riddle"
	"github.com/rivo/tview"
)

// App handles the user interface
type App struct {
	*PageHandler
	SetupForm *SetupPage
	RuleForm  *AddRulePage
	RuleList  *RulesPage
	app       *tview.Application
}

// NewApp returns a new App
func NewApp(debug bool) *App {
	if currentTheme == nil {
		LightTheme.Apply()
	}

	pages := NewPageHandler()

	setup := NewSetupForm(pages)
	pages.AddPage(setup)

	addRule := NewRuleForm(pages)
	pages.AddPage(addRule)

	rules := NewRuleList(pages)
	pages.AddPage(rules)

	results := NewResultsTree(pages)
	pages.AddPage(results)

	solverdebug := NewSolverDebugTree(pages)
	if debug {
		pages.AddPage(solverdebug)
	}

	setup.SetSaveFunc(func(setup riddle.Setup) {
		addRule.HandleSetup(setup)
		rules.HandleSetup(setup)
		results.HandleSetup(setup)
		solverdebug.HandleSetup(setup)
		if setup != nil {
			pages.SwitchToPage(1)
		}
	})
	addRule.SetSaveFunc(func(rule *riddle.Rule) {
		rules.SaveRule(rule)
		pages.ModalMessage("Saved")
	})
	rules.SetEditFunc(func(rule *riddle.Rule) {
		addRule.EditRule(rule)
		pages.SwitchToPage(1)
	})
	rules.SetSaveFunc(func(rules []riddle.Rule) {
		results.HandleRules(rules)
		solverdebug.HandleRules(rules)
	})

	app := tview.NewApplication().
		SetInputCapture(pages.InputCapture()).
		SetRoot(pages, true)

	return &App{
		PageHandler: pages,
		SetupForm:   setup,
		RuleForm:    addRule,
		RuleList:    rules,
		app:         app,
	}
}

// Run runs the user interface
func (app *App) Run() error {
	go func() {
		<-app.Quit
		app.app.Stop()
	}()

	SetConsoleTitle("Razzie's Riddle Solver")

	return app.app.SetFocus(app).Run()
}

// GetRiddle returns the current riddle
func (app *App) GetRiddle() (*riddle.Riddle, error) {
	rules := app.RuleList.GetRules()
	setup, err := app.SetupForm.GetSetup()
	if err != nil {
		return nil, err
	}

	return &riddle.Riddle{
		Setup: setup,
		Rules: rules,
	}, nil
}

// SetRiddle sets the current riddle
func (app *App) SetRiddle(r *riddle.Riddle) error {
	if err := r.Check(); err != nil {
		return err
	}

	app.SetupForm.SetSetup(r.Setup)
	app.RuleList.SetRules(r.Rules)
	app.SwitchToPage(2)
	return nil
}

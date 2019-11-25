package main

import (
	"github.com/razzie/riddle-solver/riddle"
	"github.com/razzie/riddle-solver/ui"
	"github.com/rivo/tview"
)

// App handles the user interface
type App struct {
	*ui.PageHandler
	SetupForm *ui.SetupPage
	RuleForm  *ui.AddRulePage
	RuleList  *ui.RulesPage
	app       *tview.Application
}

// NewApp returns a new App
func NewApp(debug bool) *App {
	pages := ui.NewPageHandler()

	setup := ui.NewSetupForm(pages)
	pages.AddPage(setup)

	addRule := ui.NewRuleForm(pages)
	pages.AddPage(addRule)

	rules := ui.NewRuleList(pages)
	pages.AddPage(rules)

	results := ui.NewResultsTree(pages)
	pages.AddPage(results)

	load := ui.NewLoadPage(pages)
	pages.AddPage(load)

	save := ui.NewSavePage(pages)
	pages.AddPage(save)

	solverdebug := ui.NewSolverDebugTree(pages)
	if debug {
		pages.AddPage(solverdebug)
	}

	setup.SetSaveFunc(func(setup riddle.Setup) {
		addRule.HandleSetup(setup)
		rules.HandleSetup(setup)
		results.HandleSetup(setup)
		solverdebug.HandleSetup(setup)
		if len(setup) > 0 {
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

	subapp := tview.NewApplication().
		SetInputCapture(pages.InputCapture()).
		SetRoot(pages, true)

	app := &App{
		PageHandler: pages,
		SetupForm:   setup,
		RuleForm:    addRule,
		RuleList:    rules,
		app:         subapp,
	}
	load.SetRiddleSetter(app.SetRiddle)
	save.SetRiddleGetter(app.GetRiddle)

	r, err := riddle.LoadRiddleFromFile("riddles/autosave.json")
	if err == nil {
		app.SetRiddle(r)
	}

	return app
}

// Run runs the user interface
func (app *App) Run() error {
	go func() {
		<-app.Quit
		app.autosave()
		app.app.Stop()
	}()

	ui.SetConsoleTitle("Razzie's Riddle Solver")

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

	app.autosave()

	app.SetupForm.SetSetup(r.Setup)
	app.RuleList.SetRules(r.Rules)
	if len(r.Rules) > 0 {
		app.SwitchToPage(2)
	} else {
		app.SwitchToPage(0)
	}
	return nil
}

func (app *App) autosave() {
	r, err := app.GetRiddle()
	if err == nil && len(r.Setup) > 0 {
		r.SaveToFile("riddles/autosave.json")
	}
}

package gui

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	gioapp "gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	"github.com/razzie/riddle-solver/pkg/riddle"
)

type App struct {
	theme   *material.Theme
	pages   *PageHandler
	setup   *SetupPage
	addRule *AddRulePage
	rules   *RulesPage
}

func NewApp(th *material.Theme, debug bool) *App {
	pages := NewPageHandler(th)

	setup := NewSetupPage(th, pages)
	pages.AddPage(setup)

	addRule := NewAddRulePage(th, pages)
	addRulePageIdx := pages.AddPage(addRule)

	rules := NewRulesPage(th, pages)
	rulesPageIdx := pages.AddPage(rules)

	results := NewResultsPage(th, pages)
	pages.AddPage(results)

	solverdebug := NewDebugPage(th, pages)
	if debug {
		pages.AddPage(solverdebug)
	}

	setup.SetSaveFunc(func(setup riddle.Setup) {
		addRule.HandleSetup(setup)
		rules.HandleSetup(setup)
		results.HandleSetup(setup)
		solverdebug.HandleSetup(setup)
		if len(setup) > 0 {
			pages.SwitchToPage(rulesPageIdx)
		}
	})
	addRule.SetSaveFunc(func(rule *riddle.Rule) {
		rules.SaveRule(rule)
		pages.ModalMessage("Saved")
		pages.SwitchToPage(rulesPageIdx)
	})
	rules.SetEditFunc(func(rule *riddle.Rule) {
		addRule.EditRule(rule)
		pages.SwitchToPage(addRulePageIdx)
	})
	rules.SetSaveFunc(func(rules []riddle.Rule) {
		results.HandleRules(rules)
		solverdebug.HandleRules(rules)
	})

	app := &App{
		theme:   th,
		pages:   pages,
		setup:   setup,
		addRule: addRule,
		rules:   rules,
	}
	return app
}

func (app *App) Run() error {
	go func() {
		defer os.Exit(0)
		w := gioapp.NewWindow(gioapp.Title("Razzie's Riddle Solver"))
		if err := app.loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	gioapp.Main()
	return nil
}

func (app *App) SetRiddle(r *riddle.Riddle) error {
	if err := r.Check(); err != nil {
		return err
	}

	//app.autosave()

	app.setup.SetSetup(r.Setup)
	app.rules.SetRules(r.Rules)
	return nil
}

func (app *App) loop(w *gioapp.Window) error {
	defer func() {
		if r := recover(); r != nil {
			OSMessageBox(fmt.Sprint(r, "\n", string(debug.Stack())), "Error")
			os.Exit(1)
		}
	}()

	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			app.pages.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

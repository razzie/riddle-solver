package gui

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	gioapp "gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	"github.com/gen2brain/dlgs"
	"github.com/razzie/razgio"
	"github.com/razzie/riddle-solver/pkg/riddle"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type App struct {
	theme   *material.Theme
	w       *gioapp.Window
	pages   *razgio.PageHandler
	setup   *SetupPage
	addRule *AddRulePage
	rules   *RulesPage
}

func NewApp(th *material.Theme, debug bool) *App {
	pages := razgio.NewPageHandler(th)

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

	load := NewLoadPage(th, pages)
	pages.AddPage(load)

	save := NewSavePage(th, pages)
	pages.AddPage(save)

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
	load.SetRiddleSetter(app.SetRiddle)
	save.SetRiddleGetter(app.GetRiddle)
	return app
}

func (app *App) Run() error {
	go func() {
		defer os.Exit(0)
		app.w = gioapp.NewWindow(gioapp.Title("Razzie's Riddle Solver"))
		if err := app.loop(); err != nil {
			log.Fatal(err)
		}
	}()
	gioapp.Main()
	return nil
}

func (app *App) GetRiddle() (*riddle.Riddle, error) {
	rules := app.rules.GetRules()
	setup, err := app.setup.GetSetup()
	if err != nil {
		return nil, err
	}

	return &riddle.Riddle{
		Setup: setup,
		Rules: rules,
	}, nil
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

func (app *App) loop() error {
	defer func() {
		if r := recover(); r != nil {
			dlgs.Error("Error", fmt.Sprint(r, "\n", string(debug.Stack())))
			os.Exit(1)
		}
	}()

	var ops op.Ops
	for {
		e := <-app.w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case key.Event:
			if e.Name == key.NameEscape {
				app.pages.ModalYesNo("Exit program?", app.w.Close)
			}
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			//key.InputOp{Tag: app, Hint: key.HintAny}.Add(gtx.Ops)
			//key.FocusOp{Tag: app}.Add(gtx.Ops)
			app.pages.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

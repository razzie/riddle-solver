package gui

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	"github.com/razzie/riddle-solver/pkg/riddle"
)

type App struct {
	th *material.Theme
	ph *PageHandler
}

func NewApp(th *material.Theme, debug bool) *App {
	a := &App{
		th: th,
		ph: NewPageHandler(th),
	}

	r := riddle.NewEinsteinRiddle()

	//a.ph.tabs.AddTab("test tab", nil)

	setuppage := NewSetupPage(th, a.ph)
	setuppage.SetSetup(r.Setup)
	a.ph.AddPage(setuppage)

	addrulepage := NewAddRulePage(th, a.ph)
	a.ph.AddPage(addrulepage)

	rulespage := NewRulesPage(th, a.ph)
	rulespage.SetRules(r.Rules)
	a.ph.AddPage(rulespage)

	return a
}

func (a *App) Run() error {
	go func() {
		defer os.Exit(0)
		w := app.NewWindow(app.Title("Razzie's Riddle Solver"))
		if err := a.loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
	return nil
}

func (a *App) SetRiddle(r *riddle.Riddle) error {
	return nil
}

func (a *App) loop(w *app.Window) error {
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
			a.ph.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

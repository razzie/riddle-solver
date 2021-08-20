package gui

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	"github.com/razzie/riddle-solver/pkg/riddle"
)

type App struct {
	th   *material.Theme
	tabs Tabs
}

func NewApp(th *material.Theme, debug bool) *App {
	a := &App{
		th:   th,
		tabs: Tabs{Theme: th},
	}

	for i := 1; i <= 10; i++ {
		a.tabs.tabs = append(a.tabs.tabs,
			Tab{Title: fmt.Sprintf("Tab %d", i)},
		)
	}

	return a
}

func (a *App) Run() error {
	go func() {
		defer os.Exit(0)
		w := app.NewWindow()
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
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			a.tabs.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

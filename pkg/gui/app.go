package gui

import (
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

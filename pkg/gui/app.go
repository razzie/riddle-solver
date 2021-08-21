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
	th *material.Theme
	ph *PageHandler
}

func NewApp(th *material.Theme, debug bool) *App {
	a := &App{
		th: th,
		ph: NewPageHandler(th),
	}

	for i := 1; i <= 10; i++ {
		a.ph.tabs.AddTab(fmt.Sprintf("Tab %d", i), nil)
	}
	a.ph.tabs.SetSelectFunc(func(i int) {
		if i == 3 {
			a.ph.ModalYesNo("test test test", func() {
				a.ph.ModalMessage("yes pressed")
			})
		}
	})

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

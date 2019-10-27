package main

import (
	"github.com/razzie/riddle-solver/ui"
	"github.com/rivo/tview"
)

func main() {
	root := ui.NewRootElement()
	app := tview.NewApplication().
		SetInputCapture(root.InputCapture()).
		SetRoot(root, true)

	go func() {
		<-root.Quit
		app.Stop()
	}()

	if err := app.SetFocus(root).Run(); err != nil {
		panic(err)
	}
}

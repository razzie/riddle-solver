package main

import (
	"github.com/rivo/tview"
)

func main() {
	root := NewPageHandler()
	app := tview.NewApplication().
		SetInputCapture(root.InputCapture()).
		SetRoot(root, true)

	setup := NewSetupPage()
	root.AddPage("setup", setup)

	if err := app.SetFocus(root).Run(); err != nil {
		panic(err)
	}
}

package main

import (
	"github.com/rivo/tview"
)

func main() {
	root := NewPageHandler()
	app := tview.NewApplication().
		SetInputCapture(root.InputCapture()).
		SetRoot(root, true)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

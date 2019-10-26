package main

import (
	"fmt"

	"github.com/rivo/tview"
)

func main() {
	root := NewPageHandler()
	app := tview.NewApplication().
		SetInputCapture(root.InputCapture()).
		SetRoot(root, true)

	go func() {
		<-root.Quit
		app.Stop()
	}()

	setup := NewSetupPage()
	setup.SetSaveFunc(func(setup Setup) {
		if err := setup.Check(); err != nil {
			root.ModalMessage(fmt.Sprint(err))
		}
	})
	root.AddPage("Setup", setup)

	rules := NewRulesPage()
	root.AddPage("Rules", rules)

	if err := app.SetFocus(root).Run(); err != nil {
		panic(err)
	}
}

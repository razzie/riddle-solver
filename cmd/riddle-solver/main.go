package main

import (
	"github.com/razzie/riddle-solver/solver"
	"github.com/razzie/riddle-solver/ui"
	"github.com/rivo/tview"
)

func main() {
	root := ui.NewPageHandler()
	app := tview.NewApplication().
		SetInputCapture(root.InputCapture()).
		SetRoot(root, true)

	go func() {
		<-root.Quit
		app.Stop()
	}()

	setup := ui.NewSetupPage(root)
	root.AddPage("Setup", setup)

	rules := ui.NewRulesPage(root)
	root.AddPage("Rules", rules)
	setup.SetSaveFunc(func(setup solver.Setup) {
		rules.HandleSetup(setup)
		root.SwitchToPage(1)
	})

	if err := app.SetFocus(root).Run(); err != nil {
		panic(err)
	}
}

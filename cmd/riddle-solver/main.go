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

	setup := ui.NewSetupForm(root)
	root.AddPage("Setup", setup)

	addRule := ui.NewRuleForm(root)
	root.AddPage("Add rule", addRule)

	rules := ui.NewRuleList(root)
	root.AddPage("Rules", tview.NewFrame(rules))

	setup.SetSaveFunc(func(setup solver.Setup) {
		addRule.HandleSetup(setup)
		rules.HandleSetup(setup)
		root.SwitchToPage(1)
	})
	addRule.SetSaveFunc(func(rule *solver.Rule) {
		rules.SaveRule(rule)
		root.ModalMessage("Saved")
	})
	rules.SetEditFunc(func(rule *solver.Rule) {
		addRule.EditRule(rule)
		root.SwitchToPage(1)
	})

	if err := app.SetFocus(root).Run(); err != nil {
		panic(err)
	}
}

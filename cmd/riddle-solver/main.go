package main

import (
	"flag"

	"github.com/razzie/riddle-solver/ui"
	"github.com/rivo/tview"
)

func main() {
	demo := flag.Bool("demo", false, "Einstein's 5 house riddle demo mode")
	flag.Parse()

	root := ui.NewRootElement()
	app := tview.NewApplication().
		SetInputCapture(root.InputCapture()).
		SetRoot(root, true)

	if *demo {
		SetupDemo(root)
	}

	go func() {
		<-root.Quit
		app.Stop()
	}()

	if err := app.SetFocus(root).Run(); err != nil {
		panic(err)
	}
}

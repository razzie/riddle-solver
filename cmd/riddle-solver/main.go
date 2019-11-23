package main

import (
	"flag"

	"github.com/razzie/riddle-solver/riddle"
	"github.com/razzie/riddle-solver/ui"
	"github.com/rivo/tview"
)

func main() {
	ui.SetConsoleTitle("Razzie's Riddle Solver")

	demo := flag.Bool("demo", false, "Einstein's 5 house riddle demo mode")
	theme := flag.String("theme", "light", "Specify light or dark theme")
	debug := flag.Bool("debug", false, "Enable an additional debug page")
	load := flag.String("load", "riddle.json", "Specify a riddle JSON file to load")
	flag.Parse()

	themes := map[string]*ui.Theme{
		"light": &ui.LightTheme,
		"dark":  &ui.DarkTheme,
	}

	if t, found := themes[*theme]; found {
		t.Apply()
	}

	root := ui.NewRootElement(*debug)
	app := tview.NewApplication().
		SetInputCapture(root.InputCapture()).
		SetRoot(root, true)

	if *demo {
		root.SetRiddle(riddle.EinsteinRiddle)
	} else {
		r, err := riddle.LoadRiddleFromFile(*load)
		if err == nil {
			root.SetRiddle(r)
		}
	}

	go func() {
		<-root.Quit
		r, err := root.GetRiddle()
		if err == nil {
			r.SaveToFile("riddle.json")
		}
		app.Stop()
	}()

	if err := app.SetFocus(root).Run(); err != nil {
		panic(err)
	}
}

package main

import (
	"flag"

	"github.com/razzie/riddle-solver/riddle"
	"github.com/razzie/riddle-solver/ui"
)

func main() {
	ui.SetConsoleTitle("Razzie's Riddle Solver")

	demo := flag.Bool("demo", false, "Einstein's 5 house riddle demo mode")
	theme := flag.String("theme", "light", "Specify light or dark theme")
	debug := flag.Bool("debug", false, "Enable an additional debug page")
	load := flag.String("load", "autosave.json", "Specify a riddle JSON file to load")
	flag.Parse()

	themes := map[string]*ui.Theme{
		"light": &ui.LightTheme,
		"dark":  &ui.DarkTheme,
	}

	if t, found := themes[*theme]; found {
		t.Apply()
	}

	app := ui.NewApp(*debug)

	if *demo {
		app.SetRiddle(riddle.EinsteinRiddle)
	} else {
		r, err := riddle.LoadRiddleFromFile(*load)
		if err == nil {
			app.SetRiddle(r)
		}
	}

	if err := app.Run(); err != nil {
		panic(err)
	}

	r, err := app.GetRiddle()
	if err == nil {
		r.SaveToFile("autosave.json")
	}
}

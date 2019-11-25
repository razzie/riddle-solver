package main

import (
	"flag"
	"fmt"

	"github.com/razzie/riddle-solver/riddle"
	"github.com/razzie/riddle-solver/ui"
)

func main() {
	demo := flag.Bool("demo", false, "Einstein's 5 house riddle demo mode")
	theme := flag.String("theme", "light", "Specify light or dark theme")
	debug := flag.Bool("debug", false, "Enable an additional debug page")
	load := flag.String("load", "", "Specify a riddle JSON file to load")
	flag.Parse()

	if t, ok := ui.Themes[*theme]; ok {
		t.Apply()
	} else {
		panic(fmt.Errorf("Theme not found: %s", *theme))
	}

	app := NewApp(*debug)

	if *demo {
		app.SetRiddle(riddle.NewEinsteinRiddle())
	} else if len(*load) > 0 {
		r, err := riddle.LoadRiddleFromFile(*load)
		if err == nil {
			app.SetRiddle(r)
		}
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}

package main

import (
	"flag"
	"log"

	"github.com/razzie/riddle-solver/pkg/riddle"
)

var (
	theme      string
	riddleFile string
	debug      bool
	nogui      bool
)

func init() {
	flag.StringVar(&theme, "theme", "light", "Specify light or dark theme")
	flag.StringVar(&riddleFile, "load", "", "Specify a riddle JSON file to load")
	flag.BoolVar(&debug, "debug", false, "Enable an additional debug page")
	flag.BoolVar(&nogui, "nogui", false, "Disable GUI and use terminal UI instead")
	flag.Parse()
}

func tryLoadRiddle() *riddle.Riddle {
	if len(riddleFile) > 0 {
		r, err := riddle.LoadRiddleFromFile(riddleFile)
		if err != nil {
			log.Fatalf("failed to load riddle: %v", err)
		}
		return r
	}
	return nil
}

func main() {
	var app App
	r := tryLoadRiddle()

	if nogui {
		app = getTuiApp(theme, debug)
	} else {
		app = getGuiApp(theme, debug)
	}

	if r != nil {
		app.SetRiddle(r)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

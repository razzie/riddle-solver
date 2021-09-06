package main

import (
	"flag"
	"log"

	"github.com/razzie/riddle-solver/pkg/gui"
	"github.com/razzie/riddle-solver/pkg/riddle"
)

var (
	theme      string
	riddleFile string
	debug      bool
)

func init() {
	flag.StringVar(&theme, "theme", "light", "Specify light or dark theme")
	flag.StringVar(&riddleFile, "load", "", "Specify a riddle JSON file to load")
	flag.BoolVar(&debug, "debug", false, "Enable an additional debug page")
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
	return riddle.NewEinsteinRiddle()
}

func main() {
	t, ok := gui.Themes[theme]
	if !ok {
		log.Fatalf("Theme not found: %s", theme)
	}

	app := gui.NewApp(t, debug)
	if r := tryLoadRiddle(); r != nil {
		app.SetRiddle(r)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

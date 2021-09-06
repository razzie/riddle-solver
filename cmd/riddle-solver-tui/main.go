package main

import (
	"flag"
	"log"

	"github.com/razzie/riddle-solver/pkg/riddle"
	"github.com/razzie/riddle-solver/pkg/tui"
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
	return nil
}

func main() {
	if t, ok := tui.Themes[theme]; ok {
		t.Apply()
	} else {
		log.Fatalf("Theme not found: %s", theme)
	}

	app := tui.NewApp(debug)
	if r := tryLoadRiddle(); r != nil {
		app.SetRiddle(r)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

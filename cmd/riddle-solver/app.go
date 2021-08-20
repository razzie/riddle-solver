package main

import (
	"log"

	"github.com/razzie/riddle-solver/pkg/gui"
	"github.com/razzie/riddle-solver/pkg/riddle"
	"github.com/razzie/riddle-solver/pkg/tui"
)

type App interface {
	Run() error
	SetRiddle(r *riddle.Riddle) error
}

func getTuiApp(theme string, debug bool) App {
	if t, ok := tui.Themes[theme]; ok {
		t.Apply()
	} else {
		log.Fatalf("Theme not found: %s", theme)
	}

	return tui.NewApp(debug)
}

func getGuiApp(theme string, debug bool) App {
	t, ok := gui.Themes[theme]
	if !ok {
		log.Fatalf("Theme not found: %s", theme)
	}

	return gui.NewApp(t, debug)
}

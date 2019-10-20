package main

import (
	//	"github.com/gdamore/tcell/"
	"github.com/rivo/tview"
)

func main() {
	box := tview.NewBox().SetBorder(true).SetTitle(" Razzie's Riddle Solver ")
	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}

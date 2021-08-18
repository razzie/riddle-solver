package tui

import (
	"path/filepath"

	"github.com/razzie/riddle-solver/pkg/riddle"
	"github.com/rivo/tview"
)

// LoadPage is a UI list element that contains loadable riddles
type LoadPage struct {
	Page
	list   *tview.List
	setter func(*riddle.Riddle) error
	modal  ModalHandler
}

// NewLoadPage returns a new LoadPage
func NewLoadPage(modal ModalHandler) *LoadPage {
	list := tview.NewList().ShowSecondaryText(false)
	p := &LoadPage{
		Page:  NewPage(tview.NewFrame(list), "Load"),
		list:  list,
		modal: modal,
	}
	p.Page.SetSelectFunc(p.Reset)
	return p
}

// AddRiddle adds a riddle to the list
func (p *LoadPage) AddRiddle(name string, loadFunc func() (*riddle.Riddle, error)) {
	p.list.AddItem(name, "", 0, func() {
		if p.setter != nil {
			r, err := loadFunc()
			if err != nil {
				p.modal.ModalMessage(err.Error())
				return
			}

			err = p.setter(r)
			if err != nil {
				p.modal.ModalMessage(err.Error())
				return
			}

			p.modal.ModalMessage("Riddle loaded")
		}
	})
}

// SetRiddleSetter sets the function that gets called when a riddle is selected
func (p *LoadPage) SetRiddleSetter(setter func(*riddle.Riddle) error) {
	p.setter = setter
}

// Reset resets and updates the page
func (p *LoadPage) Reset() {
	p.list.Clear()

	p.AddRiddle("Einstein's 5 house riddle [built-in[]", func() (*riddle.Riddle, error) {
		return riddle.NewEinsteinRiddle(), nil
	})

	p.AddRiddle("The Jindosh riddle [built-in[]", func() (*riddle.Riddle, error) {
		return riddle.NewJindoshRiddle(), nil
	})

	files, _ := filepath.Glob("riddles/*.json")
	for _, file := range files {
		filename := file
		p.AddRiddle(filename, func() (*riddle.Riddle, error) {
			return riddle.LoadRiddleFromFile(filename)
		})
	}
}

package tui

import (
	"github.com/rivo/tview"
)

// Page represents a selectable page handled by PageHandler
type Page interface {
	tview.Primitive
	GetName() string
	SetSelectFunc(func())
	Select()
}

type page struct {
	tview.Primitive
	name       string
	selectFunc func()
}

// NewPage returns a new Page
func NewPage(primitive tview.Primitive, name string) Page {
	return &page{
		Primitive: primitive,
		name:      name,
	}
}

// GetName returns the page's name
func (p *page) GetName() string {
	return p.name
}

// SetSelectFunc sets a callback which is called when the page is selected
func (p *page) SetSelectFunc(selectFunc func()) {
	p.selectFunc = selectFunc
}

// Select is called by PageHandler when the page is selected
func (p *page) Select() {
	if p.selectFunc != nil {
		p.selectFunc()
	}
}

package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// PageHandler is the root UI element of the application
type PageHandler struct {
	grid       *tview.Grid
	pages      *tview.Pages
	footer     *tview.TextView
	pageNames  []string
	activePage int
}

// NewPageHandler returns a new PageHandler
func NewPageHandler() *PageHandler {
	pages := tview.NewPages()
	footer := tview.NewTextView().SetDynamicColors(true)
	grid := tview.NewGrid().
		SetRows(0, 1).
		SetColumns(0).
		AddItem(pages, 0, 0, 1, 1, 0, 0, true).
		AddItem(footer, 1, 0, 1, 1, 1, 0, false)

	return &PageHandler{
		grid:   grid,
		pages:  pages,
		footer: footer}
}

// AddPage adds a publicly listed page to the frame
func (ph *PageHandler) AddPage(name string, page tview.Primitive) *PageHandler {
	ph.pages.AddPage(name, page, true, ph.pages.GetPageCount() == 0)
	ph.pageNames = append(ph.pageNames, name)
	ph.updateFooter()
	return ph
}

// SwitchToPage switches to the page with the number 'page'
func (ph *PageHandler) SwitchToPage(page int) {
	if page < len(ph.pageNames) {
		ph.pages.SwitchToPage(ph.pageNames[page])
		ph.activePage = page
		ph.updateFooter()
	}
}

// InputCapture returns a delegate that handles input capture for PageHandler
func (ph *PageHandler) InputCapture() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		return ph.handleInput(event)
	}
}

func (ph *PageHandler) handleInput(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	if key >= tcell.KeyF1 && key <= tcell.KeyF12 {
		page := int(key - tcell.KeyF1)
		ph.SwitchToPage(page)
		return nil
	}

	return event
}

func (ph *PageHandler) updateFooter() {
	var footerText string

	for i, name := range ph.pageNames {
		if i == ph.activePage {
			footerText = fmt.Sprintf("%s [[red]F%d %s[white]] ", footerText, i+1, name)
		} else {
			footerText = fmt.Sprintf("%s [F%d %s] ", footerText, i+1, name)
		}
	}

	ph.footer.SetText(footerText)
}

// Draw implements tview.Primitive.Draw
func (ph *PageHandler) Draw(screen tcell.Screen) {
	ph.grid.Draw(screen)
}

// GetRect implements tview.Primitive.GetRect
func (ph *PageHandler) GetRect() (int, int, int, int) {
	return ph.grid.GetRect()
}

// SetRect implements tview.Primitive.SetRect
func (ph *PageHandler) SetRect(x, y, width, height int) {
	ph.grid.SetRect(x, y, width, height)
}

// InputHandler implements tview.Primitive.InputHandler
func (ph *PageHandler) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return ph.pages.InputHandler()
}

// Focus implements tview.Primitive.Focus
func (ph *PageHandler) Focus(delegate func(p tview.Primitive)) {
	ph.pages.Focus(delegate)
}

// Blur implements tview.Primitive.Blur
func (ph *PageHandler) Blur() {
	ph.pages.Blur()
}

// GetFocusable implements tview.Primitive.GetFocusable
func (ph *PageHandler) GetFocusable() tview.Focusable {
	return ph.pages.GetFocusable()
}

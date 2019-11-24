package ui

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// PageHandler handles the layout of the application, pages and modal dialogs
type PageHandler struct {
	tview.Primitive
	Quit        chan bool
	pageHandler *tview.Pages
	pages       []Page
	activePage  int
	footer      *tview.TextView
	modalMsg    *tview.Modal
	modalYesNo  *tview.Modal
	modalActive bool
}

// NewPageHandler returns a new PageHandler
func NewPageHandler() *PageHandler {
	pages := tview.NewPages()
	footer := tview.NewTextView().SetRegions(true)
	grid := tview.NewGrid().
		SetRows(0, 1).
		SetColumns(0).
		AddItem(pages, 0, 0, 1, 1, 0, 0, true).
		AddItem(footer, 1, 0, 1, 1, 1, 0, false)

	modalMsg := tview.NewModal().
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.HidePage("modal_msg")
		})
	pages.AddPage("modal_msg", modalMsg, false, false)

	modalYesNo := tview.NewModal().
		AddButtons([]string{"Yes", "No"})
	pages.AddPage("modal_yes_no", modalYesNo, false, false)

	return &PageHandler{
		Primitive:   grid,
		Quit:        make(chan bool),
		pageHandler: pages,
		footer:      footer,
		modalMsg:    modalMsg,
		modalYesNo:  modalYesNo,
	}
}

// AddPage adds a publicly listed page to the frame
func (ph *PageHandler) AddPage(page Page) *PageHandler {
	ph.pageHandler.AddPage(page.GetName(), page, true, len(ph.pages) == 0)
	ph.pages = append(ph.pages, page)
	ph.updateFooter()
	return ph
}

// SwitchToPage switches to the page with the number 'page'
func (ph *PageHandler) SwitchToPage(page int) {
	if ph.modalActive || ph.activePage == page {
		return
	}

	if page < len(ph.pages) {
		ph.pageHandler.SwitchToPage(ph.pages[page].GetName())
		ph.activePage = page
		ph.updateFooter()
		ph.pages[page].Select()
	}
}

// ModalMessage displays a modal window with a message and OK button
func (ph *PageHandler) ModalMessage(msg string) {
	ph.modalMsg.SetText(msg)
	ph.pageHandler.SendToFront("modal_msg").ShowPage("modal_msg")
	ph.modalActive = true
}

// ModalYesNo displays a modal dialog with a message and yes/no options
func (ph *PageHandler) ModalYesNo(msg string, yes func()) {
	ph.modalYesNo.SetText(msg).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonIndex == 0 {
			yes()
		}
		ph.pageHandler.HidePage("modal_yes_no")
		ph.modalActive = false
	})
	ph.pageHandler.SendToFront("modal_yes_no").ShowPage("modal_yes_no")
	ph.modalActive = true
}

// InputCapture returns a function that handles input capture for PageHandler
func (ph *PageHandler) InputCapture() func(event *tcell.EventKey) *tcell.EventKey {
	return ph.handleInput
}

func (ph *PageHandler) handleInput(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	if key >= tcell.KeyF1 && key <= tcell.KeyF12 {
		page := int(key - tcell.KeyF1)
		ph.SwitchToPage(page)
		return nil

	} else if key == tcell.KeyEscape {
		if ph.modalActive {
			ph.pageHandler.HidePage("modal_msg").HidePage("modal_yes_no")
			ph.modalActive = false
		} else {
			ph.ModalYesNo("Do you really want to quit?", func() {
				ph.Quit <- true
			})
		}
		return nil
	}

	return event
}

func (ph *PageHandler) updateFooter() {
	var footerText string

	for i, page := range ph.pages {
		name := page.GetName()
		footerText += fmt.Sprintf(" [\"%s\"] F%d %s [\"\"] ", name, i+1, name)
	}

	footerText += " ESC Quit"

	ph.footer.SetText(footerText)
	ph.footer.Highlight(ph.pages[ph.activePage].GetName())
}

package gui

import (
	"time"

	"gioui.org/layout"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

type PageHandler struct {
	theme *material.Theme
	tabs  *Tabs
	pages []Page
	modal component.ModalStyle
}

func NewPageHandler(th *material.Theme) *PageHandler {
	ph := &PageHandler{
		theme: th,
		tabs:  NewTabs(th),
		modal: component.Modal(th, new(component.ModalState)),
	}
	ph.modal.VisibilityAnimation.Duration = time.Millisecond * 250
	ph.modal.VisibilityAnimation.State = component.Invisible
	return ph
}

func (ph *PageHandler) AddPage(page Page) {
	ph.pages = append(ph.pages, page)
	ph.tabs.AddTab(page.GetName(), page.Layout)
}

func (ph *PageHandler) ModalMessage(msg string) {
	ph.modal.Show(time.Now(), NewMessageBox(ph.theme, msg, func() {
		ph.modal.Disappear(time.Now())
	}))
}

func (ph *PageHandler) ModalYesNo(msg string, yesFunc func()) {
	ph.modal.Show(time.Now(), NewYesNoMessageBox(ph.theme, msg, func(yes bool) {
		if yes {
			yesFunc()
		} else {
			ph.modal.Disappear(time.Now())
		}
	}))
}

func (ph *PageHandler) Layout(gtx C) D {
	gtx.Constraints.Min = gtx.Constraints.Max
	return layout.Stack{Alignment: layout.S}.Layout(gtx,
		layout.Stacked(ph.tabs.Layout),
		layout.Stacked(func(gtx C) D {
			return ph.modal.Layout(gtx)
		}),
	)
}

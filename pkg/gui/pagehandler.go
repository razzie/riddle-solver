package gui

import (
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

type PageHandler struct {
	theme *material.Theme
	tabs  *Tabs
	pages []Page
	modal component.ModalLayer
}

func NewPageHandler(th *material.Theme) *PageHandler {
	ph := &PageHandler{
		theme: th,
		tabs:  NewTabs(th),
		modal: *component.NewModal(),
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
	mbox := NewMessageBox(ph.theme, msg, func() {
		ph.modal.Disappear(time.Now())
	})
	ph.modal.Widget = func(gtx C, th *material.Theme, anim *component.VisibilityAnimation) D {
		return mbox(gtx)
	}
	ph.modal.Appear(time.Now())
}

func (ph *PageHandler) ModalYesNo(msg string, yesFunc func()) {
	mbox := NewYesNoMessageBox(ph.theme, msg, func(yes bool) {
		if yes {
			yesFunc()
		}
		ph.modal.Disappear(time.Now())
	})
	ph.modal.Widget = func(gtx C, th *material.Theme, anim *component.VisibilityAnimation) D {
		return mbox(gtx)
	}
	ph.modal.Appear(time.Now())
}

func (ph *PageHandler) Layout(gtx C) D {
	defer op.Save(gtx.Ops).Load()
	gtx.Constraints.Min = gtx.Constraints.Max
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return ph.tabs.Layout(gtx)
		}),
		layout.Expanded(func(gtx C) D {
			return ph.modal.Layout(gtx, ph.theme)
		}),
	)
}

package gui

import (
	"fmt"
	"image"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/razzie/razgio"
	"github.com/razzie/riddle-solver/pkg/riddle"
)

type SavePage struct {
	theme   *material.Theme
	modal   razgio.ModalHandler
	name    razgio.TextField
	saveBtn widget.Clickable
	getter  func() (*riddle.Riddle, error)
}

func NewSavePage(th *material.Theme, modal razgio.ModalHandler) *SavePage {
	return &SavePage{
		theme: th,
		modal: modal,
	}
}

func (p *SavePage) GetName() string {
	return "Save"
}

func (p *SavePage) Select() {

}

func (p *SavePage) Layout(gtx C) D {
	if p.saveBtn.Clicked() {
		p.Save()
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			dims := p.name.Layout(gtx, p.theme, "Riddle name", "riddle name without .json")
			dims.Size.Y += gtx.Px(unit.Dp(12))
			return dims
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min = image.Pt(0, 0)
			return material.Button(p.theme, &p.saveBtn, "Save").Layout(gtx)
		}),
	)
}

func (p *SavePage) SetRiddleGetter(getter func() (*riddle.Riddle, error)) {
	p.getter = getter
}

func (p *SavePage) Save() {
	name := p.name.Text()
	if len(name) == 0 {
		p.modal.ModalMessage("Empty riddle name")
		return
	}

	r, err := p.getter()
	if err != nil {
		p.modal.ModalMessage(err.Error())
		return
	}

	err = r.SaveToFile(fmt.Sprintf("riddles/%s.json", name))
	if err != nil {
		p.modal.ModalMessage(err.Error())
		return
	}

	p.name.SetText("")
	p.modal.ModalMessage("Riddle saved")
}

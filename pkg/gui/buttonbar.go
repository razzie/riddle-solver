package gui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type ButtonBar struct {
	list     layout.List
	in       layout.Inset
	btns     []widget.Clickable
	btnNames []string
}

func NewButtonBar(btns ...string) ButtonBar {
	return ButtonBar{
		list:     layout.List{Axis: layout.Horizontal},
		in:       layout.UniformInset(unit.Dp(5)),
		btns:     make([]widget.Clickable, len(btns)),
		btnNames: btns,
	}
}

func (bb *ButtonBar) Clicked(idx int) bool {
	return bb.btns[idx].Clicked()
}

func (bb *ButtonBar) Layout(gtx C, th *material.Theme) D {
	return bb.list.Layout(gtx, len(bb.btns), func(gtx C, idx int) D {
		return bb.in.Layout(gtx, material.Button(th, &bb.btns[idx], bb.btnNames[idx]).Layout)
	})
}

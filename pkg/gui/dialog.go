package gui

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Dialog struct {
	Message string
	ButtonBar
}

func NewDialog(message string, buttons ...string) Dialog {
	return Dialog{
		Message:   message,
		ButtonBar: NewButtonBar(buttons...),
	}
}

func (d *Dialog) Layout(gtx C, th *material.Theme) D {
	return layout.Center.Layout(gtx, func(gtx C) D {
		in := layout.UniformInset(unit.Dp(8))
		gtx.Constraints.Max.X = gtx.Px(unit.Dp(400))
		gtx.Constraints.Max.Y = gtx.Px(unit.Dp(100))
		gtx.Constraints.Min = gtx.Constraints.Max
		rr := float32(gtx.Px(unit.Dp(8)))
		clip.UniformRRect(f32.Rectangle{Max: f32.Point{
			X: float32(gtx.Constraints.Min.X),
			Y: float32(gtx.Constraints.Min.Y),
		}}, rr).Add(gtx.Ops)
		paint.Fill(gtx.Ops, th.Bg)
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					return in.Layout(gtx, material.Body1(th, d.Message).Layout)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					return d.ButtonBar.Layout(gtx, th)
				})
			}),
		)
	})
}

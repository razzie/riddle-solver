package gui

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type ListWithScrollbar struct {
	layout.List
}

func (l *ListWithScrollbar) Layout(gtx C, th *material.Theme, len int, w layout.ListElement) D {
	return layout.Stack{Alignment: layout.NE}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return l.List.Layout(gtx, len, w)
		}),
		layout.Stacked(func(gtx C) D {
			scrollbarThickness := gtx.Px(unit.Dp(8))
			if l.Axis == layout.Vertical {
				max := gtx.Constraints.Max.Y
				if l.Position.Length < max {
					return D{}
				}
				scrollbar := f32.Rect(
					0,
					float32(l.Position.Offset)/float32(l.Position.Length),
					float32(scrollbarThickness),
					float32(gtx.Constraints.Max.Y)-(float32(l.Position.OffsetLast)/float32(l.Position.Length)),
				)
				rr := float32(gtx.Px(unit.Dp(4)))
				clip.UniformRRect(scrollbar, rr).Add(gtx.Ops)
				paint.Fill(gtx.Ops, th.ContrastBg)
				return layout.Dimensions{
					Size: image.Pt(scrollbarThickness, gtx.Constraints.Max.X),
				}
			} else {
				return D{}
			}
		}),
	)
}

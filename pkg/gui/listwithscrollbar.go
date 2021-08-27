package gui

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type ListWithScrollbar struct {
	layout.List
}

func (l *ListWithScrollbar) Layout(gtx C, th *material.Theme, len int, w layout.ListElement) D {
	fakeGtx := gtx
	fakeGtx.Ops = new(op.Ops)
	fakeGtx.Constraints.Min = image.Pt(0, 0)
	contentLen := 0
	wSizes := make([]int, len)
	for i := 0; i < len; i++ {
		size := l.Axis.Convert(w(fakeGtx, i).Size).X
		wSizes[i] = size
		contentLen += size
	}

	var stack layout.Stack
	if l.Axis == layout.Vertical {
		stack.Alignment = layout.NE
	} else {
		stack.Alignment = layout.SW
	}
	return stack.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return l.List.Layout(gtx, len, w)
		}),
		layout.Stacked(func(gtx C) D {
			max := l.Axis.Convert(gtx.Constraints.Max).X
			if contentLen < max {
				return D{}
			}

			offset := l.Position.Offset
			for i := 0; i < l.Position.First; i++ {
				offset += wSizes[i]
			}

			scale := float32(max) / float32(contentLen)
			scrollbarThickness := float32(gtx.Px(unit.Dp(8)))
			scrollbarStart := float32(offset) * scale
			scrollbarLen := float32(max) * scale

			var scrollbar f32.Rectangle
			if l.Axis == layout.Vertical {
				scrollbar = f32.Rect(
					0,
					scrollbarStart,
					scrollbarThickness,
					scrollbarStart+scrollbarLen,
				)
			} else {
				scrollbar = f32.Rect(
					scrollbarStart,
					0,
					scrollbarStart+scrollbarLen,
					scrollbarThickness,
				)
			}
			rr := scrollbarThickness / 2
			clip.UniformRRect(scrollbar, rr).Add(gtx.Ops)
			paint.Fill(gtx.Ops, th.ContrastBg)
			return D{}
		}),
	)
}

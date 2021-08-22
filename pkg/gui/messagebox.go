package gui

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type MessageBox struct {
	theme   *material.Theme
	message string
	okBtn   widget.Clickable
	okFunc  func()
}

func NewMessageBox(th *material.Theme, message string, okFunc func()) layout.Widget {
	mbox := &MessageBox{
		theme:   th,
		message: message,
		okFunc:  okFunc,
	}
	return func(gtx C) D {
		return layout.Center.Layout(gtx, mbox.Layout)
	}
}

func (mbox *MessageBox) Layout(gtx C) D {
	in := layout.UniformInset(unit.Dp(8))
	gtx.Constraints.Max.X = gtx.Px(unit.Dp(400))
	gtx.Constraints.Max.Y = gtx.Px(unit.Dp(100))
	gtx.Constraints.Min = gtx.Constraints.Max
	rr := float32(gtx.Px(unit.Dp(8)))
	clip.UniformRRect(f32.Rectangle{Max: f32.Point{
		X: float32(gtx.Constraints.Min.X),
		Y: float32(gtx.Constraints.Min.Y),
	}}, rr).Add(gtx.Ops)
	dr := image.Rectangle{Max: gtx.Constraints.Min}
	paint.FillShape(gtx.Ops,
		mbox.theme.Bg,
		clip.Rect(dr).Op(),
	)
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return in.Layout(gtx, material.Body1(mbox.theme, mbox.message).Layout)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				if mbox.okBtn.Clicked() && mbox.okFunc != nil {
					mbox.okFunc()
				}
				gtx.Constraints.Max.X = gtx.Px(unit.Dp(75))
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return in.Layout(gtx, material.Button(mbox.theme, &mbox.okBtn, "OK").Layout)
			})
		}),
	)
}

type YesNoMessageBox struct {
	theme     *material.Theme
	message   string
	yesBtn    widget.Clickable
	noBtn     widget.Clickable
	yesNoFunc func(bool)
}

func NewYesNoMessageBox(th *material.Theme, message string, yesNoFunc func(bool)) layout.Widget {
	mbox := &YesNoMessageBox{
		theme:     th,
		message:   message,
		yesNoFunc: yesNoFunc,
	}
	return func(gtx C) D {
		return layout.Center.Layout(gtx, mbox.Layout)
	}
}

func (mbox *YesNoMessageBox) Layout(gtx C) D {
	in := layout.UniformInset(unit.Dp(8))
	gtx.Constraints.Max.X = gtx.Px(unit.Dp(400))
	gtx.Constraints.Max.Y = gtx.Px(unit.Dp(100))
	gtx.Constraints.Min = gtx.Constraints.Max
	rr := float32(gtx.Px(unit.Dp(8)))
	clip.UniformRRect(f32.Rectangle{Max: f32.Point{
		X: float32(gtx.Constraints.Min.X),
		Y: float32(gtx.Constraints.Min.Y),
	}}, rr).Add(gtx.Ops)
	dr := image.Rectangle{Max: gtx.Constraints.Min}
	paint.FillShape(gtx.Ops,
		mbox.theme.Bg,
		clip.Rect(dr).Op(),
	)
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return in.Layout(gtx, material.Body1(mbox.theme, mbox.message).Layout)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				if mbox.yesNoFunc != nil {
					if mbox.yesBtn.Clicked() {
						mbox.yesNoFunc(true)
					}
					if mbox.noBtn.Clicked() {
						mbox.yesNoFunc(false)
					}
				}
				gtx.Constraints.Max.X = gtx.Px(unit.Dp(200))
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(0.5, func(gtx C) D {
						return in.Layout(gtx, material.Button(mbox.theme, &mbox.yesBtn, "Yes").Layout)
					}),
					layout.Flexed(0.5, func(gtx C) D {
						return in.Layout(gtx, material.Button(mbox.theme, &mbox.noBtn, "No").Layout)
					}),
				)
			})
		}),
	)
}

package gui

import (
	"fmt"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/razzie/razgio"
	"github.com/razzie/riddle-solver/pkg/riddle"
)

type SetupPage struct {
	theme    *material.Theme
	modal    razgio.ModalHandler
	list     razgio.ListWithScrollbar
	items    []setupItem
	buttons  razgio.ButtonBar
	saveFunc func(riddle.Setup)
	removeCh chan int
}

func NewSetupPage(th *material.Theme, modal razgio.ModalHandler) *SetupPage {
	p := &SetupPage{
		theme:    th,
		modal:    modal,
		list:     razgio.NewListWithScrollbar(),
		buttons:  razgio.NewButtonBar("Add item type", "Save / apply", "Reset"),
		removeCh: make(chan int, 1),
	}
	p.buttons.SetButtonIcon(0, razgio.GetIcons().ContentAdd)
	p.Reset()
	return p
}

func (p *SetupPage) GetName() string {
	return "Setup"
}

func (p *SetupPage) Select() {

}

func (p *SetupPage) update() {
	select {
	case idx := <-p.removeCh:
		p.items = append(p.items[:idx], p.items[idx+1:]...)
	default:
	}
}

func (p *SetupPage) Layout(gtx C) D {
	p.update()

	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	in := layout.UniformInset(unit.Dp(5))
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			if len(p.items) == 0 {
				return D{Size: gtx.Constraints.Max}
			}
			return p.list.Layout(gtx, p.theme, len(p.items), func(gtx C, idx int) D {
				return in.Layout(gtx, func(gtx C) D {
					return p.items[idx].Layout(gtx, p.theme, idx)
				})
			})
		}),
		layout.Rigid(func(gtx C) D {
			switch {
			case p.buttons.Clicked(0):
				p.Add()
			case p.buttons.Clicked(1):
				p.Save()
			case p.buttons.Clicked(2):
				p.modal.ModalYesNo("Reset setup?\nThis will remove all rules as well", p.Reset)
			}
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return p.buttons.Layout(gtx, p.theme)
		}),
	)
}

func (p *SetupPage) SetSaveFunc(saveFunc func(riddle.Setup)) {
	p.saveFunc = saveFunc
}

func (p *SetupPage) GetSetup() (riddle.Setup, error) {
	var setup = make(riddle.Setup)

	for _, item := range p.items {
		itemType := item.itemType.Text()
		values := strings.Split(item.values.Text(), ",")
		if len(itemType) == 0 {
			if len(values) > 0 && len(values[0]) > 0 {
				return nil, fmt.Errorf("cannot have values without item type")
			}
			continue
		}

		var trimmedValues []string
		for _, value := range values {
			trimmedValue := strings.Trim(value, " ")
			if len(trimmedValue) == 0 {
				continue
			}
			trimmedValues = append(trimmedValues, trimmedValue)
		}
		if len(trimmedValues) == 0 {
			continue
		}

		setup[itemType] = trimmedValues
	}

	if err := setup.Check(); err != nil {
		return nil, err
	}

	return setup, nil
}

func (p *SetupPage) SetSetup(setup riddle.Setup) {
	var newItems []setupItem
	for itemType, values := range setup {
		newItem := newSetupItem(p)
		newItem.itemType.SetText(itemType)
		newItem.values.SetText(strings.Join(values, ", "))
		newItems = append(newItems, newItem)
	}
	p.items = newItems
	p.Save()
}

func (p *SetupPage) Add() {
	p.items = append(p.items, newSetupItem(p))
}

func (p *SetupPage) remove(idx int) {
	p.removeCh <- idx
}

func (p *SetupPage) Save() {
	setup, err := p.GetSetup()
	if err != nil {
		p.modal.ModalMessage(err.Error())
		return
	}

	if p.saveFunc != nil {
		p.saveFunc(setup)
	}
}

func (p *SetupPage) Reset() {
	p.items = []setupItem{newSetupItem(p)}
	p.Save()
}

type setupItem struct {
	list       layout.List
	itemType   razgio.TextField
	values     razgio.TextField
	delete     widget.Clickable
	deleteIcon *widget.Icon
	p          *SetupPage
}

func newSetupItem(p *SetupPage) setupItem {
	return setupItem{
		list: layout.List{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		},
		itemType: razgio.TextField{
			Editor: widget.Editor{SingleLine: true},
		},
		values: razgio.TextField{
			Editor: widget.Editor{SingleLine: true},
		},
		deleteIcon: razgio.GetIcons().ActionDelete,
		p:          p,
	}
}

func (item *setupItem) Layout(gtx C, th *material.Theme, idx int) D {
	in := layout.Inset{Left: unit.Dp(8)}
	maxWidth := gtx.Constraints.Max.X
	gtx.Constraints.Min.X = maxWidth
	widgets := [...]layout.Widget{
		func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(th.TextSize.Scale(2))
			return material.Label(th, th.TextSize, fmt.Sprintf("#%d", idx+1)).Layout(gtx)
		},
		func(gtx C) D {
			gtx.Constraints.Max.X = maxWidth / 3
			return in.Layout(gtx, func(gtx C) D {
				return item.itemType.Layout(gtx, th, "item type", "e.g. color")
			})
		},
		func(gtx C) D {
			gtx.Constraints.Max.X = maxWidth / 2
			return in.Layout(gtx, func(gtx C) D {
				return item.values.Layout(gtx, th, "values", "e.g. red, green, blue")
			})
		},
		func(gtx C) D {
			if len(item.itemType.Text()) > 0 || len(item.values.Text()) > 0 {
				gtx.Queue = nil
			} else if item.delete.Clicked() {
				item.p.remove(idx)
			}
			btn := razgio.IconAndTextButton(th, &item.delete, item.deleteIcon, "")
			btn.TextSize = th.TextSize.Scale(1.5)
			btn.Inset = layout.UniformInset(unit.Dp(2))
			return in.Layout(gtx, btn.Layout)
		},
	}
	return item.list.Layout(gtx, len(widgets), func(gtx C, idx int) D {
		return widgets[idx](gtx)
	})
}

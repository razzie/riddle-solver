package gui

import (
	"fmt"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/razzie/riddle-solver/pkg/riddle"
)

var _ Page = (*SetupPage)(nil)

type SetupPage struct {
	theme    *material.Theme
	modal    ModalHandler
	list     ListWithScrollbar
	items    []setupItem
	buttons  ButtonBar
	saveFunc func(riddle.Setup)
}

func NewSetupPage(th *material.Theme, modal ModalHandler) *SetupPage {
	p := &SetupPage{
		theme:   th,
		modal:   modal,
		list:    NewListWithScrollbar(),
		buttons: NewButtonBar("Add item type", "Save / apply", "Reset"),
	}
	p.Reset()
	return p
}

func (p *SetupPage) GetName() string {
	return "Setup"
}

func (p *SetupPage) Select() {

}

func (p *SetupPage) Layout(gtx C) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	in := layout.UniformInset(unit.Dp(5))
	return p.list.Layout(gtx, p.theme, len(p.items)+1, func(gtx C, idx int) D {
		if idx < len(p.items) {
			return in.Layout(gtx, func(gtx C) D {
				return p.items[idx].Layout(gtx, p.theme, idx)
			})
		}
		switch {
		case p.buttons.Clicked(0):
			p.Add()
		case p.buttons.Clicked(1):
			p.Save()
		case p.buttons.Clicked(2):
			p.modal.ModalYesNo("Are you sure?", p.Reset)
		}
		return p.buttons.Layout(gtx, p.theme)
	})
}

func (p *SetupPage) SetSaveFunc(saveFunc func(riddle.Setup)) {
	p.saveFunc = saveFunc
}

func (p *SetupPage) GetSetup() (riddle.Setup, error) {
	var setup = make(riddle.Setup)

	for _, item := range p.items {
		itemType := item.itemType.Text()
		values := strings.Split(item.values.Text(), "")
		if len(itemType) == 0 {
			if len(values) > 0 {
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
		newItem := newSetupItem()
		newItem.itemType.SetText(itemType)
		newItem.values.SetText(strings.Join(values, ", "))
		newItems = append(newItems, newItem)
	}
	p.items = newItems
	p.Save()
}

func (p *SetupPage) Add() {
	p.items = append(p.items, setupItem{})
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
	p.items = []setupItem{newSetupItem()}
}

type setupItem struct {
	list     layout.List
	itemType component.TextField
	values   component.TextField
}

func newSetupItem() setupItem {
	return setupItem{}
}

func (item *setupItem) Layout(gtx C, th *material.Theme, idx int) D {
	in := layout.Inset{Left: unit.Dp(5)}
	maxWidth := gtx.Constraints.Max.X
	gtx.Constraints.Min.X = maxWidth
	widgets := [...]layout.Widget{
		material.Label(th, th.TextSize, fmt.Sprintf("#%d", idx+1)).Layout,
		func(gtx C) D {
			gtx.Constraints.Max.X = maxWidth / 3
			return in.Layout(gtx, func(gtx C) D {
				return item.itemType.Layout(gtx, th, "item type")
			})
		},
		func(gtx C) D {
			gtx.Constraints.Max.X = maxWidth / 2
			return in.Layout(gtx, func(gtx C) D {
				return item.values.Layout(gtx, th, "values (comma separated)")
			})
		},
	}
	return item.list.Layout(gtx, len(widgets), func(gtx C, idx int) D {
		return widgets[idx](gtx)
	})
}

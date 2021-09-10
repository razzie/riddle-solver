package gui

import (
	"io"
	"path/filepath"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/explorer"
	"gioui.org/x/richtext"
	"github.com/razzie/razgio"
	"github.com/razzie/riddle-solver/pkg/riddle"
)

type LoadPage struct {
	theme  *material.Theme
	modal  razgio.ModalHandler
	list   razgio.ListWithScrollbar
	btns   razgio.ButtonBar
	setter func(*riddle.Riddle) error
	items  []*loadPageItem
}

func NewLoadPage(th *material.Theme, modal razgio.ModalHandler) *LoadPage {
	return &LoadPage{
		theme: th,
		modal: modal,
		list:  razgio.NewListWithScrollbar(),
		btns:  razgio.NewButtonBar("Explore", "Refresh"),
	}
}

func (p *LoadPage) GetName() string {
	return "Load"
}

func (p *LoadPage) Select() {
	p.items = p.items[:0]
	addRiddle := func(name string, builtin bool, loader func() (*riddle.Riddle, error)) {
		p.items = append(p.items, &loadPageItem{
			Name:    name,
			BuiltIn: builtin,
			Loader:  loader,
		})
	}

	addRiddle("Einstein's 5 house riddle", true, func() (*riddle.Riddle, error) {
		return riddle.NewEinsteinRiddle(), nil
	})

	addRiddle("The Jindosh riddle", true, func() (*riddle.Riddle, error) {
		return riddle.NewJindoshRiddle(), nil
	})

	files, _ := filepath.Glob("riddles/*.json")
	for _, file := range files {
		filename := file
		addRiddle(filename, false, func() (*riddle.Riddle, error) {
			return riddle.LoadRiddleFromFile(filename)
		})
	}
}

func (p *LoadPage) Layout(gtx C) D {
	if p.btns.Clicked(0) {
		p.loadFromExplorer()
	}
	if p.btns.Clicked(1) {
		p.Select()
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return p.list.Layout(gtx, p.theme, len(p.items), func(gtx C, idx int) D {
				dims := p.items[idx].Layout(gtx, p)
				dims.Size.Y += gtx.Px(unit.Dp(12))
				return dims
			})
		}),
		layout.Rigid(func(gtx C) D {
			return p.btns.Layout(gtx, p.theme)
		}),
	)
}

func (p *LoadPage) SetRiddleSetter(setter func(*riddle.Riddle) error) {
	p.setter = setter
}

func (p *LoadPage) loadFromExplorer() {
	reader, err := explorer.ReadFile()
	if err != nil {
		p.modal.ModalMessage(err.Error())
		return
	}
	defer reader.Close()

	bytes, err := io.ReadAll(reader)
	if err != nil {
		p.modal.ModalMessage(err.Error())
		return
	}

	r, err := riddle.LoadRiddle(bytes)
	if err != nil {
		p.modal.ModalMessage(err.Error())
		return
	}

	p.loadRiddle(r)
}

func (p *LoadPage) loadItem(item *loadPageItem) {
	r, err := item.Loader()
	if err != nil {
		p.modal.ModalMessage(err.Error())
		return
	}
	p.loadRiddle(r)
}

func (p *LoadPage) loadRiddle(r *riddle.Riddle) {
	if p.setter != nil {
		if err := p.setter(r); err != nil {
			p.modal.ModalMessage(err.Error())
		}
	}
}

type loadPageItem struct {
	widget.Clickable
	richtext.InteractiveText
	Name    string
	Loader  func() (*riddle.Riddle, error)
	BuiltIn bool
}

func (item *loadPageItem) Layout(gtx C, p *LoadPage) D {
	th := p.theme
	normal := func(text string) richtext.SpanStyle {
		return richtext.SpanStyle{
			Content: text,
			Color:   th.Fg,
			Size:    th.TextSize,
		}
	}
	colored := func(text string) richtext.SpanStyle {
		return richtext.SpanStyle{
			Content: text,
			Color:   th.ContrastBg,
			Size:    th.TextSize,
		}
	}
	var spans []richtext.SpanStyle
	if item.BuiltIn {
		spans = []richtext.SpanStyle{
			normal(item.Name + " ["),
			colored("built-in"),
			normal("]"),
		}
	} else {
		spans = []richtext.SpanStyle{normal(item.Name)}
	}
	if item.Hovered() {
		for i := range spans {
			spans[i].Font.Weight = text.Bold
		}
	}
	if item.Clicked() {
		p.loadItem(item)
	}
	dims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			dims := richtext.Text(&item.InteractiveText, th.Shaper, spans...).Layout(gtx)
			gtx.Constraints.Min = dims.Size
			gtx.Constraints.Max = dims.Size
			item.Clickable.Layout(gtx)
			return dims
		}),
	)
	dims.Size.X = gtx.Constraints.Max.X
	return dims
}

package gui

import (
	"fmt"

	"gioui.org/widget/material"
	"github.com/razzie/razgio"
	"github.com/razzie/riddle-solver/pkg/riddle"
)

type DebugPage struct {
	theme     *material.Theme
	modal     razgio.ModalHandler
	scrollbar razgio.ListWithScrollbar
	results   razgio.Tree
	setup     riddle.Setup
	rules     []riddle.Rule
	dirty     bool
}

func NewDebugPage(th *material.Theme, modal razgio.ModalHandler) *DebugPage {
	return &DebugPage{
		theme:     th,
		modal:     modal,
		scrollbar: razgio.NewListWithScrollbar(),
		results: razgio.NewTree(
			razgio.TreeLabel{Text: "Solver internals ("},
			razgio.TreeLabel{Text: "XXX", Highlight: true},
			razgio.TreeLabel{Text: " steps)"},
		),
		dirty: true,
	}
}

func (p *DebugPage) GetName() string {
	return "Debug"
}

func (p *DebugPage) Select() {
	if !p.dirty {
		return
	}

	solver := riddle.NewSolver(p.setup, p.rules)
	steps, err := solver.Solve(solver.GuessPrimaryItemType())
	if err != nil {
		p.modal.ModalMessage(err.Error())
	}
	p.results.Label[1].Text = fmt.Sprint(steps)

	p.dirty = false
	p.results.Bool.Value = true
	p.results.ClearChildren()

	for i, entry := range solver.Entries {
		node := p.results.AddChild(razgio.TreeLabel{Text: fmt.Sprintf("Entry #%d", i+1)})
		for itemType, values := range entry {
			labelParts := make([]razgio.TreeLabel, 0, len(values)*2) // 1 + len(values) + (len(values)-1)
			labelParts = append(labelParts, razgio.TreeLabel{Text: itemType + ": "})
			for i := 0; i < len(values); i++ {
				labelParts = append(labelParts, razgio.TreeLabel{Text: values[i], Highlight: true})
				if i < len(values)-1 {
					labelParts = append(labelParts, razgio.TreeLabel{Text: ", "})
				}
			}
			node.AddChild(labelParts...)
		}
	}
}

func (p *DebugPage) Layout(gtx C) D {
	return p.scrollbar.Layout(gtx, p.theme, 1, func(gtx C, idx int) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return p.results.Layout(gtx, p.theme)
	})
}

func (p *DebugPage) HandleSetup(setup riddle.Setup) {
	p.setup = setup
	p.dirty = true
}

func (p *DebugPage) HandleRules(rules []riddle.Rule) {
	p.rules = rules
	p.dirty = true
}

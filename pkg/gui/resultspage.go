package gui

import (
	"fmt"

	"gioui.org/widget/material"
	"github.com/razzie/razgio"
	"github.com/razzie/riddle-solver/pkg/riddle"
)

type ResultsPage struct {
	theme     *material.Theme
	modal     razgio.ModalHandler
	scrollbar razgio.ListWithScrollbar
	results   razgio.Tree
	setup     riddle.Setup
	rules     []riddle.Rule
	dirty     bool
}

func NewResultsPage(th *material.Theme, modal razgio.ModalHandler) *ResultsPage {
	return &ResultsPage{
		theme:     th,
		modal:     modal,
		scrollbar: razgio.NewListWithScrollbar(),
		results:   razgio.NewTree(razgio.TreeLabel{Text: "Results"}),
		dirty:     true,
	}
}

func (p *ResultsPage) GetName() string {
	return "Results"
}

func (p *ResultsPage) Select() {
	if !p.dirty {
		return
	}

	solver := riddle.NewSolver(p.setup, p.rules)
	_, err := solver.Solve(solver.GuessPrimaryItemType())
	if err != nil {
		p.modal.ModalMessage(err.Error())
	} else if solver.IsSolved() {
		p.modal.ModalMessage("Riddle solved")
	}

	p.dirty = false
	p.results.Bool.Value = true
	p.results.ClearChildren()

	for itemType, values := range p.setup {
		itemTypeNode := p.results.AddChild(razgio.TreeLabel{Text: itemType})
		for _, val := range values {
			item := riddle.Item(fmt.Sprintf("%s:%s", itemType, val))
			valueNode := itemTypeNode.AddChild(razgio.TreeLabel{Text: val})
			for itemType, values := range solver.FindAssociatedItems(item) {
				labelParts := make([]razgio.TreeLabel, 0, len(values)*2) // 1 + len(values) + (len(values)-1)
				labelParts = append(labelParts, razgio.TreeLabel{Text: itemType + ": "})
				for i := 0; i < len(values); i++ {
					labelParts = append(labelParts, razgio.TreeLabel{Text: values[i], Highlight: true})
					if i < len(values)-1 {
						labelParts = append(labelParts, razgio.TreeLabel{Text: ", "})
					}
				}
				valueNode.AddChild(labelParts...)
			}
		}
	}
}

func (p *ResultsPage) Layout(gtx C) D {
	return p.scrollbar.Layout(gtx, p.theme, 1, func(gtx C, idx int) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return p.results.Layout(gtx, p.theme)
	})
}

func (p *ResultsPage) HandleSetup(setup riddle.Setup) {
	p.setup = setup
	p.dirty = true
}

func (p *ResultsPage) HandleRules(rules []riddle.Rule) {
	p.rules = rules
	p.dirty = true
}

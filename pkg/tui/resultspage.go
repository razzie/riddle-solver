package tui

import (
	"fmt"

	"github.com/razzie/riddle-solver/pkg/riddle"
	"github.com/rivo/tview"
)

// ResultsPage is a UI tree element that lets the user browse the riddle results
type ResultsPage struct {
	Page
	tree  *tview.TreeView
	root  *tview.TreeNode
	setup riddle.Setup
	rules []riddle.Rule
	dirty bool
	modal ModalHandler
}

// NewResultsTree returns a new ResultsTree
func NewResultsTree(modal ModalHandler) *ResultsPage {
	root := tview.NewTreeNode("Results")
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	p := &ResultsPage{
		Page:  NewPage(tview.NewFrame(tree), "Results"),
		tree:  tree,
		root:  root,
		dirty: true,
		modal: modal,
	}
	p.Page.SetSelectFunc(p.Update)
	return p
}

// Update updates the results based on the latest setup and rules
func (p *ResultsPage) Update() {
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
	p.root.ClearChildren()

	for itemType, values := range p.setup {
		itemTypeNode := tview.NewTreeNode(itemType).SetExpanded(false)
		for _, val := range values {
			item := riddle.Item(fmt.Sprintf("%s:%s", itemType, val))
			valueNode := tview.NewTreeNode(val).SetReference(item)
			itemTypeNode.AddChild(valueNode)
		}
		p.root.AddChild(itemTypeNode)
	}

	p.tree.SetSelectedFunc(func(node *tview.TreeNode) {
		children := node.GetChildren()
		if len(children) == 0 {
			reference := node.GetReference()
			if reference == nil {
				return
			}

			item := reference.(riddle.Item)
			result := solver.FindAssociatedItems(item)
			for itemType, values := range result {
				text := colorizeItems(itemType, values)
				resultNode := tview.NewTreeNode(text)
				node.AddChild(resultNode)
			}

		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	})
}

// HandleSetup updates the inner stored setup
func (p *ResultsPage) HandleSetup(setup riddle.Setup) {
	p.setup = setup
}

// HandleRules updates the inner stored rules
func (p *ResultsPage) HandleRules(rules []riddle.Rule) {
	p.rules = rules
	p.dirty = true
}

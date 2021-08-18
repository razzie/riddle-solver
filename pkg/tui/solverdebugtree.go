package tui

import (
	"fmt"

	"github.com/razzie/riddle-solver/pkg/riddle"
	"github.com/rivo/tview"
)

// SolverDebugPage shows solver internals
type SolverDebugPage struct {
	Page
	tree  *tview.TreeView
	root  *tview.TreeNode
	setup riddle.Setup
	rules []riddle.Rule
	dirty bool
	modal ModalHandler
}

// NewSolverDebugTree returns a new SolverDebugTree
func NewSolverDebugTree(modal ModalHandler) *SolverDebugPage {
	root := tview.NewTreeNode("")
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	p := &SolverDebugPage{
		Page:  NewPage(tview.NewFrame(tree), "Debug"),
		tree:  tree,
		root:  root,
		dirty: true,
		modal: modal,
	}
	p.Page.SetSelectFunc(p.Update)
	return p
}

// Update updates the results based on the latest setup and rules
func (p *SolverDebugPage) Update() {
	if !p.dirty {
		return
	}

	solver := riddle.NewSolver(p.setup, p.rules)
	steps, err := solver.Solve(solver.GuessPrimaryItemType())
	if err != nil {
		p.modal.ModalMessage(err.Error())
	}

	p.dirty = false
	p.root.
		SetText(fmt.Sprintf("Solver internals (%d steps)", steps)).
		ClearChildren()

	for i, entry := range solver.Entries {
		node := tview.NewTreeNode(fmt.Sprintf("Entry #%d", i+1)).SetExpanded(false)
		for itemType, values := range entry {
			text := colorizeItems(itemType, values)
			resultNode := tview.NewTreeNode(text)
			node.AddChild(resultNode)
		}
		p.root.AddChild(node)
	}

	p.tree.SetSelectedFunc(func(node *tview.TreeNode) {
		node.SetExpanded(!node.IsExpanded())
	})
}

// HandleSetup updates the inner stored setup
func (p *SolverDebugPage) HandleSetup(setup riddle.Setup) {
	p.setup = setup
}

// HandleRules updates the inner stored rules
func (p *SolverDebugPage) HandleRules(rules []riddle.Rule) {
	p.rules = rules
	p.dirty = true
}

package ui

import (
	"fmt"

	"github.com/razzie/riddle-solver/riddle"
	"github.com/rivo/tview"
)

// SolverDebugTree shows solver internals
type SolverDebugTree struct {
	*tview.TreeView
	root  *tview.TreeNode
	setup riddle.Setup
	rules []riddle.Rule
	dirty bool
	modal ModalHandler
}

// NewSolverDebugTree returns a new SolverDebugTree
func NewSolverDebugTree(modal ModalHandler) *SolverDebugTree {
	root := tview.NewTreeNode("")
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	return &SolverDebugTree{
		TreeView: tree,
		root:     root,
		dirty:    true,
		modal:    modal,
	}
}

// Update updates the results based on the latest setup and rules
func (t *SolverDebugTree) Update() {
	if !t.dirty {
		return
	}

	solver := riddle.NewSolver(t.setup)
	steps, err := solver.ApplyRules(t.rules)
	if err != nil {
		t.modal.ModalMessage(fmt.Sprint(err))
	}

	t.dirty = false
	t.root.
		SetText(fmt.Sprintf("Solver internals (%d steps)", steps)).
		ClearChildren()

	for i, entry := range solver.Entries {
		node := tview.NewTreeNode(fmt.Sprintf("Entry #%d", i+1)).SetExpanded(false)
		for itemType, values := range entry {
			text := colorizeItems(itemType, values)
			resultNode := tview.NewTreeNode(text)
			node.AddChild(resultNode)
		}
		t.root.AddChild(node)
	}

	t.SetSelectedFunc(func(node *tview.TreeNode) {
		node.SetExpanded(!node.IsExpanded())
	})
}

// HandleSetup updates the inner stored setup
func (t *SolverDebugTree) HandleSetup(setup riddle.Setup) {
	t.setup = setup
}

// HandleRules updates the inner stored rules
func (t *SolverDebugTree) HandleRules(rules []riddle.Rule) {
	t.rules = rules
	t.dirty = true
}

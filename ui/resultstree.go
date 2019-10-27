package ui

import (
	"github.com/razzie/riddle-solver/solver"
	"github.com/rivo/tview"
)

// ResultsTree is a UI tree element that lets the user browse the riddle results
type ResultsTree struct {
	*tview.TreeView
	root  *tview.TreeNode
	setup solver.Setup
	rules []solver.Rule
	dirty bool
}

// NewResultsTree returns a new ResultsTree
func NewResultsTree() *ResultsTree {
	root := tview.NewTreeNode("Results")
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	t := &ResultsTree{
		TreeView: tree,
		root:     root}
	t.root.SetSelectedFunc(func() { t.Update() })
	return t
}

// Update updates the results based on the latest setup and rules
func (t *ResultsTree) Update() {
	if !t.dirty {
		return
	}

	t.dirty = false
}

// HandleSetup updates the inner stored setup
func (t *ResultsTree) HandleSetup(setup solver.Setup) {
	t.setup = setup
}

// HandleRules updates the inner stored rules
func (t *ResultsTree) HandleRules(rules []solver.Rule) {
	t.rules = rules
	t.dirty = true
}

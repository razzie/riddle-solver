package ui

import (
	"github.com/razzie/riddle-solver/riddle"
	"github.com/rivo/tview"
)

// ResultsTree is a UI tree element that lets the user browse the riddle results
type ResultsTree struct {
	*tview.TreeView
	root  *tview.TreeNode
	setup riddle.Setup
	rules []riddle.Rule
	dirty bool
}

// NewResultsTree returns a new ResultsTree
func NewResultsTree() *ResultsTree {
	root := tview.NewTreeNode("Results")
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	return &ResultsTree{
		TreeView: tree,
		root:     root}
}

// Update updates the results based on the latest setup and rules
func (t *ResultsTree) Update() {
	if !t.dirty {
		return
	}

	t.dirty = false
}

// HandleSetup updates the inner stored setup
func (t *ResultsTree) HandleSetup(setup riddle.Setup) {
	t.setup = setup
}

// HandleRules updates the inner stored rules
func (t *ResultsTree) HandleRules(rules []riddle.Rule) {
	t.rules = rules
	t.dirty = true
}

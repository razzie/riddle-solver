package ui

import (
	"fmt"

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
	modal ModalHandler
}

// NewResultsTree returns a new ResultsTree
func NewResultsTree(modal ModalHandler) *ResultsTree {
	root := tview.NewTreeNode("Results")
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	return &ResultsTree{
		TreeView: tree,
		root:     root,
		dirty:    true,
		modal:    modal,
	}
}

// Update updates the results based on the latest setup and rules
func (t *ResultsTree) Update() {
	if !t.dirty {
		return
	}

	solver := riddle.NewSolver(t.setup)
	_, err := solver.ApplyRules(t.rules)
	if err != nil {
		t.modal.ModalMessage(fmt.Sprint(err))
	}

	t.dirty = false
	t.root.ClearChildren()

	for itemType, values := range t.setup {
		itemTypeNode := tview.NewTreeNode(itemType).SetExpanded(false)
		for _, val := range values {
			item := riddle.Item(fmt.Sprintf("%s:%s", itemType, val))
			valueNode := tview.NewTreeNode(val).SetReference(item)
			itemTypeNode.AddChild(valueNode)
		}
		t.root.AddChild(itemTypeNode)
	}

	t.SetSelectedFunc(func(node *tview.TreeNode) {
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
func (t *ResultsTree) HandleSetup(setup riddle.Setup) {
	t.setup = setup
}

// HandleRules updates the inner stored rules
func (t *ResultsTree) HandleRules(rules []riddle.Rule) {
	t.rules = rules
	t.dirty = true
}

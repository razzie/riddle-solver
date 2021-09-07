package gui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Tree struct {
	widget.Bool
	Name     string
	Children []*Tree
}

func (tree *Tree) AddChild(name string) *Tree {
	child := &Tree{Name: name}
	tree.Children = append(tree.Children, child)
	return child
}

func (tree *Tree) ClearChildren() {
	tree.Children = nil
}

func (tree *Tree) Layout(gtx C, th *material.Theme) D {
	return tree.layout(gtx, convertToTreeTheme(th))
}

func (tree *Tree) layout(gtx C, th *material.Theme) D {
	if len(tree.Children) == 0 {
		gtx.Queue = nil
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(material.CheckBox(th, &tree.Bool, tree.Name).Layout),
		layout.Rigid(func(gtx C) D {
			if !tree.Value {
				return D{}
			}
			children := make([]layout.FlexChild, len(tree.Children))
			for i, child := range tree.Children {
				childCopy := child
				children[i] = layout.Rigid(func(gtx C) D {
					return childCopy.layout(gtx, th)
				})
			}
			return layout.Inset{Left: th.TextSize.Scale(2)}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
			})
		}),
	)
}

func convertToTreeTheme(th *material.Theme) *material.Theme {
	clone := new(material.Theme)
	*clone = *th
	clone.Icon.CheckBoxChecked = GetIcons().CheckBoxIndeterminate
	clone.Icon.CheckBoxUnchecked = GetIcons().CheckBoxBlank
	return clone
}

package gui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/razzie/riddle-solver/pkg/riddle"
)

type AddRulePage struct {
	theme               *material.Theme
	modal               ModalHandler
	itemA               component.TextField
	itemB               component.TextField
	relation            widget.Enum
	hasCondition        widget.Bool
	conditionItemType   component.TextField
	conditionExpression component.TextField
	conditionReversible widget.Bool
	buttons             ButtonBar
	rule                *riddle.Rule
	setup               riddle.Setup
	saveFunc            func(*riddle.Rule)
}

func NewAddRulePage(th *material.Theme, modal ModalHandler) *AddRulePage {
	return &AddRulePage{
		theme:   th,
		modal:   modal,
		buttons: NewButtonBar("Save", "Reset"),
	}
}

func (p *AddRulePage) GetName() string {
	return "Add rule"
}

func (p *AddRulePage) Select() {

}

func (p *AddRulePage) Layout(gtx C) D {
	if p.buttons.Clicked(0) {
		p.Save()
	}
	if p.buttons.Clicked(1) {
		p.modal.ModalYesNo("Are you sure?", p.Reset)
	}

	items := []layout.FlexChild{
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Max.X /= 2
			return p.itemA.Layout(gtx, p.theme, "Item A")
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Max.X /= 2
			return p.itemB.Layout(gtx, p.theme, "Item B")
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(material.RadioButton(p.theme, &p.relation, riddle.RelAssociated.String(), "Associated").Layout),
				layout.Rigid(material.RadioButton(p.theme, &p.relation, riddle.RelDisassociated.String(), "Disassociated").Layout),
			)
		}),
		layout.Rigid(material.CheckBox(p.theme, &p.hasCondition, "Condition").Layout),
	}
	if p.hasCondition.Value {
		extraItems := []layout.FlexChild{
			layout.Rigid(func(gtx C) D {
				gtx.Constraints.Max.X /= 2
				return p.conditionItemType.Layout(gtx, p.theme, "Condition item type")
			}),
			layout.Rigid(func(gtx C) D {
				return p.conditionExpression.Layout(gtx, p.theme, "Condition expression")
			}),
			layout.Rigid(material.CheckBox(p.theme, &p.conditionReversible, "Reversible A <-> B").Layout),
		}
		items = append(items, extraItems...)
	}
	items = append(items, layout.Rigid(func(gtx C) D {
		return p.buttons.Layout(gtx, p.theme)
	}))
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, items...)
}

func (p *AddRulePage) HandleSetup(setup riddle.Setup) {
	p.Reset()
	p.setup = setup
}

func (p *AddRulePage) EditRule(rule *riddle.Rule) {
	p.rule = rule
	p.itemA.SetText(string(rule.ItemA))
	p.itemB.SetText(string(rule.ItemB))
	p.relation.Value = p.rule.Relation.String()
	p.hasCondition.Value = rule.HasCondition()
	p.conditionItemType.SetText(rule.ConditionItemType)
	p.conditionExpression.SetText(rule.Condition)
	p.conditionReversible.Value = rule.IsReversible
}

func (p *AddRulePage) Save() {
	var rule riddle.Rule
	rule.ItemA = riddle.Item(p.itemA.Text())
	rule.ItemB = riddle.Item(p.itemB.Text())
	switch p.relation.Value {
	case riddle.RelAssociated.String():
		rule.Relation = riddle.RelAssociated
	case riddle.RelDisassociated.String():
		rule.Relation = riddle.RelDisassociated
	default:
		p.modal.ModalMessage("Unknown relation")
		return
	}
	rule.ConditionItemType = p.conditionItemType.Text()
	rule.Condition = p.conditionExpression.Text()
	rule.IsReversible = p.conditionReversible.Value

	if err := rule.Check(p.setup); err != nil {
		p.modal.ModalMessage(err.Error())
		return
	}

	if p.rule == nil {
		p.rule = &rule
	} else {
		p.rule.ItemA = rule.ItemA
		p.rule.ItemB = rule.ItemB
		p.rule.Relation = rule.Relation
		p.rule.Condition = rule.Condition
		p.rule.ConditionItemType = rule.ConditionItemType
		p.rule.IsReversible = rule.IsReversible
	}

	if p.saveFunc != nil {
		p.saveFunc(p.rule)
		p.Reset()
	}
}

func (p *AddRulePage) SetSaveFunc(saveFunc func(*riddle.Rule)) {
	p.saveFunc = saveFunc
}

func (p *AddRulePage) Reset() {
	p.rule = nil
	p.itemA.SetText("")
	p.itemB.SetText("")
	p.relation.Value = ""
	p.hasCondition.Value = false
	p.conditionItemType.SetText("")
	p.conditionExpression.SetText("")
	p.conditionReversible.Value = false
}

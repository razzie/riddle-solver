package gui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/richtext"
	"github.com/razzie/riddle-solver/pkg/riddle"
)

type RulesPage struct {
	theme      *material.Theme
	modal      ModalHandler
	list       ListWithScrollbar
	btns       ButtonBar
	deleteIcon *widget.Icon
	rules      []ruleItem
	editFunc   func(*riddle.Rule)
	saveFunc   func([]riddle.Rule)
}

func NewRulesPage(th *material.Theme, modal ModalHandler) *RulesPage {
	p := &RulesPage{
		theme:      th,
		modal:      modal,
		list:       NewListWithScrollbar(),
		btns:       NewButtonBar("Add rule", "Reset"),
		deleteIcon: GetIcons().ActionDelete,
	}
	p.btns.SetButtonIcon(0, GetIcons().ContentAdd)
	return p
}

func (p *RulesPage) GetName() string {
	return "Rules"
}

func (p *RulesPage) Select() {

}

func (p *RulesPage) Layout(gtx C) D {
	if p.btns.Clicked(0) {
		p.editRule(new(riddle.Rule))
	}
	if p.btns.Clicked(1) {
		p.modal.ModalYesNo("Are you sure?", p.Reset)
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			if len(p.rules) == 0 {
				return D{}
			}
			return p.list.Layout(gtx, p.theme, len(p.rules), func(gtx C, idx int) D {
				dims := p.rules[idx].Layout(gtx, p)
				dims.Size.Y += gtx.Px(unit.Dp(12))
				return dims
			})
		}),
		layout.Rigid(func(gtx C) D {
			return p.btns.Layout(gtx, p.theme)
		}),
	)
}

func (p *RulesPage) HandleSetup(setup riddle.Setup) {
	toBeRemoved := make([]int, 0, len(p.rules))
	removeCount := 0

	for i, rule := range p.rules {
		if err := rule.Check(setup); err != nil {
			toBeRemoved = append(toBeRemoved, i)
		}
	}

	for _, index := range toBeRemoved {
		p.removeRule(index-removeCount, false)
		removeCount++
	}

	if removeCount > 0 {
		p.Save()
	}
}

func (p *RulesPage) GetRules() []riddle.Rule {
	rules := make([]riddle.Rule, 0, len(p.rules))
	for _, rule := range p.rules {
		rules = append(rules, *rule.Rule)
	}
	return rules
}

func (p *RulesPage) SetRules(rules []riddle.Rule) {
	p.Reset()
	for _, rule := range rules {
		heapRule := &riddle.Rule{}
		*heapRule = rule
		p.addRule(heapRule, false)
	}
	p.Save()
}

func (p *RulesPage) SaveRule(rule *riddle.Rule) {
	p.addRule(rule, true)
}

func (p *RulesPage) addRule(rule *riddle.Rule, save bool) {
	if _, found := p.findRule(rule); !found {
		p.rules = append(p.rules, ruleItem{Rule: rule})
	}

	if save {
		p.Save()
	}
}

func (p *RulesPage) findRule(rule *riddle.Rule) (index int, found bool) {
	for i, ruleItem := range p.rules {
		if ruleItem.Rule == rule {
			return i, true
		}
	}
	return -1, false
}

func (p *RulesPage) editRule(rule *riddle.Rule) {
	if p.editFunc != nil {
		p.editFunc(rule)
	}
}

func (p *RulesPage) removeRule(index int, save bool) {
	if len(p.rules) == 0 {
		return
	}

	if index < 0 {
		index = len(p.rules) + index
	}
	if index >= len(p.rules) {
		index = len(p.rules) - 1
	}
	if index < 0 {
		index = 0
	}

	p.rules = append(p.rules[:index], p.rules[index+1:]...)

	if save {
		p.Save()
	}
}

func (p *RulesPage) Save() {
	if p.saveFunc != nil {
		p.saveFunc(p.GetRules())
	}
}

func (p *RulesPage) SetEditFunc(editFunc func(*riddle.Rule)) {
	p.editFunc = editFunc
}

func (p *RulesPage) SetSaveFunc(saveFunc func([]riddle.Rule)) {
	p.saveFunc = saveFunc
}

func (p *RulesPage) Reset() {
	p.rules = nil
	p.Save()
}

type ruleItem struct {
	widget.Clickable
	richtext.InteractiveText
	*riddle.Rule
	deleteBtn widget.Clickable
}

func (rule *ruleItem) Layout(gtx C, p *RulesPage) D {
	th := p.theme
	separator := richtext.SpanStyle{
		Content: " - ",
		Color:   th.Fg,
		Size:    th.TextSize,
	}
	normal := func(text string) richtext.SpanStyle {
		return richtext.SpanStyle{
			Content: text,
			Color:   th.Fg,
			Size:    th.TextSize,
		}
	}
	colored := func(text string) richtext.SpanStyle {
		return richtext.SpanStyle{
			Content: text,
			Color:   th.ContrastBg,
			Size:    th.TextSize,
		}
	}
	itemTypeA, itemValueA := rule.ItemA.Split()
	itemTypeB, itemValueB := rule.ItemB.Split()
	spans := []richtext.SpanStyle{
		normal(itemTypeA + ":"),
		colored(itemValueA),
		separator,
		normal(itemTypeB + ":"),
		colored(itemValueB),
		separator,
		normal(rule.Relation.String()),
	}
	if rule.HasCondition() {
		extraSpans := []richtext.SpanStyle{
			normal(" if A and B is "),
			colored(rule.ConditionItemType),
			normal(" and "),
			colored(rule.Condition),
		}
		spans = append(spans, extraSpans...)
		if rule.IsReversible {
			spans = append(spans, normal(" [reversible]"))
		}
	}

	if rule.Hovered() {
		for i := range spans {
			spans[i].Font.Weight = text.Bold
		}
	}
	if rule.Clicked() {
		p.editRule(rule.Rule)
	}

	deleteBtn := IconAndTextButton(th, &rule.deleteBtn, p.deleteIcon, "")
	deleteBtn.Inset = layout.UniformInset(unit.Dp(1))
	if rule.deleteBtn.Clicked() {
		p.modal.ModalYesNo("Are you sure?", func() {
			idx, _ := p.findRule(rule.Rule)
			p.removeRule(idx, true)
		})
	}

	dims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			dims := richtext.Text(&rule.InteractiveText, th.Shaper, spans...).Layout(gtx)
			gtx.Constraints.Min = dims.Size
			gtx.Constraints.Max = dims.Size
			rule.Clickable.Layout(gtx)
			return dims
		}),
		layout.Rigid(func(gtx C) D {
			return D{Size: image.Pt(gtx.Px(unit.Dp(5)), 0)}
		}),
		layout.Rigid(deleteBtn.Layout),
	)
	dims.Size.X = gtx.Constraints.Max.X
	return dims
}

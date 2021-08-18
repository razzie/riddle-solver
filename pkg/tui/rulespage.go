package tui

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/razzie/riddle-solver/pkg/riddle"
	"github.com/rivo/tview"
)

// RulesPage is a UI element that contains the list of rules
type RulesPage struct {
	Page
	list     *tview.List
	rules    []*riddle.Rule
	editFunc func(*riddle.Rule)
	saveFunc func([]riddle.Rule)
	modal    ModalHandler
}

// NewRuleList returns a new RuleList
func NewRuleList(modal ModalHandler) *RulesPage {
	list := tview.NewList().ShowSecondaryText(false)
	l := &RulesPage{
		Page:  NewPage(tview.NewFrame(list), "Rules"),
		list:  list,
		modal: modal,
	}
	list.SetInputCapture(l.handleInput)
	return l
}

// HandleSetup filters the list based on the new setup
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

// GetRules returns the rules list in a slice
func (p *RulesPage) GetRules() []riddle.Rule {
	rules := make([]riddle.Rule, 0, len(p.rules))
	for _, rule := range p.rules {
		rules = append(rules, *rule)
	}
	return rules
}

// SetRules resets the list and inserts the provided rules
func (p *RulesPage) SetRules(rules []riddle.Rule) {
	p.Reset()
	for _, rule := range rules {
		heapRule := &riddle.Rule{}
		*heapRule = rule
		p.addRule(heapRule, false)
	}
	p.Save()
}

// SaveRule adds a new rule to the list or updates an existing one
func (p *RulesPage) SaveRule(rule *riddle.Rule) {
	p.addRule(rule, true)
}

func (p *RulesPage) addRule(rule *riddle.Rule, save bool) {
	text := fmt.Sprintf("%s - %s - %s",
		colorizeItem(rule.ItemA),
		colorizeItem(rule.ItemB),
		rule.Relation.String())

	if rule.HasCondition() {
		text += fmt.Sprintf(" if A and B is %s and %s",
			colorize(rule.ConditionItemType),
			colorize(rule.Condition))
		if rule.IsReversible {
			text += " [reversible[]"
		}
	}

	selected := func() {
		if p.editFunc != nil {
			p.editFunc(rule)
		}
	}

	index, found := p.findRule(rule)
	if found {
		p.list.RemoveItem(index)
		p.list.InsertItem(index, text, "", 0, selected)
	} else {
		p.rules = append(p.rules, rule)
		p.list.InsertItem(-1, text, "", 0, selected)
	}

	if save {
		p.Save()
	}
}

func (p *RulesPage) findRule(rule *riddle.Rule) (index int, found bool) {
	for i, r := range p.rules {
		if r == rule {
			return i, true
		}
	}
	return -1, false
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
	p.list.RemoveItem(index)

	if save {
		p.Save()
	}
}

// Save calls the save function/callback with the list of rules
func (p *RulesPage) Save() {
	if p.saveFunc != nil {
		p.saveFunc(p.GetRules())
	}
}

func (p *RulesPage) handleInput(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	if key >= tcell.KeyDelete {
		remove := func() {
			p.removeRule(p.list.GetCurrentItem(), true)
		}
		p.modal.ModalYesNo("Do you really want to remove this rule?", remove)
		return nil
	}

	return event
}

// SetEditFunc sets a function that gets called on the selected rule
func (p *RulesPage) SetEditFunc(editFunc func(*riddle.Rule)) {
	p.editFunc = editFunc
}

// SetSaveFunc sets a function that gets the list of all rules upon an update
func (p *RulesPage) SetSaveFunc(saveFunc func([]riddle.Rule)) {
	p.saveFunc = saveFunc
}

// Reset resets the list
func (p *RulesPage) Reset() {
	p.rules = nil
	p.list.Clear()
}

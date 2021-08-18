package tui

import (
	"github.com/razzie/riddle-solver/pkg/riddle"
	"github.com/rivo/tview"
)

// AddRulePage is an input form to enter data for new rules
type AddRulePage struct {
	Page
	form                *tview.Form
	itemA               *tview.InputField
	itemB               *tview.InputField
	relation            *tview.DropDown
	hasCondition        *tview.Checkbox
	conditionExpr       *tview.InputField
	conditionItemType   *tview.InputField
	conditionReversible *tview.Checkbox
	rule                *riddle.Rule
	setup               riddle.Setup
	saveFunc            func(*riddle.Rule)
	modal               ModalHandler
}

// NewRuleForm returns a new RuleForm
func NewRuleForm(modal ModalHandler) *AddRulePage {
	itemA := tview.NewInputField().
		SetLabel("Item A").
		SetPlaceholder("e.g. color:red").
		SetFieldWidth(30)
	itemB := tview.NewInputField().
		SetLabel("Item B").
		SetPlaceholder("e.g. color:red").
		SetFieldWidth(30)
	relation := tview.NewDropDown().
		SetLabel("Relation").
		SetOptions([]string{"associated", "disassociated"}, nil).
		SetCurrentOption(0).
		SetFieldWidth(15)
	conditionItemType := tview.NewInputField().
		SetLabel("Condition item type").
		SetPlaceholder("e.g. position").
		SetFieldWidth(30)
	conditionExpr := tview.NewInputField().
		SetLabel("Condition expression").
		SetPlaceholder("e.g. (A == B - 1) || (A == B + 1)").
		SetFieldWidth(50)
	conditionReversible := tview.NewCheckbox().
		SetLabel("Reversible A <-> B")
	hasCondition := tview.NewCheckbox().
		SetLabel("Condition")
	form := tview.NewForm().
		SetLabelColor(tview.Styles.PrimaryTextColor).
		SetFieldTextColor(tview.Styles.SecondaryTextColor).
		SetButtonTextColor(tview.Styles.SecondaryTextColor).
		AddFormItem(itemA).
		AddFormItem(itemB).
		AddFormItem(relation).
		AddFormItem(hasCondition)

	f := &AddRulePage{
		Page:                NewPage(form, "Add rule"),
		form:                form,
		itemA:               itemA,
		itemB:               itemB,
		relation:            relation,
		hasCondition:        hasCondition,
		conditionExpr:       conditionExpr,
		conditionItemType:   conditionItemType,
		conditionReversible: conditionReversible,
		modal:               modal,
	}
	form.AddButton("Save", f.Save)
	form.AddButton("Reset", f.Reset)
	f.hasCondition.SetChangedFunc(f.showConditionFields)

	return f
}

// HandleSetup configured the autocomplete and dropdown fields
func (p *AddRulePage) HandleSetup(setup riddle.Setup) {
	p.Reset()
	p.setup = setup
	autocompleteItems := getAutocompleteItemsFunc(setup.GetItems())
	autocompleteItemTypes := getAutocompleteFunc(setup.GetItemTypes())
	p.itemA.SetAutocompleteFunc(autocompleteItems)
	p.itemB.SetAutocompleteFunc(autocompleteItems)
	p.conditionItemType.SetAutocompleteFunc(autocompleteItemTypes)
}

// EditRule sets up the form for editing an existing rule
// The given pointer will be supplied to the save function later, unless the user
// resets the form.
func (p *AddRulePage) EditRule(rule *riddle.Rule) {
	p.rule = rule
	p.itemA.SetText(string(rule.ItemA))
	p.itemB.SetText(string(rule.ItemB))
	p.relation.SetCurrentOption(int(rule.Relation))
	p.conditionExpr.SetText(rule.Condition)
	p.conditionItemType.SetText(rule.ConditionItemType)
	p.conditionReversible.SetChecked(rule.IsReversible)
	if p.hasCondition.IsChecked() != rule.HasCondition() {
		p.showConditionFields(!p.hasCondition.IsChecked())
		p.hasCondition.SetChecked(rule.HasCondition())
	}
}

// Save calls the save function on the currently edited or new rule
func (p *AddRulePage) Save() {
	var rule riddle.Rule
	rule.ItemA = riddle.Item(p.itemA.GetText())
	rule.ItemB = riddle.Item(p.itemB.GetText())
	relation, _ := p.relation.GetCurrentOption()
	rule.Relation = riddle.Relation(relation)
	rule.Condition = p.conditionExpr.GetText()
	rule.ConditionItemType = p.conditionItemType.GetText()
	rule.IsReversible = p.conditionReversible.IsChecked()

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

// SetSaveFunc sets a function that gets called on save
func (p *AddRulePage) SetSaveFunc(saveFunc func(*riddle.Rule)) {
	p.saveFunc = saveFunc
}

// Reset resets the form
func (p *AddRulePage) Reset() {
	p.rule = nil
	p.form.SetFocus(0)
	p.itemA.SetText("")
	p.itemB.SetText("")
	p.relation.SetCurrentOption(0)
	p.conditionExpr.SetText("")
	p.conditionItemType.SetText("")
	p.conditionReversible.SetChecked(false)
	if p.hasCondition.IsChecked() {
		p.hasCondition.SetChecked(false)
		p.form.RemoveFormItem(6)
		p.form.RemoveFormItem(5)
		p.form.RemoveFormItem(4)
	}
}

func (p *AddRulePage) showConditionFields(show bool) {
	if show {
		p.form.AddFormItem(p.conditionItemType)
		p.form.AddFormItem(p.conditionExpr)
		p.form.AddFormItem(p.conditionReversible)
	} else {
		p.form.RemoveFormItem(6)
		p.form.RemoveFormItem(5)
		p.form.RemoveFormItem(4)
	}
}

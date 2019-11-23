package ui

import (
	"fmt"

	"github.com/razzie/riddle-solver/riddle"
	"github.com/rivo/tview"
)

// RuleForm is an input form to enter data for new rules
type RuleForm struct {
	tview.Primitive
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
func NewRuleForm(modal ModalHandler) *RuleForm {
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

	f := &RuleForm{
		Primitive:           form,
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
func (f *RuleForm) HandleSetup(setup riddle.Setup) {
	f.Reset()
	f.setup = setup
	autocompleteItems := getAutocompleteItemsFunc(setup.GetItems())
	autocompleteItemTypes := getAutocompleteFunc(setup.GetItemTypes())
	f.itemA.SetAutocompleteFunc(autocompleteItems)
	f.itemB.SetAutocompleteFunc(autocompleteItems)
	f.conditionItemType.SetAutocompleteFunc(autocompleteItemTypes)
}

// EditRule sets up the form for editing an existing rule
// The given pointer will be supplied to the save function later, unless the user
// resets the form.
func (f *RuleForm) EditRule(rule *riddle.Rule) {
	f.rule = rule
	f.itemA.SetText(string(rule.ItemA))
	f.itemB.SetText(string(rule.ItemB))
	f.relation.SetCurrentOption(int(rule.Relation))
	f.conditionExpr.SetText(rule.Condition)
	f.conditionItemType.SetText(rule.ConditionItemType)
	f.conditionReversible.SetChecked(rule.IsReversible)
	if f.hasCondition.IsChecked() != rule.HasCondition() {
		f.showConditionFields(!f.hasCondition.IsChecked())
		f.hasCondition.SetChecked(rule.HasCondition())
	}
}

// Save calls the save function on the currently edited or new rule
func (f *RuleForm) Save() {
	var rule riddle.Rule
	rule.ItemA = riddle.Item(f.itemA.GetText())
	rule.ItemB = riddle.Item(f.itemB.GetText())
	relation, _ := f.relation.GetCurrentOption()
	rule.Relation = riddle.Relation(relation)
	rule.Condition = f.conditionExpr.GetText()
	rule.ConditionItemType = f.conditionItemType.GetText()
	rule.IsReversible = f.conditionReversible.IsChecked()

	if err := rule.Check(f.setup); err != nil {
		f.modal.ModalMessage(fmt.Sprint(err))
		return
	}

	if f.rule == nil {
		f.rule = &rule
	} else {
		f.rule.ItemA = rule.ItemA
		f.rule.ItemB = rule.ItemB
		f.rule.Relation = rule.Relation
		f.rule.Condition = rule.Condition
		f.rule.ConditionItemType = rule.ConditionItemType
		f.rule.IsReversible = rule.IsReversible
	}

	if f.saveFunc != nil {
		f.saveFunc(f.rule)
		f.Reset()
	}
}

// SetSaveFunc sets a function that gets called on save
func (f *RuleForm) SetSaveFunc(saveFunc func(*riddle.Rule)) {
	f.saveFunc = saveFunc
}

// Reset resets the form
func (f *RuleForm) Reset() {
	f.rule = nil
	f.form.SetFocus(0)
	f.itemA.SetText("")
	f.itemB.SetText("")
	f.relation.SetCurrentOption(0)
	f.conditionExpr.SetText("")
	f.conditionItemType.SetText("")
	f.conditionReversible.SetChecked(false)
	if f.hasCondition.IsChecked() {
		f.hasCondition.SetChecked(false)
		f.form.RemoveFormItem(6)
		f.form.RemoveFormItem(5)
		f.form.RemoveFormItem(4)
	}
}

func (f *RuleForm) showConditionFields(show bool) {
	if show {
		f.form.AddFormItem(f.conditionItemType)
		f.form.AddFormItem(f.conditionExpr)
		f.form.AddFormItem(f.conditionReversible)
	} else {
		f.form.RemoveFormItem(6)
		f.form.RemoveFormItem(5)
		f.form.RemoveFormItem(4)
	}
}

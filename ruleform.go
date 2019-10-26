package main

import (
	"github.com/rivo/tview"
)

// RuleForm is an input form to enter data for new rules
type RuleForm struct {
	*tview.Form
	itemA             *tview.InputField
	itemB             *tview.InputField
	relation          *tview.DropDown
	conditionItemType *tview.DropDown
	condition         *tview.InputField
	rule              *Rule
	saveFunc          func(*Rule)
}

// NewRuleForm returns a new RuleForm
func NewRuleForm() *RuleForm {
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
		SetOptions([]string{"associated", "disassociated", "unknown"}, nil).
		SetCurrentOption(2).
		SetFieldWidth(20)
	condition := tview.NewInputField().
		SetLabel("Condition (optional)").
		SetPlaceholder("e.g. (A == B - 1) || (A == B + 1)").
		SetFieldWidth(50)
	conditionItemType := tview.NewInputField().
		SetLabel("Condition item type").
		SetPlaceholder("e.g. position (or leave empty)").
		SetFieldWidth(30)
	form := tview.NewForm().
		SetLabelColor(tview.Styles.PrimaryTextColor).
		AddFormItem(itemA).
		AddFormItem(itemB).
		AddFormItem(relation).
		AddFormItem(conditionItemType).
		AddFormItem(condition)

	f := &RuleForm{Form: form}
	f.AddButton("Save", func() { f.Save() })
	f.AddButton("Reset", func() { f.Reset() })

	return f
}

// HandleSetup configured the autocomplete and dropdown fields
func (f *RuleForm) HandleSetup(setup Setup) {

}

// EditRule sets up the form for editing an existing rule
// The given pointer will be supplied to the save function later, unless the user
// resets the form.
func (f *RuleForm) EditRule(rule *Rule) {

}

// Save calls the save function on the currently edited or new rule
func (f *RuleForm) Save() {

}

// SetSaveFunc sets a function that gets called on save
func (f *RuleForm) SetSaveFunc(saveFunc func(*Rule)) {

}

// Reset resets the form
func (f *RuleForm) Reset() {

}

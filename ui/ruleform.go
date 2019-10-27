package ui

import (
	"fmt"
	"strings"

	"github.com/razzie/riddle-solver/solver"
	"github.com/rivo/tview"
)

// RuleForm is an input form to enter data for new rules
type RuleForm struct {
	*tview.Form
	itemA             *tview.InputField
	itemB             *tview.InputField
	relation          *tview.DropDown
	condition         *tview.InputField
	conditionItemType *tview.InputField
	rule              *solver.Rule
	setup             solver.Setup
	saveFunc          func(*solver.Rule)
	modal             ModalHandler
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

	f := &RuleForm{
		Form:              form,
		itemA:             itemA,
		itemB:             itemB,
		relation:          relation,
		condition:         condition,
		conditionItemType: conditionItemType,
		modal:             modal}
	f.AddButton("Save", func() { f.Save() })
	f.AddButton("Reset", func() { f.Reset() })

	return f
}

func getAutocompleteFunc(words []string) func(string) []string {
	return func(currentText string) (results []string) {
		if len(currentText) == 0 {
			return
		}

		for _, word := range words {
			if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) {
				results = append(results, word)
			}
		}

		if len(results) <= 1 {
			results = nil
		}

		return
	}
}

// HandleSetup configured the autocomplete and dropdown fields
func (f *RuleForm) HandleSetup(setup solver.Setup) {
	f.Reset()
	f.setup = setup
	autocompleteItems := getAutocompleteFunc(setup.GetItems())
	autocompleteItemTypes := getAutocompleteFunc(setup.GetItemTypes())
	f.itemA.SetAutocompleteFunc(autocompleteItems)
	f.itemB.SetAutocompleteFunc(autocompleteItems)
	f.conditionItemType.SetAutocompleteFunc(autocompleteItemTypes)
}

// EditRule sets up the form for editing an existing rule
// The given pointer will be supplied to the save function later, unless the user
// resets the form.
func (f *RuleForm) EditRule(rule *solver.Rule) {
	f.rule = rule
	f.itemA.SetText(rule.ItemA)
	f.itemB.SetText(rule.ItemB)
	f.relation.SetCurrentOption(int(rule.Relation))
	f.condition.SetText(rule.Condition)
	f.conditionItemType.SetText(rule.ConditionItemType)
}

// Save calls the save function on the currently edited or new rule
func (f *RuleForm) Save() {
	var rule solver.Rule
	rule.ItemA = f.itemA.GetText()
	rule.ItemB = f.itemB.GetText()
	relation, _ := f.relation.GetCurrentOption()
	rule.Relation = solver.Relation(relation)
	rule.Condition = f.condition.GetText()
	rule.ConditionItemType = f.conditionItemType.GetText()

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
	}

	if f.saveFunc != nil {
		f.saveFunc(f.rule)
		f.Reset()
	}
}

// SetSaveFunc sets a function that gets called on save
func (f *RuleForm) SetSaveFunc(saveFunc func(*solver.Rule)) {
	f.saveFunc = saveFunc
}

// Reset resets the form
func (f *RuleForm) Reset() {
	f.rule = nil
	f.itemA.SetText("")
	f.itemB.SetText("")
	f.relation.SetCurrentOption(0)
	f.condition.SetText("")
	f.conditionItemType.SetText("")
}

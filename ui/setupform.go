package ui

import (
	"fmt"
	"strings"

	"github.com/razzie/riddle-solver/solver"
	"github.com/rivo/tview"
)

// SetupForm is a form where the user can input riddle item types and values
type SetupForm struct {
	*tview.Form
	saveFunc       func(solver.Setup)
	itemCount      int
	itemTypeFields []*tview.InputField
	valuesFields   []*tview.InputField
	modal          ModalHandler
}

// NewSetupForm returns a new SetupForm
func NewSetupForm(modal ModalHandler) *SetupForm {
	f := &SetupForm{
		Form:  tview.NewForm(),
		modal: modal}
	f.SetLabelColor(tview.Styles.PrimaryTextColor).
		AddButton("Add item type", func() { f.addNewItemType() }).
		AddButton("Save / apply", func() { f.Save() }).
		AddButton("Reset", func() { f.Reset() })
	f.Reset()
	return f
}

func (f *SetupForm) addNewItemType() {
	f.itemCount++

	itemTypeField := tview.NewInputField().
		SetLabel(fmt.Sprintf("#%-2d item type", f.itemCount)).
		SetPlaceholder("e.g. color").
		SetFieldWidth(20)
	f.AddFormItem(itemTypeField)
	f.itemTypeFields = append(f.itemTypeFields, itemTypeField)

	valuesField := tview.NewInputField().
		SetLabel("    values").
		SetPlaceholder("e.g. red, green, blue").
		SetFieldWidth(40)
	f.AddFormItem(valuesField)
	f.valuesFields = append(f.valuesFields, valuesField)
}

// SetSaveFunc sets a function that gets called when data is saved
func (f *SetupForm) SetSaveFunc(saveFunc func(solver.Setup)) {
	f.saveFunc = saveFunc
}

// Save collects all the form data and passes it to the save function
func (f *SetupForm) Save() {
	var setup = make(solver.Setup)

	for i := 0; i < f.itemCount; i++ {
		itemType := f.itemTypeFields[i].GetText()
		if len(itemType) == 0 {
			if len(f.valuesFields[i].GetText()) > 0 {
				f.modal.ModalMessage("Cannot have values without item type")
				return
			}

			continue
		}

		values := strings.Split(f.valuesFields[i].GetText(), ",")
		var trimmedValues []string
		for _, value := range values {
			trimmedValue := strings.Trim(value, " ")
			if len(trimmedValue) == 0 {
				continue
			}
			trimmedValues = append(trimmedValues, trimmedValue)
		}
		if len(trimmedValues) == 0 {
			continue
		}

		setup[itemType] = trimmedValues
	}

	if err := setup.Check(); err != nil {
		f.modal.ModalMessage(fmt.Sprint(err))
		return
	}

	if f.saveFunc != nil {
		f.saveFunc(setup)
	}
}

// Reset resets the form to its initial state
func (f *SetupForm) Reset() {
	f.itemCount = 0
	f.Clear(false)
	f.itemTypeFields = nil
	f.valuesFields = nil
	f.addNewItemType()
	f.addNewItemType()
}

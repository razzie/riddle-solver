package ui

import (
	"fmt"
	"strings"

	"github.com/razzie/riddle-solver/riddle"
	"github.com/rivo/tview"
)

// SetupForm is a form where the user can input riddle item types and values
type SetupForm struct {
	*tview.Form
	saveFunc       func(riddle.Setup)
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
		AddButton("Add item type", func() { f.addItemTypeField() }).
		AddButton("Save / apply", func() { f.Save() }).
		AddButton("Reset", func() { f.Reset() })
	f.Reset()
	return f
}

func (f *SetupForm) addItemTypeField() {
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

// AddItemType adds a new item type field with the provided values or uses an existing empty one
func (f *SetupForm) AddItemType(itemType string, values ...string) {
	for i := 0; i < f.itemCount; i++ {
		if len(f.itemTypeFields[i].GetText()) == 0 && len(f.valuesFields[i].GetText()) == 0 {
			f.itemTypeFields[i].SetText(itemType)
			f.valuesFields[i].SetText(strings.Join(values, ", "))
			return
		}
	}

	f.addItemTypeField()
	f.itemTypeFields[f.itemCount-1].SetText(itemType)
	f.valuesFields[f.itemCount-1].SetText(strings.Join(values, ", "))
}

// SetSaveFunc sets a function that gets called when data is saved
func (f *SetupForm) SetSaveFunc(saveFunc func(riddle.Setup)) {
	f.saveFunc = saveFunc
}

// Save collects all the form data and passes it to the save function
func (f *SetupForm) Save() {
	var setup = make(riddle.Setup)

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
	f.addItemTypeField()
	f.addItemTypeField()

	if f.saveFunc != nil {
		f.saveFunc(nil)
	}
}

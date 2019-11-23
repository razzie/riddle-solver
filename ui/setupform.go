package ui

import (
	"fmt"
	"strings"

	"github.com/razzie/riddle-solver/riddle"
	"github.com/rivo/tview"
)

// SetupForm is a form where the user can input riddle item types and values
type SetupForm struct {
	tview.Primitive
	form           *tview.Form
	saveFunc       func(riddle.Setup)
	itemCount      int
	itemTypeFields []*tview.InputField
	valuesFields   []*tview.InputField
	modal          ModalHandler
}

// NewSetupForm returns a new SetupForm
func NewSetupForm(modal ModalHandler) *SetupForm {
	form := tview.NewForm()
	f := &SetupForm{
		Primitive: form,
		form:      form,
		modal:     modal,
	}
	form.SetLabelColor(tview.Styles.PrimaryTextColor).
		SetFieldTextColor(tview.Styles.SecondaryTextColor).
		SetButtonTextColor(tview.Styles.SecondaryTextColor).
		AddButton("Add item type", f.addItemTypeField).
		AddButton("Save / apply", f.Save).
		AddButton("Reset", f.Reset)
	f.Reset()
	return f
}

func (f *SetupForm) addItemTypeField() {
	f.itemCount++

	itemTypeField := tview.NewInputField().
		SetLabel(fmt.Sprintf("#%-2d item type", f.itemCount)).
		SetPlaceholder("e.g. color").
		SetFieldWidth(20)
	f.form.AddFormItem(itemTypeField)
	f.itemTypeFields = append(f.itemTypeFields, itemTypeField)

	valuesField := tview.NewInputField().
		SetLabel("    values").
		SetPlaceholder("e.g. red, green, blue").
		SetFieldWidth(40)
	f.form.AddFormItem(valuesField)
	f.valuesFields = append(f.valuesFields, valuesField)
}

func (f *SetupForm) addItemType(itemType string, values ...string) {
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

// AddItemType adds a new item type field with the provided values or uses an existing empty one
func (f *SetupForm) AddItemType(itemType string, values ...string) {
	f.addItemType(itemType, values...)
	f.Save()
}

// SetSaveFunc sets a function that gets called when data is saved
func (f *SetupForm) SetSaveFunc(saveFunc func(riddle.Setup)) {
	f.saveFunc = saveFunc
}

// GetSetup returns the current setup
func (f *SetupForm) GetSetup() (riddle.Setup, error) {
	var setup = make(riddle.Setup)

	for i := 0; i < f.itemCount; i++ {
		itemType := f.itemTypeFields[i].GetText()
		if len(itemType) == 0 {
			if len(f.valuesFields[i].GetText()) > 0 {
				return nil, fmt.Errorf("Cannot have values without item type")
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
		return nil, err
	}

	return setup, nil
}

// SetSetup resets the form and inserts the values from the provided setup
func (f *SetupForm) SetSetup(setup riddle.Setup) {
	f.Reset()
	for itemType, values := range setup {
		f.addItemType(itemType, values...)
	}
	f.Save()
}

// Save collects all the form data and passes it to the save function
func (f *SetupForm) Save() {
	setup, err := f.GetSetup()
	if err != nil {
		f.modal.ModalMessage(fmt.Sprint(err))
		return
	}

	if f.saveFunc != nil {
		f.saveFunc(setup)
	}
}

// Reset resets the form to its initial state
func (f *SetupForm) Reset() {
	f.form.SetFocus(0)
	f.itemCount = 0
	f.form.Clear(false)
	f.itemTypeFields = nil
	f.valuesFields = nil
	f.addItemTypeField()
	f.addItemTypeField()

	if f.saveFunc != nil {
		f.saveFunc(nil)
	}
}

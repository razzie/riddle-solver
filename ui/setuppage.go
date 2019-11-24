package ui

import (
	"fmt"
	"strings"

	"github.com/razzie/riddle-solver/riddle"
	"github.com/rivo/tview"
)

// SetupPage is a form where the user can input riddle item types and values
type SetupPage struct {
	Page
	form           *tview.Form
	saveFunc       func(riddle.Setup)
	itemCount      int
	itemTypeFields []*tview.InputField
	valuesFields   []*tview.InputField
	modal          ModalHandler
}

// NewSetupForm returns a new SetupForm
func NewSetupForm(modal ModalHandler) *SetupPage {
	form := tview.NewForm()
	f := &SetupPage{
		Page:  NewPage(form, "Setup"),
		form:  form,
		modal: modal,
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

func (p *SetupPage) addItemTypeField() {
	p.itemCount++

	itemTypeField := tview.NewInputField().
		SetLabel(fmt.Sprintf("#%-2d item type", p.itemCount)).
		SetPlaceholder("e.g. color").
		SetFieldWidth(20)
	p.form.AddFormItem(itemTypeField)
	p.itemTypeFields = append(p.itemTypeFields, itemTypeField)

	valuesField := tview.NewInputField().
		SetLabel("    values").
		SetPlaceholder("e.g. red, green, blue").
		SetFieldWidth(40)
	p.form.AddFormItem(valuesField)
	p.valuesFields = append(p.valuesFields, valuesField)
}

func (p *SetupPage) addItemType(itemType string, values ...string) {
	for i := 0; i < p.itemCount; i++ {
		if len(p.itemTypeFields[i].GetText()) == 0 && len(p.valuesFields[i].GetText()) == 0 {
			p.itemTypeFields[i].SetText(itemType)
			p.valuesFields[i].SetText(strings.Join(values, ", "))
			return
		}
	}

	p.addItemTypeField()
	p.itemTypeFields[p.itemCount-1].SetText(itemType)
	p.valuesFields[p.itemCount-1].SetText(strings.Join(values, ", "))
}

// AddItemType adds a new item type field with the provided values or uses an existing empty one
func (p *SetupPage) AddItemType(itemType string, values ...string) {
	p.addItemType(itemType, values...)
	p.Save()
}

// SetSaveFunc sets a function that gets called when data is saved
func (p *SetupPage) SetSaveFunc(saveFunc func(riddle.Setup)) {
	p.saveFunc = saveFunc
}

// GetSetup returns the current setup
func (p *SetupPage) GetSetup() (riddle.Setup, error) {
	var setup = make(riddle.Setup)

	for i := 0; i < p.itemCount; i++ {
		itemType := p.itemTypeFields[i].GetText()
		if len(itemType) == 0 {
			if len(p.valuesFields[i].GetText()) > 0 {
				return nil, fmt.Errorf("Cannot have values without item type")
			}

			continue
		}

		values := strings.Split(p.valuesFields[i].GetText(), ",")
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
func (p *SetupPage) SetSetup(setup riddle.Setup) {
	p.Reset()
	for itemType, values := range setup {
		p.addItemType(itemType, values...)
	}
	p.Save()
}

// Save collects all the form data and passes it to the save function
func (p *SetupPage) Save() {
	setup, err := p.GetSetup()
	if err != nil {
		p.modal.ModalMessage(fmt.Sprint(err))
		return
	}

	if p.saveFunc != nil {
		p.saveFunc(setup)
	}
}

// Reset resets the form to its initial state
func (p *SetupPage) Reset() {
	p.form.SetFocus(0)
	p.itemCount = 0
	p.form.Clear(false)
	p.itemTypeFields = nil
	p.valuesFields = nil
	p.addItemTypeField()
	p.addItemTypeField()

	if p.saveFunc != nil {
		p.saveFunc(nil)
	}
}

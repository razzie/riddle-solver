package ui

import (
	"fmt"
	"strings"

	"github.com/razzie/riddle-solver/solver"
	"github.com/rivo/tview"
)

// SetupPage is a form where the user can input riddle item types and values
type SetupPage struct {
	*tview.Form
	saveFunc       func(solver.Setup)
	itemCount      int
	itemTypeFields []*tview.InputField
	valuesFields   []*tview.InputField
	modal          ModalHandler
}

// NewSetupPage returns a new SetupPage
func NewSetupPage(modal ModalHandler) *SetupPage {
	p := &SetupPage{
		Form:  tview.NewForm(),
		modal: modal}
	p.SetLabelColor(tview.Styles.PrimaryTextColor).
		AddButton("Add item type", func() { p.addNewItemType() }).
		AddButton("Save / apply", func() { p.Save() }).
		AddButton("Reset", func() { p.Reset() })
	p.Reset()
	return p
}

func (p *SetupPage) addNewItemType() {
	p.itemCount++

	itemTypeField := tview.NewInputField().
		SetLabel(fmt.Sprintf("#%-2d item type", p.itemCount)).
		SetPlaceholder("e.g. color").
		SetFieldWidth(20)
	p.AddFormItem(itemTypeField)
	p.itemTypeFields = append(p.itemTypeFields, itemTypeField)

	valuesField := tview.NewInputField().
		SetLabel("    values").
		SetPlaceholder("e.g. red, green, blue").
		SetFieldWidth(40)
	p.AddFormItem(valuesField)
	p.valuesFields = append(p.valuesFields, valuesField)
}

// SetSaveFunc sets a function that gets called when data is saved
func (p *SetupPage) SetSaveFunc(saveFunc func(solver.Setup)) {
	p.saveFunc = saveFunc
}

// Save collects all the form data and passes it to the save function
func (p *SetupPage) Save() {
	var setup = make(solver.Setup)

	for i := 0; i < p.itemCount; i++ {
		itemType := p.itemTypeFields[i].GetText()
		if len(itemType) == 0 {
			if len(p.valuesFields[i].GetText()) > 0 {
				p.modal.ModalMessage("Cannot have values without item type")
				return
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
		p.modal.ModalMessage(fmt.Sprint(err))
		return
	}

	if p.saveFunc != nil {
		p.saveFunc(setup)
	}
}

// Reset resets the form to its initial state
func (p *SetupPage) Reset() {
	p.itemCount = 0
	p.Clear(false)
	p.itemTypeFields = nil
	p.valuesFields = nil
	p.addNewItemType()
	p.addNewItemType()
}

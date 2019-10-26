package main

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

// SetupPage is a form where the user can input riddle item types and values
type SetupPage struct {
	*tview.Form
	delegate       func(Setup)
	itemCount      int
	itemTypeFields []*tview.InputField
	valuesFields   []*tview.InputField
}

// NewSetupPage returns a new SetupPage
func NewSetupPage() *SetupPage {
	p := &SetupPage{Form: tview.NewForm()}
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

// SetSaveDelegate sets a delegate that gets called when data is saved
func (p *SetupPage) SetSaveDelegate(delegate func(Setup)) {
	p.delegate = delegate
}

// Save collects all the form data and passes it to the delegate
func (p *SetupPage) Save() {
	var data = make(Setup)

	for i := 0; i < p.itemCount; i++ {
		itemType := p.itemTypeFields[i].GetText()
		if len(itemType) == 0 {
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

		data[itemType] = trimmedValues
	}

	if p.delegate != nil {
		p.delegate(data)
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

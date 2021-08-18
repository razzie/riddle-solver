package tui

import (
	"fmt"

	"github.com/razzie/riddle-solver/pkg/riddle"
	"github.com/rivo/tview"
)

// SavePage is UI form element for saving the current riddle
type SavePage struct {
	Page
	form   *tview.Form
	name   *tview.InputField
	getter func() (*riddle.Riddle, error)
	modal  ModalHandler
}

// NewSavePage returns a new SavePage
func NewSavePage(modal ModalHandler) *SavePage {
	form := tview.NewForm().
		SetLabelColor(tview.Styles.PrimaryTextColor).
		SetFieldTextColor(tview.Styles.SecondaryTextColor).
		SetButtonTextColor(tview.Styles.SecondaryTextColor)

	name := tview.NewInputField().
		SetLabel("Riddle name").
		SetPlaceholder("riddle name without .json").
		SetFieldWidth(30)

	p := &SavePage{
		Page:  NewPage(form, "Save"),
		form:  form,
		name:  name,
		modal: modal,
	}
	p.form.AddFormItem(name).AddButton("Save", p.Save)
	return p
}

// SetRiddleGetter sets the function that returns the current riddle for saving
func (p *SavePage) SetRiddleGetter(getter func() (*riddle.Riddle, error)) {
	p.getter = getter
}

// Save saves the current riddle
func (p *SavePage) Save() {
	name := p.name.GetText()
	if len(name) == 0 {
		p.modal.ModalMessage("Empty riddle name")
		return
	}

	r, err := p.getter()
	if err != nil {
		p.modal.ModalMessage(err.Error())
		return
	}

	err = r.SaveToFile(fmt.Sprintf("riddles/%s.json", name))
	if err != nil {
		p.modal.ModalMessage(err.Error())
		return
	}

	p.modal.ModalMessage("Riddle saved")
}

package gui

import (
	"sync"

	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	iconsOnce sync.Once
	icns      *Icons
)

type Icons struct {
	ContentAdd            *widget.Icon
	ActionDelete          *widget.Icon
	CheckBoxBlank         *widget.Icon
	CheckBoxIndeterminate *widget.Icon
}

func GetIcons() *Icons {
	iconsOnce.Do(func() {
		icns = new(Icons)
		icns.ContentAdd, _ = widget.NewIcon(icons.ContentAdd)
		icns.ActionDelete, _ = widget.NewIcon(icons.ActionDelete)
		icns.CheckBoxBlank, _ = widget.NewIcon(icons.ToggleCheckBoxOutlineBlank)
		icns.CheckBoxIndeterminate, _ = widget.NewIcon(icons.ToggleIndeterminateCheckBox)
	})
	return icns
}

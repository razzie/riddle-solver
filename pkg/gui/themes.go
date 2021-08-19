package gui

import (
	"gioui.org/font/gofont"
	"gioui.org/widget/material"
)

var Themes = map[string]*material.Theme{
	"light": LightTheme,
}

var LightTheme = material.NewTheme(gofont.Collection())

package tui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Theme contains the colors used by the UI
type Theme struct {
	TextColor                 string
	FieldTextColor            string
	HighlightTextColor        string
	FieldPlaceholderTextColor string
	BackgroundColor           string
	FieldBackgroundColor      string
	DropdownBackgroundColor   string
}

// Apply applies the theme
func (theme *Theme) Apply() {
	currentTheme = theme

	tview.Styles.PrimaryTextColor = tcell.GetColor(theme.TextColor)
	tview.Styles.SecondaryTextColor = tcell.GetColor(theme.FieldTextColor)
	tview.Styles.ContrastSecondaryTextColor = tcell.GetColor(theme.FieldPlaceholderTextColor)
	tview.Styles.InverseTextColor = tcell.GetColor(theme.HighlightTextColor)

	tview.Styles.PrimitiveBackgroundColor = tcell.GetColor(theme.BackgroundColor)
	tview.Styles.ContrastBackgroundColor = tcell.GetColor(theme.FieldBackgroundColor)
	tview.Styles.MoreContrastBackgroundColor = tcell.GetColor(theme.DropdownBackgroundColor)

	tview.Styles.BorderColor = tview.Styles.PrimaryTextColor
	tview.Styles.TitleColor = tview.Styles.PrimaryTextColor
	tview.Styles.GraphicsColor = tview.Styles.PrimaryTextColor
}

var currentTheme = &LightTheme

// Themes contains the string-to-theme map
var Themes = map[string]*Theme{
	"dark":  &DarkTheme,
	"light": &LightTheme,
}

// DarkTheme represents a dark theme
var DarkTheme = Theme{
	TextColor:                 "white",
	FieldTextColor:            "white",
	HighlightTextColor:        "yellow",
	FieldPlaceholderTextColor: "green",
	BackgroundColor:           "black",
	FieldBackgroundColor:      "blue",
	DropdownBackgroundColor:   "grey",
}

// LightTheme represents a light theme
var LightTheme = Theme{
	TextColor:                 "black",
	FieldTextColor:            "black",
	HighlightTextColor:        "red",
	FieldPlaceholderTextColor: "white",
	BackgroundColor:           "white",
	FieldBackgroundColor:      "lightgrey",
	DropdownBackgroundColor:   "red",
}

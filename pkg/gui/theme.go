package gui

type Theme struct{}

var Themes = map[string]*Theme{
	"dark":  &DarkTheme,
	"light": &LightTheme,
}

var DarkTheme = Theme{}
var LightTheme = Theme{}

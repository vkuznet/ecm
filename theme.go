package main

import (
	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var GreyStyle = tview.Theme{
	PrimitiveBackgroundColor:    tcell.ColorGrey,
	ContrastBackgroundColor:     tcell.ColorSilver,
	MoreContrastBackgroundColor: tcell.ColorGreen,
	BorderColor:                 tcell.ColorWhite,
	TitleColor:                  tcell.ColorWhite,
	GraphicsColor:               tcell.ColorWhite,
	PrimaryTextColor:            tcell.ColorBlack,
	SecondaryTextColor:          tcell.ColorYellow,
	TertiaryTextColor:           tcell.ColorGreen,
	InverseTextColor:            tcell.ColorBlue,
	ContrastSecondaryTextColor:  tcell.ColorDarkCyan,
}

// helper function to setup the theme for our app
func setTheme(theme string) {
	if theme == "grey" {
		tview.Styles = GreyStyle
	}
}

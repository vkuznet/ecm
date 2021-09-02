package main

import (
	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// GreyStyle represents grey style theme
var GreyStyle = tview.Theme{
	PrimitiveBackgroundColor:    tcell.ColorLightGrey,
	ContrastBackgroundColor:     tcell.ColorSilver,
	MoreContrastBackgroundColor: tcell.ColorGreen,
	//     BorderColor:                 tcell.ColorMediumOrchid,
	BorderColor:                tcell.ColorSnow,
	TitleColor:                 tcell.ColorSteelBlue,
	GraphicsColor:              tcell.ColorWhite,
	PrimaryTextColor:           tcell.ColorBlack,
	SecondaryTextColor:         tcell.ColorSteelBlue,
	TertiaryTextColor:          tcell.ColorGreen,
	InverseTextColor:           tcell.ColorBlue,
	ContrastSecondaryTextColor: tcell.ColorDarkCyan,
}

// TitleColor represents title color used in widgets
var TitleColor = tcell.ColorWhite

// helper function to setup the theme for our app
func setTheme(theme string) {
	if theme == "grey" {
		tview.Styles = GreyStyle
		TitleColor = GreyStyle.TitleColor
	}
}

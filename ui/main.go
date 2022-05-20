package main

import (
	"image/color"
	"log"

	fyne "fyne.io/fyne/v2"
	app "fyne.io/fyne/v2/app"
	theme "fyne.io/fyne/v2/theme"
)

// various sizes of our application
var windowSize, inputSize fyne.Size
var fontColor color.Color
var mobile bool

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	a := app.NewWithID("io.github.vkuznet")
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("ECM")

	appSettings()
	LoginWindow(a, w)

	w.ShowAndRun()
}

// helper function to set application settings
func appSettings() {
	if mobile {
		windowSize = fyne.NewSize(100, 400)
		inputSize = fyne.NewSize(50, 50)
		fontColor = color.White
	} else {
		windowSize = fyne.NewSize(700, 400)
		inputSize = fyne.NewSize(300, 50)
		fontColor = color.White
	}
}

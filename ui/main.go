package main

import (
	"fmt"
	"log"
	"os"

	fyne "fyne.io/fyne/v2"
	app "fyne.io/fyne/v2/app"
	theme "fyne.io/fyne/v2/theme"
)

// various sizes of our application
var windowSize, inputSize, rowSize fyne.Size
var appKind string

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	a := app.NewWithID("io.github.vkuznet")
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("ECM")

	appSettings(a)
	LoginWindow(a, w)

	w.ShowAndRun()
}

// helper function to set application preferences/settings
func appSettings(app fyne.App) {
	// set some values for our app preferences
	pref := app.Preferences()

	// default values
	cipher := "aes"
	vdir := fmt.Sprintf("%s/.ecm/Primary", os.Getenv("HOME"))
	fontSize := "Normal"

	// rowSize represents main row size used in our UI containers
	rowSize = fyne.NewSize(340, 40)

	// make changes depending on application kind
	if appKind == "desktop" {
		windowSize = fyne.NewSize(900, 600)
		inputSize = fyne.NewSize(300, 50)
	} else {
		vdir = app.Storage().RootURI().Path()
		windowSize = fyne.NewSize(300, 600)
		inputSize = fyne.NewSize(50, 50)
	}

	// save preferences
	pref.SetString("VaultCipher", cipher)
	pref.SetString("VaultDirectory", vdir)
	pref.SetString("FontSize", fontSize)

	// write ecmconfig
	WriteSyncConfig(app)
}

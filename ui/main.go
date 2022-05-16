package main

import (
	"log"

	fyne "fyne.io/fyne/v2"
	app "fyne.io/fyne/v2/app"
	theme "fyne.io/fyne/v2/theme"
)

// various sizes of our application
var windowSize, inputSize fyne.Size

func main() {
	// use New method for generic app
	//     a := app.New()
	// use NewWithID for preferences
	a := app.NewWithID("io.github.vkuznet")
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("ECM application")

	// TODO: make configuration
	mobile := false
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if mobile {
		windowSize = fyne.NewSize(700, 400)
		inputSize = fyne.NewSize(300, 50)
	} else {
		windowSize = fyne.NewSize(700, 400)
		inputSize = fyne.NewSize(300, 50)
	}

	// original idea
	//     MainWindow(a, w)
	//     w.ShowAndRun()

	LoginWindow(a, w)
	w.ShowAndRun()

	//     w.SetContent(Create(a, w))
	//     w.Resize(fyne.NewSize(700, 400))
	//     w.SetMaster()
	//     w.ShowAndRun()
}

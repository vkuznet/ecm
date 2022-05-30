package main

import (
	"log"

	app "fyne.io/fyne/v2/app"
	theme "fyne.io/fyne/v2/theme"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	a := app.NewWithID("io.github.vkuznet")
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("ECM")

	appSettings(a)
	LoginWindow(a, w)

	w.ShowAndRun()
}

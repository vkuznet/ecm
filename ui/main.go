package main

import (
	"context"
	"log"

	app "fyne.io/fyne/v2/app"
	theme "fyne.io/fyne/v2/theme"
)

func main() {
	// set our logging settings
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// setup fyne app and main window
	a := app.NewWithID("io.github.vkuznet.ecm")
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("ECM")

	// setup sync config
	syncConfig(a)

	// load custom images
	setCustomImages()

	// setup app error handler
	setupAppError()

	// setup app settings
	appSettings(a)

	// start login window
	LoginWindow(a, w)

	// start autologout loop
	ctx, cancel := context.WithCancel(context.Background())
	go autoLogout(a, w, ctx)
	defer cancel() // when we quit our app cancel() will be called and quite our goroutine

	// start our app
	w.ShowAndRun()
}

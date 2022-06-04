package main

import (
	"context"
	"log"

	app "fyne.io/fyne/v2/app"
	theme "fyne.io/fyne/v2/theme"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	a := app.NewWithID("io.github.vkuznet")
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("ECM")

	setupAppError()
	appSettings(a)
	LoginWindow(a, w)

	// start autoLogout goroutine which can be killed gracefully
	//     done := make(chan os.Signal, 1)
	//     signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// start autologout loop
	//     go autoLogout(a, w, done)
	ctx, cancel := context.WithCancel(context.Background())
	go autoLogout(a, w, ctx)
	defer cancel() // when we quit our app cancel() will be called and quite our goroutine

	// start our app
	w.ShowAndRun()

	//     <-done
	// here our app will exit since autoLogout goroutine will be gone
}

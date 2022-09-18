package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	app "fyne.io/fyne/v2/app"
	theme "fyne.io/fyne/v2/theme"
)

// version of the code
var gitVersion, gitTag string

// ecmInfo function returns version string of the server
func ecmInfo() string {
	goVersion := runtime.Version()
	tstamp := time.Now().Format("2006-02-01")
	return fmt.Sprintf("ecm git=%s tag=%s go=%s date=%s", gitVersion, gitTag, goVersion, tstamp)
}

func main() {
	var config string
	flag.StringVar(&config, "config", "", "config file name")
	var version bool
	flag.BoolVar(&version, "version", false, "show version")
	var prefs bool
	flag.BoolVar(&prefs, "prefs", false, "show prefs")
	flag.Parse()
	if version {
		fmt.Println(ecmInfo())
		os.Exit(0)

	}
	if config != "" {
		err := ParseConfig(config)
		if err != nil {
			log.Fatal(err)
		}
	}

	// setup fyne app and main window
	a := app.NewWithID("io.github.vkuznet")
	if prefs {
		printPrefs(a)
		os.Exit(0)
	}
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

	// init cloud provider
	initDropbox()

	// start internal web server on non-desktop app
	//     if appKind != "desktop" {
	if appKind != "BLA" {
		ctx, cancel := context.WithCancel(context.Background())
		go authServer(a, ctx)
		defer cancel() // when we quit our app cancel() will be called and quite our goroutine
	}

	// start our app
	w.ShowAndRun()
}

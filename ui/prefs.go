package main

import (
	"fmt"
	"os"

	fyne "fyne.io/fyne/v2"
	canvas "fyne.io/fyne/v2/canvas"
)

// various sizes of our application
var windowSize, inputSize, rowSize fyne.Size
var appKind, appTheme string
var gitImage, docImage, webImage, lockImage *canvas.Image

// helper function to set application preferences/settings
func appSettings(app fyne.App) {
	// set some values for our app preferences
	pref := app.Preferences()

	// default values
	cipher := "aes"
	vdir := fmt.Sprintf("%s/.ecm/Primary", os.Getenv("HOME"))
	fontSize := "Normal"
	appTheme = pref.String("AppTheme")
	if appTheme == "" {
		appTheme = "dark"
	}

	// rowSize represents main row size used in our UI containers
	rowSize = fyne.NewSize(340, 40)
	// input size represents size of the input field which is shorter by row size
	inputSize = fyne.NewSize(340, 40)

	// make changes depending on application kind
	if appKind == "desktop" {
		windowSize = fyne.NewSize(900, 600)
	} else {
		vdir = app.Storage().RootURI().Path()
		windowSize = fyne.NewSize(300, 600)
		// on mobile input size should be short of rowSize
		inputSize = fyne.NewSize(300, 40)
	}

	// save preferences
	pref.SetString("VaultCipher", cipher)
	pref.SetString("VaultDirectory", vdir)
	pref.SetString("FontSize", fontSize)
	pref.SetString("AppTheme", appTheme)

	// set images
	setCustomImages()

	// write ecmconfig
	WriteSyncConfig(app)
}

// helper function to set custom images for our app theme
func setCustomImages() {
	if appTheme == "light" {
		gitImage = canvas.NewImageFromResource(resourceGithubBlackSvg)
		webImage = canvas.NewImageFromResource(resourceWebBlackSvg)
		docImage = canvas.NewImageFromResource(resourceDocBlackSvg)
		lockImage = canvas.NewImageFromResource(resourceLockBlackSvg)
	} else {
		gitImage = canvas.NewImageFromResource(resourceGithubWhiteSvg)
		webImage = canvas.NewImageFromResource(resourceWebWhiteSvg)
		docImage = canvas.NewImageFromResource(resourceDocWhiteSvg)
		lockImage = canvas.NewImageFromResource(resourceLockWhiteSvg)
	}
	gitImage.SetMinSize(fyne.NewSize(40, 40))
	webImage.SetMinSize(fyne.NewSize(35, 35))
	docImage.SetMinSize(fyne.NewSize(40, 40))
	lockImage.SetMinSize(fyne.NewSize(40, 40))
}

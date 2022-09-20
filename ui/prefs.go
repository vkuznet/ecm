package main

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	fyne "fyne.io/fyne/v2"
	theme "fyne.io/fyne/v2/theme"
)

// various sizes of our application
var windowSize, inputSize, rowSize fyne.Size
var appKind, appTheme, fontSize string

// var gitImage, docImage, webImage, lockImage, syncImage, passImage, listImage *canvas.Image
// var rightArrowImage, leftArrowImage *canvas.Image
var btnColor, greenColor, blueColor, redColor, authColor, grayColor color.NRGBA

// helper function to set application preferences/settings
// the app preference file has the following attributes
// AppTheme: dark or light
// FontSize: Normal
// VaultCipher: aes
// VaultDirectory: /path/.ecm/Primary
// VaultName: Primary
// cloud: dropbox:ECM
// local: http://...
func appSettings(app fyne.App) {
	// set some values for our app preferences
	pref := app.Preferences()

	// default values
	cipher := getPrefValue(pref, "VaultCipher", "aes")
	vdir := getPrefValue(pref, "VaultDirectory", fmt.Sprintf("%s/.ecm/Primary", os.Getenv("HOME")))
	if fontSize == "" {
		fontSize = getPrefValue(pref, "FontSize", "Normal")
	}
	if appTheme == "" {
		appTheme = getPrefValue(pref, "AppTheme", "dark")
	}

	// color for our buttons
	redColor = color.NRGBA{0xff, 0x26, 0x00, 0xff}
	authColor = color.NRGBA{0x94, 0x17, 0x51, 0xff}
	btnColor = color.NRGBA{0x79, 0x79, 0x79, 0xff}
	greenColor = color.NRGBA{0x00, 0x8f, 0x00, 0xff}
	blueColor = color.NRGBA{0x04, 0x33, 0xff, 0xff}
	grayColor = color.NRGBA{0xc0, 0xc0, 0xc0, 0xff}

	// rowSize represents main row size used in our UI containers
	rowSize = fyne.NewSize(340, 40)
	// input size represents size of the input field which is shorter by row size
	inputSize = fyne.NewSize(340, 40)

	// default vault name
	vname := "Primary"

	// make changes depending on application kind
	if appKind == "desktop" {
		windowSize = fyne.NewSize(900, 600)
		// on desktop app we'll use name of vault dir
		arr := strings.Split(vdir, "/")
		if len(arr) > 1 {
			vname = arr[len(arr)-1]
		}
	} else {
		vdir = app.Storage().RootURI().Path()
		windowSize = fyne.NewSize(300, 600)
		// on mobile input size should be short of rowSize
		inputSize = fyne.NewSize(300, 40)
	}

	// save preferences
	pref.SetString("VaultCipher", cipher)
	pref.SetString("VaultDirectory", vdir)
	pref.SetString("VaultName", vname)
	pref.SetString("FontSize", fontSize)
	pref.SetString("AppTheme", appTheme)
	if autoThreshold != nil {
		thr, err := autoThreshold.Get()
		if err == nil {
			pref.SetString("Autologout", thr)
		}
	}

	// set images
	//     setCustomImages()
}

// helper function to set custom images for our app theme
func setCustomImages() {
	theme.NewThemedResource(resourceGithubSvg)
	theme.NewThemedResource(resourceWebSvg)
	theme.NewThemedResource(resourceDocSvg)
	theme.NewThemedResource(resourceLockSvg)
	theme.NewThemedResource(resourceSyncSvg)
	theme.NewThemedResource(resourcePassSvg)
	theme.NewThemedResource(resourceListSvg)
	theme.NewThemedResource(resourceLeftArrowSvg)
	theme.NewThemedResource(resourceRightArrowSvg)
}

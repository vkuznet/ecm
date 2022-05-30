package main

import (
	"log"

	"fyne.io/fyne/v2"
	container "fyne.io/fyne/v2/container"
	widget "fyne.io/fyne/v2/widget"
	vt "github.com/vkuznet/ecm/vault"
)

// define our vault
var _vault *vt.Vault

// ErrWindow represents generic error window content
func ErrWindow(app fyne.App, w fyne.Window, err error) {
	log.Println("ERROR", err)
	w.SetContent(widget.NewLabel("Error window"))
}

// AppWindow represents application window
func AppWindow(app fyne.App, w fyne.Window) {
	w.SetContent(Create(app, w))
	w.Resize(windowSize)
	w.SetMaster()
	// custom theme
	//     app.Settings().SetTheme(&grayTheme{})
}

// Create will stitch together all ui components
func Create(app fyne.App, window fyne.Window) *container.AppTabs {
	uiRecords := newUIVaultRecords(app, window)
	return &container.AppTabs{Items: []*container.TabItem{
		uiRecords.tabItem(),
		newUIRecord(app, window).tabItem(),
		newUIPassword(app, window).tabItem(),
		newUISync(app, window, uiRecords).tabItem(),
		newUISettings(app, window).tabItem(),
	}}
}

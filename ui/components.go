package main

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	container "fyne.io/fyne/v2/container"
	binding "fyne.io/fyne/v2/data/binding"
	theme "fyne.io/fyne/v2/theme"
	widget "fyne.io/fyne/v2/widget"
	vt "github.com/vkuznet/ecm/vault"
)

// define our vault
var _vault *vt.Vault

// global app tabs which keeps all app tabs
var appTabs *container.AppTabs

// global app error string
var appError binding.String

// global label for error widget
var appErrorLabel *widget.Label

// AppWindow represents application window
func AppWindow(app fyne.App, w fyne.Window) {
	w.SetContent(Create(app, w))
	w.Resize(windowSize)
	w.SetMaster()
	// custom theme
	//     app.Settings().SetTheme(&grayTheme{})
}

// helper function to unify error messages
func errorMessage(msg string, err error) {
	tstamp := time.Now().Format(time.RFC3339)
	text := fmt.Sprintf("%s %s, error: %v", tstamp, msg, err)
	log.Println("###", text, "appError", appError)
	if appError == nil {
		appError = binding.NewString()
	}
	appError.Set(text)
}

// helper function to setup app error label
func setupAppError() {
	appError = binding.NewString()
	appError.Set("ECM error window")
	appErrorLabel = widget.NewLabelWithData(appError)
}

// errorTabItem continer
func errorTabItem(app fyne.App, w fyne.Window) *container.TabItem {
	content := container.NewVBox(
		appErrorLabel,
	)
	return &container.TabItem{Text: "Error", Icon: theme.ErrorIcon(), Content: content}
}

// Create will stitch together all ui components
func Create(app fyne.App, window fyne.Window) *container.AppTabs {
	uiRecords := newUIVaultRecords(app, window)
	appTabs = &container.AppTabs{Items: []*container.TabItem{
		uiRecords.tabItem(),
		newUIRecord(app, window).tabItem(),
		newUIPassword(app, window).tabItem(),
		newUISync(app, window, uiRecords).tabItem(),
		newUISettings(app, window).tabItem(),
		errorTabItem(app, window),
		logoutTabItem(app, window),
	}}
	return appTabs
}

// helper function to refresh all app widgets/canvases
func appRefresh(app fyne.App, window fyne.Window) {
	// update custom image pointers
	setCustomImages()
	// refresh all widgets
	content := window.Content()
	canvas := window.Canvas()
	canvas.Refresh(content)
}

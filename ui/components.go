package main

import (
	"fmt"
	"log"
	"strings"
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
var appLogEntry binding.String

// global label for error widget
var appLogLabel *widget.Label

// AppWindow represents application window
func AppWindow(app fyne.App, w fyne.Window) {
	w.SetContent(Create(app, w))
	w.Resize(windowSize)
	w.SetMaster()
	// custom theme
	//     app.Settings().SetTheme(&grayTheme{})
}

// helper function to unify error messages
func appLog(level, msg string, err error) {
	tstamp := time.Now().Format(time.RFC3339)
	text := fmt.Sprintf("%s %s %s %v", tstamp, level, msg, err)
	if appLogEntry == nil {
		appLogEntry = binding.NewString()
	}
	log.Println(text)

	// get previous message and keep log growing to some size
	var messages []string
	emsg, err := appLogEntry.Get()
	if err == nil {
		for _, m := range strings.Split(emsg, "\n") {
			messages = append(messages, m)
		}
	}
	messages = append(messages, text)
	// reverse messages to show last message first
	rarr := reverse(messages)
	info := strings.Join(rarr, "\n")
	if len(rarr) > 100 {
		info = strings.Join(rarr[0:9], "\n")
	}
	appLogEntry.Set(info)
}

// helper function to reverse array of strings
func reverse(arr []string) []string {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

// helper function to setup app error label
func setupAppError() {
	appLogEntry = binding.NewString()
	appLogEntry.Set("ECM error window")
	appLogLabel = widget.NewLabelWithData(appLogEntry)
	appLogLabel.Wrapping = fyne.TextWrapBreak
}

// logTabItem continer
func logTabItem(app fyne.App, w fyne.Window) *container.TabItem {
	content := container.NewVBox(
		appLogLabel,
	)
	logContainer := container.NewScroll(content)
	return &container.TabItem{Text: "Log", Icon: theme.InfoIcon(), Content: logContainer}
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
		logTabItem(app, window),
		logoutTabItem(app, window),
	}}
	return appTabs
}

// helper function to refresh all app widgets/canvases
func appRefresh(app fyne.App, window fyne.Window) {
	// refresh all widgets
	if appTabs != nil && appTabs.Items != nil {
		uiRecords.Refresh()
		for _, tab := range appTabs.Items {
			tab.Content.Refresh()
		}
	}
}

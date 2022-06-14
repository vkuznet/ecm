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

	pref := app.Preferences()
	if pref.String("AppTheme") == "light" {
		app.Settings().SetTheme(&grayTheme{})
	}
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
	year := time.Now().Year()
	sep := fmt.Sprintf("\n%d", year)
	emsg, err := appLogEntry.Get()
	if err == nil {
		for _, m := range strings.Split(emsg, sep) {
			if !strings.HasPrefix(m, fmt.Sprintf("%d", year)) {
				messages = append(messages, fmt.Sprintf("%d%s", year, m))
			} else {
				messages = append(messages, m)
			}
		}
	}
	messages = append(messages, text)
	// reverse messages to show last message first
	rarr := reverse(messages)
	info := strings.Join(rarr, "\n")
	nLines := 100 // number of log line to keep
	if len(rarr) > nLines {
		info = strings.Join(rarr[0:nLines], "\n")
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
	appLogEntry.Set("ECM log window")
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

// keep appRecords global as we'll need to update them
var appRecords *vaultRecords

// Create will stitch together all ui components
func Create(app fyne.App, window fyne.Window) *container.AppTabs {
	appRecords = newUIVaultRecords(app, window)
	appTabs = &container.AppTabs{Items: []*container.TabItem{
		appRecords.tabItem(),
		newUIRecord(app, window).tabItem(),
		newUIPassword(app, window).tabItem(),
		newUISync(app, window, appRecords).tabItem(),
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

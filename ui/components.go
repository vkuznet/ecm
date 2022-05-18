package main

import (
	"fyne.io/fyne/v2"
	canvas "fyne.io/fyne/v2/canvas"
	container "fyne.io/fyne/v2/container"
	layout "fyne.io/fyne/v2/layout"
	widget "fyne.io/fyne/v2/widget"
)

// ErrWindow represents generic error window content
func ErrWindow(app fyne.App, w fyne.Window) {
	w.SetContent(widget.NewLabel("Error window"))
}

// AppWindow represents application window
func AppWindow(app fyne.App, w fyne.Window) {
	w.SetContent(Create(app, w))
	w.Resize(windowSize)
	w.SetMaster()
}

// LoginWindow represents login window
func LoginWindow(app fyne.App, w fyne.Window) {
	password := widget.NewPasswordEntry()
	password.OnSubmitted = func(p string) {
		//         log.Println("password", p)
		var err error
		if err != nil {
			ErrWindow(app, w)
		} else {
			AppWindow(app, w)
		}
	}
	password.Resize(inputSize)
	label := ""
	formItem := widget.NewFormItem(label, password)

	form := &widget.Form{
		Items: []*widget.FormItem{formItem},
		OnSubmit: func() {
			var err error
			if err != nil {
				ErrWindow(app, w)
			} else {
				AppWindow(app, w)
			}
		},
	}
	text := canvas.NewText("Encrypted Content", fontColor)
	text.Alignment = fyne.TextAlignCenter
	spacer := &layout.Spacer{}

	// set final container
	content := container.NewVBox(
		spacer,
		text,
		form,
		spacer,
	)

	// set window settings
	w.SetContent(content)
	w.Resize(windowSize)
	w.Canvas().Focus(password)
	w.SetMaster()
}

// Create will stitch together all ui components
func Create(app fyne.App, window fyne.Window) *container.AppTabs {
	return &container.AppTabs{Items: []*container.TabItem{
		newVaultRecords(app, window).tabItem(),
		newRecord(app, window).tabItem(),
		newPassword(app, window).tabItem(),
		newSettings(app, window).tabItem(),
	}}
}

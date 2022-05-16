package main

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	canvas "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	layout "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	widget "fyne.io/fyne/v2/widget"
)

func ErrWindow(app fyne.App, w fyne.Window) {
	w.SetContent(widget.NewLabel("Error window"))
}

func AppWindow(app fyne.App, w fyne.Window) {
	w.SetContent(Create(app, w))
	w.Resize(windowSize)
	w.SetMaster()
}

func LoginWindow(app fyne.App, w fyne.Window) {
	password := widget.NewPasswordEntry()
	password.OnSubmitted = func(p string) {
		log.Println("password", p)
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
	text := canvas.NewText("Encrypted Content", color.White)
	text.Alignment = fyne.TextAlignCenter

	//     grid := container.New(layout.NewGridWrapLayout(windowSize), text, form)
	//     content := container.New(layout.NewCenterLayout(), grid)
	//     content := container.New(layout.NewCenterLayout(), text, form)
	spacer := &layout.Spacer{}
	content := container.NewVBox(
		spacer,
		text,
		form,
		spacer,
	)

	w.SetContent(content)
	w.Resize(windowSize)
	w.Canvas().Focus(password)
	w.SetMaster()
}

// Record represents new Record button
type Record struct {
	window fyne.Window
	app    fyne.App
}

func newRecord(a fyne.App, w fyne.Window) *Record {
	return &Record{app: a, window: w}
}
func (r *Record) buildUI() *fyne.Container {
	return container.NewVBox(&widget.Label{})
}
func (r *Record) tabItem() *container.TabItem {
	return &container.TabItem{Text: "", Icon: theme.ContentAddIcon(), Content: r.buildUI()}
}

// Password represents new Password button
type Password struct {
	window fyne.Window
	app    fyne.App
}

func newPassword(a fyne.App, w fyne.Window) *Password {
	return &Password{app: a, window: w}
}
func (r *Password) buildUI() *fyne.Container {
	return container.NewVBox(&widget.Label{})
}
func (r *Password) tabItem() *container.TabItem {
	return &container.TabItem{Text: "", Icon: theme.VisibilityIcon(), Content: r.buildUI()}
}

// Settings represents new Settings button
type Settings struct {
	preferences fyne.Preferences
	window      fyne.Window
	app         fyne.App
}

func newSettings(a fyne.App, w fyne.Window) *Settings {
	return &Settings{app: a, window: w, preferences: a.Preferences()}
}
func (r *Settings) buildUI() *fyne.Container {
	return container.NewVBox(&widget.Label{})
}
func (r *Settings) tabItem() *container.TabItem {
	return &container.TabItem{Text: "", Icon: theme.SettingsIcon(), Content: r.buildUI()}
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

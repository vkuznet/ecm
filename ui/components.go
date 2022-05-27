package main

import (
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	canvas "fyne.io/fyne/v2/canvas"
	container "fyne.io/fyne/v2/container"
	layout "fyne.io/fyne/v2/layout"
	widget "fyne.io/fyne/v2/widget"
	vt "github.com/vkuznet/ecm/vault"
	"golang.org/x/exp/errors"
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
}

func checkVault() {
	if appKind == "desktop" {
		if _, err := os.Stat(_vault.Directory); errors.Is(err, os.ErrNotExist) {
			err := _vault.Create(_vault.Directory)
			if err != nil {
				log.Println("unable to create vault directory", _vault.Directory, "error: ", err)
			}
		}
	}
}

// LoginWindow represents login window
func LoginWindow(app fyne.App, w fyne.Window) {
	// get vault records
	if _vault == nil {
		pref := app.Preferences()
		cipher := pref.String("VaultCipher")
		vdir := pref.String("VaultDirectory")
		_vault = &vt.Vault{Directory: vdir, Cipher: cipher, Start: time.Now()}
	}

	password := widget.NewPasswordEntry()
	password.OnSubmitted = func(p string) {
		_vault.Secret = p
		checkVault()
		err := _vault.Read()
		if err != nil {
			ErrWindow(app, w, err)
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
			_vault.Secret = password.Text
			checkVault()
			err := _vault.Read()
			if err != nil {
				ErrWindow(app, w, err)
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
	uiRecords := newUIVaultRecords(app, window)
	return &container.AppTabs{Items: []*container.TabItem{
		uiRecords.tabItem(),
		newUIRecord(app, window).tabItem(),
		newUIPassword(app, window).tabItem(),
		newUISync(app, window, uiRecords).tabItem(),
		newUISettings(app, window).tabItem(),
	}}
}

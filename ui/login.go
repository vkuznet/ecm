package main

import (
	"image/color"
	"log"
	"net/url"
	"os"
	"time"

	"fyne.io/fyne/v2"
	container "fyne.io/fyne/v2/container"
	layout "fyne.io/fyne/v2/layout"
	widget "fyne.io/fyne/v2/widget"
	vt "github.com/vkuznet/ecm/vault"
	"golang.org/x/exp/errors"
)

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

// helper function to create appropriate login button with custom text and icon
func loginButton(app fyne.App, w fyne.Window, entry *widget.Entry) *widget.Button {
	btn := &widget.Button{
		Text: "",
		//         Icon: theme.LoginIcon(),
		Icon: lockImage.Resource,
		OnTapped: func() {
			_vault.Secret = entry.Text
			checkVault()
			err := _vault.Read()
			if err != nil {
				ErrWindow(app, w, err)
			} else {
				AppWindow(app, w)
			}
		},
	}
	return btn
}

// helper function to create appropriate button with custom icon
func tapButton(resource *fyne.Resource, link string) *widget.Button {
	return &widget.Button{
		Text: "",
		Icon: *resource,
		OnTapped: func() {
			openURL(link)
		},
	}
}

// LoginWindow represents login window
func LoginWindow(app fyne.App, w fyne.Window) {
	// custom theme
	//     app.Settings().SetTheme(&grayTheme{})

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
	password.PlaceHolder = "Enter Master Password"
	//     password.Resize(inputSize)

	loginContainer := container.NewGridWrap(inputSize, password)

	btnColor := color.NRGBA{0x79, 0x79, 0x79, 0xff}
	btn := loginButton(app, w, password)
	btnContainer := colorButtonContainer(btn, btnColor)
	loginRowContainer := container.NewHBox(
		loginContainer, btnContainer,
	)
	// add image
	spacer := &layout.Spacer{}
	webImg := tapButton(&webImage.Resource, "https://github.com/vkuznet/ecm")
	docImg := tapButton(&docImage.Resource, "https://github.com/vkuznet/ecm")
	gitImg := tapButton(&gitImage.Resource, "https://github.com/vkuznet/ecm")
	contentContainer := container.NewVBox(
		spacer,
		loginRowContainer,
		//         container.NewHBox(spacer, webImage, gitImage, docImage, spacer),
		container.NewHBox(spacer, webImg, gitImg, docImg, spacer),
		spacer,
	)
	content := container.NewCenter(contentContainer)

	//     hyperlink := &widget.Hyperlink{Text: version, URL: releaseURL, TextStyle: fyne.TextStyle{Bold: true}

	// set window settings
	w.SetContent(content)
	w.Resize(windowSize)
	w.Canvas().Focus(password)
	w.SetMaster()
}

// helper function to open given URL link
func openURL(link string) error {
	u, err := url.Parse(link)
	if err != nil {
		fyne.LogError("Failed to parse url link ", err)
		return err
	}
	err = fyne.CurrentApp().OpenURL(u)
	if err != nil {
		fyne.LogError("Failed to open url", err)
		return err
	}
	return nil
}

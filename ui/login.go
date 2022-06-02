package main

import (
	"log"
	"net/url"
	"os"
	"time"

	"fyne.io/fyne/v2"
	container "fyne.io/fyne/v2/container"
	layout "fyne.io/fyne/v2/layout"
	theme "fyne.io/fyne/v2/theme"
	widget "fyne.io/fyne/v2/widget"
	vt "github.com/vkuznet/ecm/vault"
	"golang.org/x/exp/errors"
)

// main container represents main view of the app
var passwordEntry *widget.Entry

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
				errorMessage("unable to read vault records", err)
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

	passwordEntry = widget.NewPasswordEntry()
	passwordEntry.OnSubmitted = func(p string) {
		_vault.Secret = p
		checkVault()
		err := _vault.Read()
		if err != nil {
			errorMessage("unable to read vault records", err)
		} else {
			AppWindow(app, w)
		}
	}
	passwordEntry.PlaceHolder = "Enter Master Password"

	loginContainer := container.NewGridWrap(inputSize, passwordEntry)

	btn := loginButton(app, w, passwordEntry)
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
		container.NewHBox(spacer, webImg, gitImg, docImg, spacer),
		spacer,
	)
	content := container.NewCenter(contentContainer)

	//     hyperlink := &widget.Hyperlink{Text: version, URL: releaseURL, TextStyle: fyne.TextStyle{Bold: true}

	// set window settings
	w.SetContent(content)
	w.Resize(windowSize)
	w.Canvas().Focus(passwordEntry)
	w.SetMaster()
}

// logout continer
func logoutTabItem(app fyne.App, w fyne.Window) *container.TabItem {
	logout := logoutContainer(app, w)
	return &container.TabItem{Text: "Logout", Icon: theme.LogoutIcon(), Content: logout}
}
func logoutContainer(app fyne.App, w fyne.Window) *fyne.Container {
	btn := &widget.Button{
		Text: "Logout",
		Icon: theme.LogoutIcon(),
		OnTapped: func() {
			_vault.Secret = ""
			_vault.Records = nil
			passwordEntry.Text = ""
			appTabs = nil
			LoginWindow(app, w)
		},
	}
	label := widget.NewLabel("Confirm logout to reset vault access")
	content := container.NewCenter(
		container.NewVBox(
			label,
			colorButtonContainer(btn, redColor),
		),
	)
	return content
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

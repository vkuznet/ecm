package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	canvas "fyne.io/fyne/v2/canvas"
	container "fyne.io/fyne/v2/container"
	binding "fyne.io/fyne/v2/data/binding"
	layout "fyne.io/fyne/v2/layout"
	theme "fyne.io/fyne/v2/theme"
	widget "fyne.io/fyne/v2/widget"
	vt "github.com/vkuznet/ecm/vault"
	"golang.org/x/exp/errors"
)

// foreground time in seconds since epoch
var foregroundTime int64

// var autoLogoutQuit chan bool

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
		Icon: resourceLockSvg,
		OnTapped: func() {
			_vault.Secret = entry.Text
			checkVault()
			err := _vault.Read()
			if err != nil {
				appLog("ERROR", "unable to read vault records", err)
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
			appLog("ERROR", "unable to read vault records", err)
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
	gitImage := canvas.NewImageFromResource(resourceGithubSvg)
	docImage := canvas.NewImageFromResource(resourceDocSvg)
	webImage := canvas.NewImageFromResource(resourceWebSvg)

	webImg := tapButton(&webImage.Resource, "https://vkuznet.github.io/ecm/")
	docImg := tapButton(&docImage.Resource, "https://github.com/vkuznet/ecm")
	gitImg := tapButton(&gitImage.Resource, "https://github.com/vkuznet/ecm")
	contentContainer := container.NewVBox(
		spacer,
		loginRowContainer,
		container.NewHBox(spacer, webImg, gitImg, docImg, spacer),
		spacer,
	)
	content := container.NewCenter(contentContainer)

	// set autologout functions
	app.Lifecycle().SetOnEnteredForeground(func() {
		// when our app goes to foreground, i.e. it is a primary visible window
		foregroundTime = 0
	})
	app.Lifecycle().SetOnExitedForeground(func() {
		// when our app goes to background
		foregroundTime = time.Now().Unix()
	})

	// read sync config and dump it to the log
	if appKind != "desktop" {
		sconf := syncPath(app)
		appLog("INFO", sconf, nil)
		msg := fmt.Sprintf("Vault at %s has %d records", _vault.Directory, len(_vault.Records))
		appLog("INFO", msg, nil)
		err := logSyncConfig(app)
		if err != nil {
			appLog("ERROR", "unable to read sync config", err)
		}
	}

	// set window settings
	w.SetContent(content)
	w.Resize(windowSize)
	w.Canvas().Focus(passwordEntry)
	w.SetMaster()
}

// autologout threshold to use
var autoThreshold binding.String

// helper function to perform logout, should be run as goroutine
func autoLogout(app fyne.App, w fyne.Window, ctx context.Context) {
	if autoThreshold == nil {
		autoThreshold = binding.NewString()
	}
	autoThreshold.Set("60") // default autologout value
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// check if our app is sleeping and auto-logout if necessary
			now := time.Now().Unix()
			strThr, err := autoThreshold.Get()
			if err != nil {
				log.Println("unable to get autoThreshold valut", err)
				continue
			}
			thr, err := strconv.Atoi(strThr)
			if err == nil && foregroundTime > 0 && now-foregroundTime > int64(thr) {
				log.Println("autologin reset")
				_vault.Secret = ""
				_vault.Records = nil
				passwordEntry.Text = ""
				appTabs = nil
				foregroundTime = 0
				LoginWindow(app, w)
			}
			time.Sleep(time.Duration(5) * time.Second)
			//             log.Println("times", now, foregroundTime, now-foregroundTime)
		}
	}
}

// logout continer
func logoutTabItem(app fyne.App, w fyne.Window) *container.TabItem {
	logout := logoutContainer(app, w)
	return &container.TabItem{Text: "Logout", Icon: theme.LogoutIcon(), Content: logout}
}

// helper function to build logout container
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

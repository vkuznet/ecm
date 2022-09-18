package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	container "fyne.io/fyne/v2/container"
	theme "fyne.io/fyne/v2/theme"
	widget "fyne.io/fyne/v2/widget"
	crypt "github.com/vkuznet/ecm/crypt"
)

// helper function to fetch attribute value from preference and assign default
func getPrefValue(pref fyne.Preferences, key, value string) string {
	prefValue := pref.String(key)
	if prefValue != "" {
		value = prefValue
	}
	return value
}

// helper function to create bold label widget
func newBoldLabel(text string) *widget.Label {
	return &widget.Label{Text: text, TextStyle: fyne.TextStyle{Bold: true}}
}

// Settings represents new Settings button
type Settings struct {
	preferences     fyne.Preferences
	window          fyne.Window
	app             fyne.App
	theme           *widget.Select
	vaultCipher     *widget.Select
	vaultDirectory  *widget.Entry
	vaultAutologout *widget.Entry
	fontSize        *widget.Select
}

func newUISettings(a fyne.App, w fyne.Window) *Settings {
	return &Settings{app: a, window: w, preferences: a.Preferences()}
}
func (r *Settings) getPreferences() {
}

func (r *Settings) onAutologoutChanged(v string) {
	autoThreshold.Set(v)
	r.preferences.SetString("Autologout", v)
}
func (r *Settings) onSyncURLChanged(v string) {
	r.preferences.SetString("SyncURL", v)
}
func (r *Settings) onFontSizeChanged(v string) {
	// TODO: find out how to change font size based on provided value, e.g. tiny
	// for example it may require to get theme
	// r.app.Settings().Theme()
	// and adjust it accoringly
	r.preferences.SetString("FontSize", v)
}
func (r *Settings) onPasswordLengthChanged(v string) {
	r.preferences.SetString("PasswordLength", v)
}
func (r *Settings) onThemeChanged(v string) {
	// update global theme name
	appTheme = v
	// update app setttings
	if v == "dark" {
		r.app.Settings().SetTheme(theme.DarkTheme())
	} else if v == "ligth" {
		//         r.app.Settings().SetTheme(theme.LightTheme())
		r.app.Settings().SetTheme(&grayTheme{})
	} else {
		//         r.app.Settings().SetTheme(theme.LightTheme())
		r.app.Settings().SetTheme(&grayTheme{})
	}
	r.app.Preferences().SetString("AppTheme", appTheme)
	appRefresh(r.app, r.window)
}
func (r *Settings) onVaultCipherChanged(v string) {
	_vault.Cipher = v
	r.app.Preferences().SetString("VaultCipher", v)
}
func (r *Settings) onVaultDirectoryChanged(v string) {
	if _, err := os.Stat(_vault.Directory); !os.IsNotExist(err) {
		_vault.Directory = v
		_vault.Records = nil
		err := _vault.Read()
		if err != nil {
			appLog("ERROR", "fail to read vault record", err)
		} else {
			msg := fmt.Sprintf("Read vault %s, found %d records", v, len(_vault.Records))
			appLog("INFO", msg, nil)
			// refresh ui records
			if appRecords != nil {
				appRecords.Refresh()
			}
		}
	}
	r.app.Preferences().SetString("VaultDirectory", v)
}

func (r *Settings) buildUI() *container.Scroll {

	pref := r.app.Preferences()
	//     fontSize := pref.String("FontSize")
	vaultCipher := pref.String("VaultCipher")

	// set initial values of internal data
	vaultDirectory := pref.String("VaultDirectory")
	r.vaultDirectory = &widget.Entry{Text: vaultDirectory, OnSubmitted: r.onVaultDirectoryChanged}

	// set autologout settings
	r.vaultAutologout = widget.NewEntryWithData(autoThreshold)
	r.vaultAutologout.OnSubmitted = r.onAutologoutChanged

	r.vaultCipher = widget.NewSelect(crypt.SupportedCiphers, r.onVaultCipherChanged)
	r.vaultCipher.SetSelected(vaultCipher)

	fontSizes := []string{"Tiny", "Small", "Large", "Normal", "Huge"}
	r.fontSize = widget.NewSelect(fontSizes, r.onFontSizeChanged)
	r.fontSize.SetSelected(fontSize)

	themeNames := []string{"dark", "light"}
	r.theme = widget.NewSelect(themeNames, r.onThemeChanged)
	r.theme.SetSelected(appTheme)

	// TODO: add selection of font sizes

	// setup preferences
	r.getPreferences()

	//     uiContainer := appearance.NewSettings().LoadAppearanceScreen(r.window)
	uiContainer := container.NewVBox(
		newBoldLabel("Theme"),
		r.theme,
		newBoldLabel("Font size"),
		r.fontSize,
	)

	vaultContainer := container.NewVBox(
		newBoldLabel("Vault autologout"),
		r.vaultAutologout,
		newBoldLabel("Vault cipher"),
		r.vaultCipher,
		newBoldLabel("Vault directory"),
		r.vaultDirectory,
	)

	return container.NewScroll(container.NewVBox(
		newBoldLabel("Vault Settings"),
		vaultContainer,
		newBoldLabel("User interface"),
		uiContainer,
	))
}
func (r *Settings) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Settings", Icon: theme.SettingsIcon(), Content: r.buildUI()}
}

// helper function to print preferences area/files
func printPrefs(a fyne.App) {
	dir := a.Storage().RootURI().Path()
	fmt.Println(dir)
	if files, err := ioutil.ReadDir(dir); err == nil {
		for _, finfo := range files {
			if finfo.IsDir() {
				continue
			}
			fmt.Println(finfo.Name())
			fname := filepath.Join(dir, finfo.Name())
			file, err := os.Open(fname)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			data, err := ioutil.ReadAll(file)
			if err == nil {
				fmt.Println(string(data))
			}
		}
	}
}

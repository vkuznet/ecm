package main

import (
	"fyne.io/fyne/v2"
	container "fyne.io/fyne/v2/container"
	theme "fyne.io/fyne/v2/theme"
	widget "fyne.io/fyne/v2/widget"
	crypt "github.com/vkuznet/ecm/crypt"
)

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
	//     fontSize       *widget.Select
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
		r.app.Settings().SetTheme(theme.LightTheme())
	} else {
		r.app.Settings().SetTheme(theme.LightTheme())
	}
	appRefresh(r.app, r.window)
}
func (r *Settings) onVaultCipherChanged(v string) {
	_vault.Cipher = v
	r.app.Preferences().SetString("VaultCipher", v)
}
func (r *Settings) onVaultDirectoryChanged(v string) {
	_vault.Directory = v
	_vault.Records = nil
	err := _vault.Read()
	if err != nil {
		errorMessage("fail to read vault record", err)
	}
	r.app.Preferences().SetString("VaultDirectory", v)
}

func (r *Settings) buildUI() *container.Scroll {

	pref := r.app.Preferences()
	//     fontSize := pref.String("FontSize")
	vaultCipher := pref.String("VaultCipher")
	vaultDirectory := pref.String("VaultDirectory")

	// set initial values of internal data
	r.vaultDirectory = &widget.Entry{Text: vaultDirectory, OnSubmitted: r.onVaultDirectoryChanged}

	// set autologout settings
	r.vaultAutologout = widget.NewEntryWithData(autoThreshold)
	r.vaultAutologout.OnSubmitted = r.onAutologoutChanged

	r.vaultCipher = widget.NewSelect(crypt.SupportedCiphers, r.onVaultCipherChanged)
	r.vaultCipher.SetSelected(vaultCipher)

	//     fontSizes := []string{"Tiny", "Small", "Large", "Normal", "Huge"}
	//     r.fontSize = widget.NewSelect(fontSizes, r.onFontSizeChanged)
	//     r.fontSize.SetSelected(fontSize)

	themeNames := []string{"dark", "light"}
	r.theme = widget.NewSelect(themeNames, r.onThemeChanged)
	r.theme.SetSelected(appTheme)
	//     r.theme.SetSelected("light")

	// TODO: add selection of font sizes

	// setup preferences
	r.getPreferences()

	//     uiContainer := appearance.NewSettings().LoadAppearanceScreen(r.window)
	uiContainer := container.NewVBox(
		newBoldLabel("Theme"),
		r.theme,
		//         newBoldLabel("Font size"),
		//         r.fontSize,
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

package main

import (
	"log"

	"fyne.io/fyne/v2"
	container "fyne.io/fyne/v2/container"
	theme "fyne.io/fyne/v2/theme"
	widget "fyne.io/fyne/v2/widget"
	crypt "github.com/vkuznet/ecm/crypt"
)

func newBoldLabel(text string) *widget.Label {
	return &widget.Label{Text: text, TextStyle: fyne.TextStyle{Bold: true}}
}
func onOrOff(on bool) string {
	if on {
		return "On"
	}
	return "Off"
}

// Settings represents new Settings button
type Settings struct {
	preferences    fyne.Preferences
	window         fyne.Window
	app            fyne.App
	theme          *widget.Select
	verifyRadio    *widget.RadioGroup
	vaultCipher    *widget.Select
	vaultDirectory *widget.Entry
	syncURL        *widget.Entry
	fontSize       *widget.Select
}

func newUISettings(a fyne.App, w fyne.Window) *Settings {
	return &Settings{app: a, window: w, preferences: a.Preferences()}
}
func (r *Settings) getPreferences() {
	verify := r.preferences.Bool("Verify")
	r.verifyRadio.Selected = onOrOff(verify)
}

func (r *Settings) onVerifyChanged(selected string) {
	enabled := selected == "On"
	r.app.Preferences().SetBool("Verify", enabled)
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
	if v == "dark" {
		r.app.Settings().SetTheme(theme.DarkTheme())
	} else if v == "ligth" {
		r.app.Settings().SetTheme(theme.LightTheme())
	} else {
		r.app.Settings().SetTheme(theme.LightTheme())
	}
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
		// TODO: redirect to ErrWindow
		log.Println("fail to read vault record", err)
	}
	r.app.Preferences().SetString("VaultDirectory", v)
}

func (r *Settings) buildUI() *container.Scroll {

	pref := r.app.Preferences()
	fontSize := pref.String("FontSize")
	vaultCipher := pref.String("VaultCipher")
	vaultDirectory := pref.String("VaultDirectory")

	// set initial values of internal data
	onOffOptions := []string{"On", "Off"}
	r.verifyRadio = &widget.RadioGroup{
		Options: onOffOptions, Horizontal: true, Required: true, OnChanged: r.onVerifyChanged}
	r.syncURL = &widget.Entry{PlaceHolder: "hostname", OnChanged: r.onSyncURLChanged}
	r.vaultDirectory = &widget.Entry{Text: vaultDirectory, OnSubmitted: r.onVaultDirectoryChanged}

	r.vaultCipher = widget.NewSelect(crypt.SupportedCiphers, r.onVaultCipherChanged)
	r.vaultCipher.SetSelected(vaultCipher)

	fontSizes := []string{"Tiny", "Small", "Large", "Normal", "Huge"}
	r.fontSize = widget.NewSelect(fontSizes, r.onFontSizeChanged)
	r.fontSize.SetSelected(fontSize)

	themeNames := []string{"dark", "light"}
	r.theme = widget.NewSelect(themeNames, r.onThemeChanged)
	r.theme.SetSelected("dark")
	//     r.theme.SetSelected("light")

	// TODO: add selection of font sizes

	// setup preferences
	r.getPreferences()

	//     uiContainer := appearance.NewSettings().LoadAppearanceScreen(r.window)
	uiContainer := container.NewVBox(
		container.NewGridWithColumns(2,
			newBoldLabel("Theme"), r.theme,
			newBoldLabel("Font size"), r.fontSize,
		),
	)

	vaultContainer := container.NewVBox(
		container.NewGridWithColumns(2,
			newBoldLabel("Verify before accepting"), r.verifyRadio,
			newBoldLabel("Vault cipher"), r.vaultCipher,
		),
		&widget.Accordion{Items: []*widget.AccordionItem{
			{Title: "Advanced", Detail: container.NewGridWithColumns(2,
				newBoldLabel("Vault directory"), r.vaultDirectory,
				newBoldLabel("Sync URL"), r.syncURL,
			)},
		}},
	)

	return container.NewScroll(container.NewVBox(
		&widget.Card{Title: "Vault Settings", Content: vaultContainer},
		&widget.Card{Title: "User Interface", Content: uiContainer},
	))
}
func (r *Settings) tabItem() *container.TabItem {
	return &container.TabItem{Text: "", Icon: theme.SettingsIcon(), Content: r.buildUI()}
}

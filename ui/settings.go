package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	container "fyne.io/fyne/v2/container"
	theme "fyne.io/fyne/v2/theme"
	widget "fyne.io/fyne/v2/widget"
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
	passwordLength *widget.Entry
	storagePath    *widget.Entry
	syncURL        *widget.Entry
}

func newSettings(a fyne.App, w fyne.Window) *Settings {
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
func (r *Settings) onStoragePathChanged(v string) {
	r.preferences.SetString("StoragePath", v)
}
func (r *Settings) onPasswordLengthChanged(v string) {
	r.preferences.SetString("PasswordLength", v)
}
func (r *Settings) onThemeChanged(v string) {
	if v == "dark" {
		r.app.Settings().SetTheme(theme.DarkTheme())
		fontColor = color.White
	} else if v == "ligth" {
		r.app.Settings().SetTheme(theme.LightTheme())
		fontColor = color.Black
	} else {
		r.app.Settings().SetTheme(theme.LightTheme())
		fontColor = color.Black
	}
	//     canvas.Refresh(r.app)
}

func (r *Settings) buildUI() *container.Scroll {

	// set initial values of internal data
	onOffOptions := []string{"On", "Off"}
	r.verifyRadio = &widget.RadioGroup{
		Options: onOffOptions, Horizontal: true, Required: true, OnChanged: r.onVerifyChanged}
	r.passwordLength = &widget.Entry{PlaceHolder: "length", OnChanged: r.onPasswordLengthChanged}
	r.syncURL = &widget.Entry{PlaceHolder: "hostname", OnChanged: r.onSyncURLChanged}
	r.storagePath = &widget.Entry{PlaceHolder: "~/Dropbox", OnChanged: r.onStoragePathChanged}
	themeNames := []string{"dark", "light"}
	r.theme = widget.NewSelect(themeNames, r.onThemeChanged)
	r.theme.SetSelected("dark")

	// TODO: add selection of font sizes

	// setup preferences
	r.getPreferences()

	//     uiContainer := appearance.NewSettings().LoadAppearanceScreen(r.window)
	uiContainer := container.NewVBox(
		container.NewGridWithColumns(2,
			newBoldLabel("Theme"), r.theme,
		),
	)

	vaultContainer := container.NewVBox(
		container.NewGridWithColumns(2,
			newBoldLabel("Verify before accepting"), r.verifyRadio,
			newBoldLabel("Password Length"), r.passwordLength,
			newBoldLabel("StoragePath"), r.storagePath,
		),
		&widget.Accordion{Items: []*widget.AccordionItem{
			{Title: "Advanced", Detail: container.NewGridWithColumns(2,
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

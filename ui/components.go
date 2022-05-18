package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	canvas "fyne.io/fyne/v2/canvas"
	container "fyne.io/fyne/v2/container"
	layout "fyne.io/fyne/v2/layout"
	theme "fyne.io/fyne/v2/theme"
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

// Record represents new Record button
type Record struct {
	window     fyne.Window
	app        fyne.App
	Name       *widget.Entry
	URL        *widget.Entry
	Login      *widget.Entry
	Password   *widget.Entry
	Tags       *widget.Entry
	JSON       *widget.Entry
	Note       *widget.Entry
	CardNumber *widget.Entry
	Code       *widget.Entry
	Date       *widget.Entry
	Phone      *widget.Entry
	Encryption *widget.Entry
	FromHost   *widget.Entry
	ToHost     *widget.Entry
	File       *widget.Entry
}

func newRecord(a fyne.App, w fyne.Window) *Record {
	return &Record{app: a, window: w}
}

func (r *Record) LoginForm() {
}
func (r *Record) JSONForm() {
}
func (r *Record) NoteForm() {
}
func (r *Record) CardForm() {
}
func (r *Record) VaultForm() {
}
func (r *Record) SyncForm() {
}
func (r *Record) FileUploadForm() {
}
func (r *Record) SelectCipher(v string) {
}

func (r *Record) buildUI() *container.Scroll {
	name := &widget.Entry{PlaceHolder: "name"}
	url := &widget.Entry{PlaceHolder: "e.g. http://cnn.com"}
	login := &widget.Entry{PlaceHolder: "login name"}
	password := widget.NewPasswordEntry()
	tags := &widget.Entry{PlaceHolder: "tag1,tag2"}
	json := &widget.Entry{PlaceHolder: "{\"foo\":1}", MultiLine: true}
	note := &widget.Entry{PlaceHolder: "some note", MultiLine: true}
	card := &widget.Entry{PlaceHolder: "123"}
	code := &widget.Entry{PlaceHolder: "123"}
	date := &widget.Entry{PlaceHolder: "mm/dd/YYYY"}
	phone := &widget.Entry{PlaceHolder: "+1-888-888-8888"}
	ciphers := []string{"AES", "NaCl"}
	encryption := widget.NewSelect(ciphers, r.SelectCipher)
	fromHost := &widget.Entry{PlaceHolder: "www.www.com"}
	toHost := &widget.Entry{PlaceHolder: "a.b.com"}
	// TODO: change file to be drop area or select from file system
	file := &widget.Entry{PlaceHolder: "file"}

	loginForm := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Name", name),
			widget.NewFormItem("URL", url),
			widget.NewFormItem("Login", login),
			widget.NewFormItem("Password", password),
			widget.NewFormItem("Tags", tags),
		},
		OnSubmit: r.LoginForm,
	}
	loginContainer := container.NewVBox(loginForm)

	jsonForm := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Name", name),
			widget.NewFormItem("JSON", json),
		},
		OnSubmit: r.JSONForm,
	}
	jsonContainer := container.NewVBox(jsonForm)

	noteForm := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Name", name),
			widget.NewFormItem("Note", note),
		},
		OnSubmit: r.NoteForm,
	}
	noteContainer := container.NewVBox(noteForm)

	cardForm := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Name", name),
			widget.NewFormItem("Card Number", card),
			widget.NewFormItem("Verification code", code),
			widget.NewFormItem("Date", date),
			widget.NewFormItem("Phone", phone),
			widget.NewFormItem("Tags", tags),
		},
		OnSubmit: r.CardForm,
	}
	cardContainer := container.NewVBox(cardForm)

	syncForm := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("From host", fromHost),
			widget.NewFormItem("To host", toHost),
		},
		OnSubmit: r.SyncForm,
	}
	syncContainer := container.NewVBox(syncForm)

	vaultForm := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Name", name),
			widget.NewFormItem("Encryption", encryption),
		},
		OnSubmit: r.VaultForm,
	}
	vaultContainer := container.NewVBox(vaultForm)

	fileForm := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("File", file),
		},
		OnSubmit: r.FileUploadForm,
	}
	fileContainer := container.NewVBox(fileForm)

	return container.NewScroll(container.NewVBox(
		&widget.Accordion{Items: []*widget.AccordionItem{
			{Title: "Login Record", Detail: loginContainer},
			{Title: "JSON Record", Detail: jsonContainer},
			{Title: "Note Record", Detail: noteContainer},
			{Title: "Card Record", Detail: cardContainer},
			{Title: "File upload", Detail: fileContainer},
			{Title: "Vault", Detail: vaultContainer},
			{Title: "Sync", Detail: syncContainer},
		}},
	))
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
func (r *Password) GeneratePassword() {
}
func (r *Password) CharactersChange(v string) {
}
func (r *Password) buildUI() *fyne.Container {
	// widgets
	spacer := &layout.Spacer{}
	genPassword := &widget.Entry{}
	//     icon := widget.NewIcon(theme.ContentCopyIcon())
	icon := &widget.Button{
		Text: "Copy",
		Icon: theme.ContentCopyIcon(),
		OnTapped: func() {
			text := genPassword.Text
			r.window.Clipboard().SetContent(text)
		},
	}
	size := &widget.Entry{Text: "16"}
	names := []string{"letters", "letters+digits"}
	characters := widget.NewSelect(names, r.CharactersChange)

	// form widget
	passForm := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Characters", characters),
			widget.NewFormItem("Size", size),
		},
		SubmitText: "Generate password",
		OnSubmit: func() {
			genPassword.Text = "some new password"
			genPassword.Refresh()
		},
	}

	// password container
	passContainer := container.NewVBox(
		container.NewGridWithColumns(3,
			newBoldLabel("Generated password"), genPassword, icon,
		),
	)

	// final container
	return container.NewVBox(
		passForm,
		spacer,
		passContainer,
	)
}
func (r *Password) tabItem() *container.TabItem {
	return &container.TabItem{Text: "", Icon: theme.VisibilityIcon(), Content: r.buildUI()}
}

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

// Create will stitch together all ui components
func Create(app fyne.App, window fyne.Window) *container.AppTabs {
	return &container.AppTabs{Items: []*container.TabItem{
		newVaultRecords(app, window).tabItem(),
		newRecord(app, window).tabItem(),
		newPassword(app, window).tabItem(),
		newSettings(app, window).tabItem(),
	}}
}

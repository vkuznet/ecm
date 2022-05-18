package main

import (
	"fyne.io/fyne/v2"
	container "fyne.io/fyne/v2/container"
	theme "fyne.io/fyne/v2/theme"
	widget "fyne.io/fyne/v2/widget"
)

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

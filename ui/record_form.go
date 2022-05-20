package main

import (
	"log"

	"fyne.io/fyne/v2"
	container "fyne.io/fyne/v2/container"
	binding "fyne.io/fyne/v2/data/binding"
	theme "fyne.io/fyne/v2/theme"
	widget "fyne.io/fyne/v2/widget"
	vt "github.com/vkuznet/ecm/vault"
)

// Record represents new Record button
type Record struct {
	window fyne.Window
	app    fyne.App

	// bindings
	Name       binding.String
	URL        binding.String
	Login      binding.String
	Password   binding.String
	Tags       binding.String
	JSON       binding.String
	Note       binding.String
	CardNumber binding.String
	Code       binding.String
	Date       binding.String
	Phone      binding.String
	Encryption binding.String
	FromHost   binding.String
	ToHost     binding.String
	File       binding.String
}

func newRecord(a fyne.App, w fyne.Window) *Record {
	return &Record{
		app:        a,
		window:     w,
		Name:       binding.NewString(),
		URL:        binding.NewString(),
		Login:      binding.NewString(),
		Password:   binding.NewString(),
		Tags:       binding.NewString(),
		JSON:       binding.NewString(),
		Note:       binding.NewString(),
		CardNumber: binding.NewString(),
		Code:       binding.NewString(),
		Date:       binding.NewString(),
		Phone:      binding.NewString(),
		Encryption: binding.NewString(),
		FromHost:   binding.NewString(),
		ToHost:     binding.NewString(),
		File:       binding.NewString(),
	}
}

func (r *Record) ResetBindings() {
	r.Name = binding.NewString()
	r.URL = binding.NewString()
	r.Login = binding.NewString()
	r.Password = binding.NewString()
	r.Tags = binding.NewString()
	r.JSON = binding.NewString()
	r.Note = binding.NewString()
	r.CardNumber = binding.NewString()
	r.Code = binding.NewString()
	r.Date = binding.NewString()
	r.Phone = binding.NewString()
	r.Encryption = binding.NewString()
	r.FromHost = binding.NewString()
	r.ToHost = binding.NewString()
	r.File = binding.NewString()
}
func (r *Record) LoginForm() {
	rec := vt.NewVaultRecord("login")
	if val, err := r.Name.Get(); err == nil {
		rec.Map["Name"] = val
	}
	if val, err := r.Login.Get(); err == nil {
		rec.Map["Login"] = val
	}
	if val, err := r.Password.Get(); err == nil {
		rec.Map["Password"] = val
	}
	if val, err := r.Tags.Get(); err == nil {
		rec.Map["Tags"] = val
	}
	if val, err := r.Note.Get(); err == nil {
		rec.Map["Note"] = val
	}
	if val, err := r.URL.Get(); err == nil {
		rec.Map["URL"] = val
	}
	log.Println("New vault record", rec.String())
	err := _vault.Update(*rec)
	if err != nil {
		log.Println("ERROR", "unable to write vault record")
	}
	//     r.ResetBindings()
	r.window.SetContent(Create(r.app, r.window))
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

	name := widget.NewEntryWithData(r.Name)
	name.PlaceHolder = "record name"

	url := widget.NewEntryWithData(r.URL)
	url.PlaceHolder = "e.g. http://abc.com"

	login := widget.NewEntryWithData(r.Login)
	login.PlaceHolder = "login name"

	password := widget.NewEntryWithData(r.Password)
	password.Password = true

	tags := widget.NewEntryWithData(r.Tags)
	tags.PlaceHolder = "tag1,tag2,..."

	note := widget.NewEntryWithData(r.Note)
	note.PlaceHolder = "some text"
	note.MultiLine = true

	json := &widget.Entry{PlaceHolder: "{\"foo\":1}", MultiLine: true}
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

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

// LoginRecord represents login UI record
type LoginRecord struct {
	Name     binding.String
	URL      binding.String
	Login    binding.String
	Password binding.String
	Note     binding.String
	Tags     binding.String
}

// JSONRecord represents json UI record
type JSONRecord struct {
	Name binding.String
	JSON binding.String
}

// CardRecord represents card UI record
type CardRecord struct {
	Name       binding.String
	CardNumber binding.String
	Code       binding.String
	Date       binding.String
	Phone      binding.String
	Tags       binding.String
}

// NoteRecord represents note UI record
type NoteRecord struct {
	Name binding.String
	Note binding.String
	Tags binding.String
}

// SyncRecord represents sycn UI record
type SyncRecord struct {
	Name    binding.String
	FromURI binding.String
	ToURI   binding.String
}

// UploadRecord represents sycn UI record
type UploadRecord struct {
	Name binding.String
	File binding.String
}

// Record represents new UI Record
type Record struct {
	window fyne.Window
	app    fyne.App

	// binding records
	LoginRecord  *LoginRecord
	JSONRecord   *JSONRecord
	NoteRecord   *NoteRecord
	CardRecord   *CardRecord
	SyncRecord   *SyncRecord
	UploadRecord *UploadRecord
}

func newLoginRecord() *LoginRecord {
	return &LoginRecord{
		Name:     binding.NewString(),
		URL:      binding.NewString(),
		Login:    binding.NewString(),
		Password: binding.NewString(),
		Note:     binding.NewString(),
		Tags:     binding.NewString(),
	}
}
func newJSONRecord() *JSONRecord {
	return &JSONRecord{
		Name: binding.NewString(),
		JSON: binding.NewString(),
	}
}
func newNoteRecord() *NoteRecord {
	return &NoteRecord{
		Name: binding.NewString(),
		Note: binding.NewString(),
		Tags: binding.NewString(),
	}
}
func newCardRecord() *CardRecord {
	return &CardRecord{
		Name:       binding.NewString(),
		CardNumber: binding.NewString(),
		Code:       binding.NewString(),
		Date:       binding.NewString(),
		Phone:      binding.NewString(),
		Tags:       binding.NewString(),
	}
}
func newSyncRecord() *SyncRecord {
	return &SyncRecord{
		FromURI: binding.NewString(),
		ToURI:   binding.NewString(),
	}
}
func newUploadRecord() *UploadRecord {
	return &UploadRecord{
		File: binding.NewString(),
	}
}

func newRecord(a fyne.App, w fyne.Window) *Record {
	return &Record{
		app:          a,
		window:       w,
		LoginRecord:  newLoginRecord(),
		JSONRecord:   newJSONRecord(),
		NoteRecord:   newNoteRecord(),
		CardRecord:   newCardRecord(),
		SyncRecord:   newSyncRecord(),
		UploadRecord: newUploadRecord(),
	}
}

// helper function used by all form records to update provided vault record
func (r *Record) updateVaultRecord(rec *vt.VaultRecord) {
	if _vault.Verbose > 0 {
		log.Println("New vault record", rec.String())
	}
	err := _vault.Update(*rec)
	if err != nil {
		log.Println("ERROR", "unable to write vault record")
	}
	r.window.SetContent(Create(r.app, r.window))
}
func (r *Record) LoginForm() {
	rec := vt.NewVaultRecord("login")
	if val, err := r.LoginRecord.Name.Get(); err == nil {
		rec.Map["Name"] = val
	}
	if val, err := r.LoginRecord.Login.Get(); err == nil {
		rec.Map["Login"] = val
	}
	if val, err := r.LoginRecord.Password.Get(); err == nil {
		rec.Map["Password"] = val
	}
	if val, err := r.LoginRecord.Tags.Get(); err == nil {
		rec.Map["Tags"] = val
	}
	if val, err := r.LoginRecord.Note.Get(); err == nil {
		rec.Map["Note"] = val
	}
	if val, err := r.LoginRecord.URL.Get(); err == nil {
		rec.Map["URL"] = val
	}
	r.updateVaultRecord(rec)
}
func (r *Record) JSONForm() {
	rec := vt.NewVaultRecord("json")
	if val, err := r.LoginRecord.Name.Get(); err == nil {
		rec.Map["Name"] = val
	}
	if val, err := r.LoginRecord.Name.Get(); err == nil {
		rec.Map["JSON"] = val
	}
	r.updateVaultRecord(rec)
}
func (r *Record) NoteForm() {
	rec := vt.NewVaultRecord("note")
	if val, err := r.NoteRecord.Name.Get(); err == nil {
		rec.Map["Name"] = val
	}
	if val, err := r.NoteRecord.Name.Get(); err == nil {
		rec.Map["Note"] = val
	}
	if val, err := r.NoteRecord.Tags.Get(); err == nil {
		rec.Map["Tags"] = val
	}
	r.updateVaultRecord(rec)
}
func (r *Record) CardForm() {
	rec := vt.NewVaultRecord("card")
	if val, err := r.CardRecord.Name.Get(); err == nil {
		rec.Map["Name"] = val
	}
	if val, err := r.CardRecord.Name.Get(); err == nil {
		rec.Map["Card"] = val
	}
	if val, err := r.CardRecord.Name.Get(); err == nil {
		rec.Map["Code"] = val
	}
	if val, err := r.CardRecord.Name.Get(); err == nil {
		rec.Map["Date"] = val
	}
	if val, err := r.CardRecord.Name.Get(); err == nil {
		rec.Map["Phone"] = val
	}
	if val, err := r.CardRecord.Tags.Get(); err == nil {
		rec.Map["Tags"] = val
	}
	r.updateVaultRecord(rec)
}
func (r *Record) SyncForm() {
	fromURI, _ := r.SyncRecord.FromURI.Get()
	toURI, _ := r.SyncRecord.ToURI.Get()
	// TODO: implement sync action
	log.Println("sync is not implemented", fromURI, toURI)
}
func (r *Record) UploadForm() {
	fname, _ := r.UploadRecord.File.Get()
	// TODO: implement file upload
	log.Println("file upload not implemented", fname)
}

func (r *Record) buildUI() *container.Scroll {

	// login form container
	loginEntryName := widget.NewEntryWithData(r.LoginRecord.Name)
	loginEntryName.PlaceHolder = "record name"

	loginEntryUrl := widget.NewEntryWithData(r.LoginRecord.URL)
	loginEntryUrl.PlaceHolder = "e.g. http://abc.com"

	loginEntryLogin := widget.NewEntryWithData(r.LoginRecord.Login)
	loginEntryLogin.PlaceHolder = "login name"

	loginEntryPassword := widget.NewEntryWithData(r.LoginRecord.Password)
	loginEntryPassword.Password = true

	loginEntryTags := widget.NewEntryWithData(r.LoginRecord.Tags)
	loginEntryTags.PlaceHolder = "tag1,tag2,..."

	loginEntryNote := widget.NewEntryWithData(r.LoginRecord.Note)
	loginEntryNote.PlaceHolder = "some text"
	loginEntryNote.MultiLine = true

	loginForm := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Name", loginEntryName),
			widget.NewFormItem("URL", loginEntryUrl),
			widget.NewFormItem("Login", loginEntryLogin),
			widget.NewFormItem("Password", loginEntryPassword),
			widget.NewFormItem("Tags", loginEntryTags),
		},
		OnSubmit: r.LoginForm,
	}
	loginContainer := container.NewVBox(loginForm)

	// json form container
	jsonEntryName := widget.NewEntryWithData(r.JSONRecord.Name)
	jsonEntryName.PlaceHolder = "record name"
	jsonEntryJson := &widget.Entry{PlaceHolder: "{\"foo\":1}", MultiLine: true}

	jsonForm := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Name", jsonEntryName),
			widget.NewFormItem("JSON", jsonEntryJson),
		},
		OnSubmit: r.JSONForm,
	}
	jsonContainer := container.NewVBox(jsonForm)

	// note form container
	noteEntryName := widget.NewEntryWithData(r.NoteRecord.Name)
	noteEntryName.PlaceHolder = "record name"
	noteEntryNote := widget.NewEntryWithData(r.NoteRecord.Note)
	noteEntryNote.PlaceHolder = "record name"
	noteEntryTags := widget.NewEntryWithData(r.NoteRecord.Tags)
	noteEntryTags.PlaceHolder = "tag1,tag2,..."
	noteForm := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Name", noteEntryName),
			widget.NewFormItem("Note", noteEntryNote),
			widget.NewFormItem("Tags", noteEntryTags),
		},
		OnSubmit: r.NoteForm,
	}
	noteContainer := container.NewVBox(noteForm)

	// card form container
	cardEntryName := widget.NewEntryWithData(r.CardRecord.Name)
	cardEntryName.PlaceHolder = "record name"
	cardEntryCard := widget.NewEntryWithData(r.CardRecord.CardNumber)
	cardEntryCard.PlaceHolder = "1234-4567-xxxx"
	cardEntryCode := widget.NewEntryWithData(r.CardRecord.Code)
	cardEntryCode.PlaceHolder = "123"
	cardEntryDate := widget.NewEntryWithData(r.CardRecord.Date)
	cardEntryDate.PlaceHolder = "mm/dd/YYYY"
	cardEntryPhone := widget.NewEntryWithData(r.CardRecord.Phone)
	cardEntryPhone.PlaceHolder = "+1-888-888-8888"
	cardEntryTags := widget.NewEntryWithData(r.CardRecord.Tags)
	cardEntryTags.PlaceHolder = "tag1,tag2,..."
	cardForm := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Name", cardEntryName),
			widget.NewFormItem("Card Number", cardEntryCard),
			widget.NewFormItem("CVC Number", cardEntryCode),
			widget.NewFormItem("Date", cardEntryDate),
			widget.NewFormItem("Phone", cardEntryPhone),
			widget.NewFormItem("Tags", cardEntryTags),
		},
		OnSubmit: r.CardForm,
	}
	cardContainer := container.NewVBox(cardForm)

	// sync form container
	fromURI := widget.NewEntryWithData(r.SyncRecord.FromURI)
	fromURI.PlaceHolder = "from URI"
	toURI := widget.NewEntryWithData(r.SyncRecord.ToURI)
	toURI.PlaceHolder = "to URI"
	syncForm := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("From URI", fromURI),
			widget.NewFormItem("To URI", toURI),
		},
		OnSubmit: r.SyncForm,
	}
	syncContainer := container.NewVBox(syncForm)

	// TODO: change file to be drop area or select from file system
	// upload form container
	fileEntryFile := widget.NewEntryWithData(r.UploadRecord.File)
	fileEntryFile.PlaceHolder = "file name"
	fileForm := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("File", fileEntryFile),
		},
		OnSubmit: r.UploadForm,
	}
	fileContainer := container.NewVBox(fileForm)

	return container.NewScroll(container.NewVBox(
		&widget.Accordion{Items: []*widget.AccordionItem{
			{Title: "Login Record", Detail: loginContainer},
			{Title: "JSON Record", Detail: jsonContainer},
			{Title: "Note Record", Detail: noteContainer},
			{Title: "Card Record", Detail: cardContainer},
			{Title: "File upload", Detail: fileContainer},
			{Title: "Sync", Detail: syncContainer},
		}},
	))
}
func (r *Record) tabItem() *container.TabItem {
	return &container.TabItem{Text: "", Icon: theme.ContentAddIcon(), Content: r.buildUI()}
}

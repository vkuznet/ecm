package main

import (
	"fmt"

	fyne "fyne.io/fyne/v2"
	container "fyne.io/fyne/v2/container"
	theme "fyne.io/fyne/v2/theme"
	widget "fyne.io/fyne/v2/widget"
	utils "github.com/vkuznet/ecm/utils"
	vt "github.com/vkuznet/ecm/vault"
)

type vaultRecords struct {
	window  fyne.Window
	app     fyne.App
	records map[string]Entry
}

func newUIVaultRecords(a fyne.App, w fyne.Window) *vaultRecords {
	return &vaultRecords{app: a, window: w}
}

// helper function to get vault record name or ID
func recordName(rec vt.VaultRecord) string {
	name := rec.ID
	if v, ok := rec.Map["Name"]; ok {
		name = fmt.Sprintf("%s/ID:%s", v, rec.ID)
	}
	return name
}

// helper function to build recordsList
func (a *vaultRecords) buildRecordsList(records []vt.VaultRecord) *widget.Accordion {
	entries := widget.NewAccordion()
	for _, rec := range records {
		entries.Append(widget.NewAccordionItem(recordName(rec), a.recordContainer(rec)))
	}
	return entries
}

// global variable to keep accordion records
// we will refresh it dyring sync process
var uiRecords *widget.Accordion

// Refresh refresh records in UI
func (a *vaultRecords) Refresh() {
	uiRecords.Items = nil
	for _, rec := range _vault.Records {
		uiRecords.Append(widget.NewAccordionItem(recordName(rec), a.recordContainer(rec)))
	}
	uiRecords.Refresh()
}

func (a *vaultRecords) buildUI() *container.Scroll {

	// build initial set of accordion records
	accRecords := a.buildRecordsList(_vault.Records)
	uiRecords = accRecords

	// setup search entry
	search := widget.NewEntry()
	search.OnSubmitted = func(v string) {
		// reset items of accordion
		// see https://yourbasic.org/golang/clear-slice/
		accRecords.Items = nil
		for _, rec := range _vault.Find(v) {
			accRecords.Append(widget.NewAccordionItem(recordName(rec), a.recordContainer(rec)))
		}
		accRecords.Refresh()
	}
	search.PlaceHolder = "search keyword"
	label := ""
	formItem := widget.NewFormItem(label, search)

	form := &widget.Form{
		Items: []*widget.FormItem{formItem},
	}

	// return final container with search and accordion records
	return container.NewScroll(container.NewVBox(
		form,
		container.NewVBox(accRecords),
	))
}

// helper function to create record form item representation
// func (a *vaultRecords) formItem(key, val string) *widget.FormItem {
func (a *vaultRecords) formItem(vrec vt.VaultRecord, key, val string) *widget.FormItem {
	if key == "Password" || key == "password" {
		rec := widget.NewPasswordEntry()
		//         rec := widget.NewEntryWithData(val)
		//         rec.Password = true
		rec.Text = val
		rec.Refresh()
		recContainer := container.NewVBox(
			container.NewGridWithColumns(2,
				rec, a.copyIcon(rec),
			),
		)
		return widget.NewFormItem(key, recContainer)
	}
	rec := &widget.Entry{
		Text: val,
		OnChanged: func(v string) {
			vrec.Map[key] = v
		},
	}
	//     rec := widget.NewEntryWithData(val)
	return widget.NewFormItem(key, rec)
}

// helper function to create record representation
func (a *vaultRecords) recordContainer(record vt.VaultRecord) *fyne.Container {
	// create entry object
	var items []*widget.FormItem
	for _, k := range vt.OrderedKeys {
		if v, ok := record.Map[k]; ok {
			items = append(items, a.formItem(record, k, v))
		}
	}
	// other keys can follow ordered ones
	for k, v := range record.Map {
		if !utils.InList(k, vt.OrderedKeys) {
			items = append(items, a.formItem(record, k, v))
		}
	}
	form := &widget.Form{
		Items:      items,
		SubmitText: "Update",
		OnSubmit: func() {
			record.WriteRecord(_vault.Directory, _vault.Secret, _vault.Cipher, _vault.Verbose)
		},
	}
	recContainer := container.NewVBox(form)
	return recContainer
}

// helper function to create appropriate copy icon
func (a *vaultRecords) copyIcon(entry *widget.Entry) *widget.Button {
	icon := &widget.Button{
		Text: "Copy",
		Icon: theme.ContentCopyIcon(),
		OnTapped: func() {
			text := entry.Text
			a.window.Clipboard().SetContent(text)
		},
	}
	return icon
}

func (a *vaultRecords) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Records", Icon: theme.ListIcon(), Content: a.buildUI()}
}

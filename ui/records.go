package main

import (
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	utils "github.com/vkuznet/ecm/utils"
	vt "github.com/vkuznet/ecm/vault"
)

type vaultRecords struct {
	window  fyne.Window
	app     fyne.App
	records map[string]Entry
}

func newVaultRecords(a fyne.App, w fyne.Window) *vaultRecords {
	return &vaultRecords{app: a, window: w}
}

func (a *vaultRecords) buildUI() *container.Scroll {
	search := widget.NewEntry()
	search.PlaceHolder = "search keyword"
	label := ""
	formItem := widget.NewFormItem(label, search)

	form := &widget.Form{
		Items: []*widget.FormItem{formItem},
	}

	// list of entries
	entries := widget.NewAccordion()
	for _, rec := range _vault.Records {
		entries.Append(widget.NewAccordionItem(rec.ID, a.recordContainer(rec.Map)))
	}

	return container.NewScroll(container.NewVBox(
		form,
		entries,
	))
}

// helper function to create record form item representation
func (a *vaultRecords) formItem(key, val string) *widget.FormItem {
	if key == "Password" || key == "password" {
		rec := widget.NewPasswordEntry()
		rec.Text = val
		rec.Refresh()
		recContainer := container.NewVBox(
			container.NewGridWithColumns(2,
				rec, a.copyIcon(rec),
			),
		)
		return widget.NewFormItem(key, recContainer)
	}
	rec := &widget.Entry{Text: val, OnChanged: func(v string) {}}
	return widget.NewFormItem(key, rec)
}

// helper function to create record representation
func (a *vaultRecords) recordContainer(record vt.Record) *fyne.Container {
	// create entry object
	orderedKeys := []string{
		"Name", "Login", "Password", "URL", "Tags", "Note",
	}
	var items []*widget.FormItem
	for _, k := range orderedKeys {
		if v, ok := record[k]; ok {
			items = append(items, a.formItem(k, v))
		}
	}
	// other keys can follow ordered ones
	for k, v := range record {
		if !utils.InList(k, orderedKeys) {
			items = append(items, a.formItem(k, v))
		}
	}
	form := &widget.Form{
		Items:      items,
		SubmitText: "Update",
		OnSubmit:   func() { /* TODO: update entry record */ },
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

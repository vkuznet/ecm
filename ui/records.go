package main

import (
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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
	rec := &widget.Entry{Text: val, OnChanged: func(v string) {}}
	recContainer := container.NewVBox(
		container.NewGridWithColumns(2,
			rec, a.copyIcon(rec),
		),
	)
	return widget.NewFormItem(key, recContainer)
}

// helper function to create record representation
func (a *vaultRecords) recordContainer(record vt.Record) *fyne.Container {
	// create entry object
	var items []*widget.FormItem
	for k, v := range record {
		items = append(items, a.formItem(k, v))
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

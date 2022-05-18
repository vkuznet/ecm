package main

import (
	"fmt"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type vaultRecords struct {
	window  fyne.Window
	app     fyne.App
	records map[string]Entry
}

func newVaultRecords(a fyne.App, w fyne.Window) *vaultRecords {
	return &vaultRecords{app: a, window: w}
}

func makeRecords() map[string]Entry {
	records := make(map[string]Entry)
	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("key-%d", i)
		records[key] = Entry{Text: key}
	}
	return records
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
	items := makeRecords()
	for key, entry := range items {
		entries.Append(widget.NewAccordionItem(key, a.recordContainer(entry)))
		//         entries.Append(widget.NewAccordionItem(key, &widget.Entry{Text: entry.Text}))
	}

	return container.NewScroll(container.NewVBox(
		form,
		entries,
	))
}

// helper function to create record representation
func (a *vaultRecords) recordContainer(entry Entry) *fyne.Container {
	// create entry object
	name := &widget.Entry{Text: entry.Text, OnChanged: func(v string) {}}
	login := &widget.Entry{Text: "some login", OnChanged: func(v string) {}}
	recContainer := container.NewVBox(
		container.NewGridWithColumns(3,
			newBoldLabel("Name"), name, a.copyIcon(name),
			newBoldLabel("Login"), login, a.copyIcon(login),
		),
	)
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

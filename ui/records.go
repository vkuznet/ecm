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
		entries.Append(widget.NewAccordionItem(key, &widget.Entry{Text: entry.Text}))
	}
	//     icon := container.NewHBox(spacer, a.icon, spacer),

	return container.NewScroll(container.NewVBox(
		form,
		entries,
	))
}

func (a *vaultRecords) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Records", Icon: theme.ListIcon(), Content: a.buildUI()}
}

package main

import (
	"fmt"
	"image/color"
	"strings"

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
		hash := strings.Split(rec.ID, "-")
		name = fmt.Sprintf("%s : %s", v, hash[0])
	}
	return name
}

// helper function to provide row container
func (a *vaultRecords) rowContainer(rec vt.VaultRecord) *fyne.Container {
	var objects []fyne.CanvasObject
	for _, k := range vt.OrderedKeys {
		if v, ok := rec.Map[k]; ok {
			objects = append(objects, a.singleRow(k, v))
		}
	}
	for k, v := range rec.Map {
		if !utils.InList(k, vt.OrderedKeys) {
			objects = append(objects, a.singleRow(k, v))
		}
	}
	// TODO: need a button with OnTapped acition
	//     btn := copyButton(a.window, "Update", "", theme.MenuIcon())
	btnColor := color.NRGBA{0x79, 0x79, 0x79, 0xff}
	btn := copyButton(a.window, "Update", "", theme.MenuIcon())
	btnContainer := colorButtonContainer(btn, btnColor)

	objects = append(objects, btnContainer)
	return container.NewVBox(objects...)
}

// helper function to create single row container
func (a *vaultRecords) singleRow(key, val string) *fyne.Container {
	label := widget.NewLabel(key)
	entry := widget.NewEntry()
	entry.Text = val
	if key == "Password" || key == "password" {
		entry = widget.NewPasswordEntry()
		entry.Text = val
		entry.Refresh()
	}
	btn := container.NewVBox(
		copyButton(a.window, "", val, theme.ContentCopyIcon()),
	)

	// specify explicitly size of our elements in a container
	labelSize := fyne.NewSize(100, 40)
	entrySize := fyne.NewSize(200, 40)
	buttonSize := fyne.NewSize(40, 40)

	//     label.Resize(labelSize)
	labelContainer := container.NewGridWrap(labelSize, label)
	entryContainer := container.NewGridWrap(entrySize, entry)
	buttonContainer := container.NewGridWrap(buttonSize, btn)
	return container.NewHBox(
		labelContainer, entryContainer, buttonContainer,
	)
}

// func NewGridWithColumns(cols int, objects ...fyne.CanvasObject) *fyne.Container {
//     return fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(cols), objects...)
// }

// helper function to build recordsList
func (a *vaultRecords) buildRecordsList(records []vt.VaultRecord) *widget.Accordion {
	entries := widget.NewAccordion()
	for _, rec := range records {
		//         entries.Append(widget.NewAccordionItem(recordName(rec), a.recordContainer(rec)))
		entries.Append(widget.NewAccordionItem(recordName(rec), a.rowContainer(rec)))
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
		//         uiRecords.Append(widget.NewAccordionItem(recordName(rec), a.recordContainer(rec)))
		uiRecords.Append(widget.NewAccordionItem(recordName(rec), a.rowContainer(rec)))
	}
	uiRecords.Refresh()
}

// helper function to create appropriate copy button with custom text and icon
func (a *vaultRecords) searchButton(entry *widget.Entry) *widget.Button {
	return &widget.Button{
		Text: "",
		Icon: theme.SearchIcon(),
		OnTapped: func() {
			key := entry.Text
			// reset items of accordion
			// see https://yourbasic.org/golang/clear-slice/
			uiRecords.Items = nil
			for _, rec := range _vault.Find(key) {
				//                 uiRecords.Append(widget.NewAccordionItem(recordName(rec), a.recordContainer(rec)))
				uiRecords.Append(widget.NewAccordionItem(recordName(rec), a.rowContainer(rec)))
			}
			uiRecords.Refresh()
		},
	}
}

// helper function to build list form UI
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
			//             accRecords.Append(widget.NewAccordionItem(recordName(rec), a.recordContainer(rec)))
			accRecords.Append(widget.NewAccordionItem(recordName(rec), a.rowContainer(rec)))
		}
		accRecords.Refresh()
	}
	search.PlaceHolder = "search keyword"
	searchContainer := container.NewGridWrap(rowSize, search)

	// TODO: assign OnTapped action to perform search across records see OnSubmitted function
	btnColor := color.NRGBA{0x79, 0x79, 0x79, 0xff}
	btn := a.searchButton(search)
	btnContainer := colorButtonContainer(btn, btnColor)
	searchRowContainer := container.NewHBox(
		searchContainer, btnContainer,
	)

	// return final container with search and accordion records
	return container.NewScroll(container.NewVBox(
		container.NewVBox(searchRowContainer),
		container.NewVBox(accRecords),
	))
}

// helper function to create record form item representation
// func (a *vaultRecords) formItem(key, val string) *widget.FormItem {
func (a *vaultRecords) formItem(vrec vt.VaultRecord, key, val string) *widget.FormItem {
	if key == "Password" || key == "password" {
		rec := widget.NewPasswordEntry()
		rec.Text = val
		rec.Refresh()
		recContainer := container.NewVBox(
			container.NewGridWithColumns(2,
				rec, copyButton(a.window, key, val, theme.ContentCopyIcon()),
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
func (a *vaultRecords) copyIcon(entry *widget.Entry, txt string) *widget.Button {
	icon := &widget.Button{
		Text: txt,
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

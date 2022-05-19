package main

import (
	"fmt"
	"image/color"
	"log"
	"net/url"

	fyne "fyne.io/fyne/v2"
	canvas "fyne.io/fyne/v2/canvas"
	container "fyne.io/fyne/v2/container"
	layout "fyne.io/fyne/v2/layout"
	widget "fyne.io/fyne/v2/widget"
)

func ErrorWindow(w fyne.Window) {
	w.SetContent(widget.NewLabel("Error window"))
}
func ListWindow(w fyne.Window) {
	// add new list widget
	var data = []string{"record1", "record2", "record3"}
	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Records")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i])
		})

	// set final grid object
	grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(700, 500)),
		list)
	w.SetContent(grid)
}
func TableWindow(w fyne.Window) {
	// add new table
	var data = [][]string{[]string{"top left", "top right"},
		[]string{"bottom left", "bottom right"}}
	list := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col])
		})

	// set final grid object
	grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(700, 500)),
		list)
	w.SetContent(grid)
}

func MainWindow(app fyne.App, w fyne.Window) {
	//     records := make(map[string]Entry)
	//     for i := 0; i < 20; i++ {
	//         key := fmt.Sprintf("key-%d", i)
	//         records[key] = Entry{Text: key}
	//     }

	password := widget.NewPasswordEntry()
	password.OnSubmitted = func(p string) {
		log.Println("password", p)
		var err error
		if err != nil {
			ErrorWindow(w)
		} else {
			//                         ListWindow(w)
			//             Accordion(w, records)
			w.SetContent(Create(app, w))
			w.Resize(fyne.NewSize(700, 400))
			w.SetMaster()
		}
	}
	label := ""
	formItem := widget.NewFormItem(label, password)

	form := &widget.Form{
		Items: []*widget.FormItem{formItem},
		OnSubmit: func() {
			var err error
			if err != nil {
				ErrorWindow(w)
			} else {
				//                 ListWindow(w)
				//                 Accordion(w, records)
				w.SetContent(Create(app, w))
				w.Resize(fyne.NewSize(700, 400))
				w.SetMaster()
			}
		},
	}
	text := canvas.NewText("Encrypted Content", color.White)
	text.Alignment = fyne.TextAlignCenter

	// set final grid object
	//     gridSize := fyne.NewSize(500, 100)
	gridSize := fyne.NewSize(300, 300)
	grid := container.New(
		layout.NewGridWrapLayout(gridSize),
		text, form)
	w.SetContent(grid)

	w.Canvas().Focus(password)
}

type Entry struct {
	Text string
}

func Accordion(w fyne.Window, items map[string]Entry) {
	search := widget.NewEntry()
	label := ""
	formItem := widget.NewFormItem(label, search)

	form := &widget.Form{
		Items: []*widget.FormItem{formItem},
	}

	// list of entries
	entries := widget.NewAccordion()
	for key, entry := range items {
		entries.Append(widget.NewAccordionItem(key, &widget.Entry{Text: entry.Text}))
	}
	//     w.SetContent(entries)

	// set final grid object
	//     gridSize := fyne.NewSize(500, 100)
	gridSize := fyne.NewSize(300, 300)
	grid := container.New(
		layout.NewGridWrapLayout(gridSize),
		form, entries)
	w.SetContent(grid)
}

func makeAccordionTab(_ fyne.Window) fyne.CanvasObject {
	link, err := url.Parse("https://fyne.io/")
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}
	ac := widget.NewAccordion(
		widget.NewAccordionItem("A", widget.NewHyperlink("One", link)),
		widget.NewAccordionItem("B", widget.NewLabel("Two")),
		&widget.AccordionItem{
			Title:  "C",
			Detail: widget.NewLabel("Three"),
		},
	)
	ac.Append(widget.NewAccordionItem("D", &widget.Entry{Text: "Four"}))
	return ac
}
func makeRecords() map[string]Entry {
	records := make(map[string]Entry)
	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("key-%d", i)
		records[key] = Entry{Text: key}
	}
	return records
}

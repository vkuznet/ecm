package main

import (
	"image/color"

	fyne "fyne.io/fyne/v2"
	canvas "fyne.io/fyne/v2/canvas"
	container "fyne.io/fyne/v2/container"
	layout "fyne.io/fyne/v2/layout"
	widget "fyne.io/fyne/v2/widget"
)

// helper function to make custom entry button
func entryButton(bname string) *fyne.Container {
	btn := widget.NewButton(bname, nil)
	btn_color := canvas.NewRectangle(
		color.NRGBA{0xd6, 0xd6, 0xd6, 0xff})
	return container.New(
		layout.NewMaxLayout(),
		btn_color,
		btn,
	)
}

// helper function to create appropriate copy button with custom text and icon
func copyButton(w fyne.Window, entry *widget.Entry, txt string, icon fyne.Resource) *widget.Button {
	return &widget.Button{
		Text: txt,
		Icon: icon,
		OnTapped: func() {
			text := entry.Text
			w.Clipboard().SetContent(text)
		},
	}
}

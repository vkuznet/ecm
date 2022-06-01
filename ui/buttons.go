package main

import (
	"image/color"

	fyne "fyne.io/fyne/v2"
	canvas "fyne.io/fyne/v2/canvas"
	container "fyne.io/fyne/v2/container"
	layout "fyne.io/fyne/v2/layout"
	widget "fyne.io/fyne/v2/widget"
)

// for more help see
// https://blogvali.com/category/fyne-golang-gui/

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
func copyButton(w fyne.Window, label, txt string, icon fyne.Resource) *widget.Button {
	return &widget.Button{
		Text: label,
		Icon: icon,
		OnTapped: func() {
			text := txt
			if txt != "" {
				w.Clipboard().SetContent(text)
			}
		},
	}
}

// helper function to wrap button into color container
func colorButtonContainer(btn *widget.Button, nrgba color.NRGBA) *fyne.Container {
	btn_color := canvas.NewRectangle(nrgba)
	return container.New(
		layout.NewMaxLayout(),
		btn_color,
		btn,
	)
}

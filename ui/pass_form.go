package main

import (
	"fyne.io/fyne/v2"
	container "fyne.io/fyne/v2/container"
	layout "fyne.io/fyne/v2/layout"
	theme "fyne.io/fyne/v2/theme"
	widget "fyne.io/fyne/v2/widget"
)

// Password represents new Password button
type Password struct {
	window fyne.Window
	app    fyne.App
}

func newPassword(a fyne.App, w fyne.Window) *Password {
	return &Password{app: a, window: w}
}
func (r *Password) GeneratePassword() {
}
func (r *Password) CharactersChange(v string) {
}
func (r *Password) buildUI() *fyne.Container {
	// widgets
	spacer := &layout.Spacer{}
	genPassword := &widget.Entry{}
	icon := &widget.Button{
		Text: "Copy",
		Icon: theme.ContentCopyIcon(),
		OnTapped: func() {
			text := genPassword.Text
			r.window.Clipboard().SetContent(text)
		},
	}
	size := &widget.Entry{Text: "16"}
	names := []string{"letters", "letters+digits"}
	characters := widget.NewSelect(names, r.CharactersChange)

	// form widget
	passForm := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Characters", characters),
			widget.NewFormItem("Size", size),
		},
		SubmitText: "Generate password",
		OnSubmit: func() {
			genPassword.Text = "some new password"
			genPassword.Refresh()
		},
	}

	// password container
	passContainer := container.NewVBox(
		container.NewGridWithColumns(3,
			newBoldLabel("Generated password"), genPassword, icon,
		),
	)

	// final container
	return container.NewVBox(
		passForm,
		spacer,
		passContainer,
	)
}
func (r *Password) tabItem() *container.TabItem {
	return &container.TabItem{Text: "", Icon: theme.VisibilityIcon(), Content: r.buildUI()}
}

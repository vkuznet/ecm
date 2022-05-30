package main

import (
	"log"
	"strings"

	"fyne.io/fyne/v2"
	container "fyne.io/fyne/v2/container"
	binding "fyne.io/fyne/v2/data/binding"
	//     layout "fyne.io/fyne/v2/layout"
	//     theme "fyne.io/fyne/v2/theme"
	widget "fyne.io/fyne/v2/widget"
	crypt "github.com/vkuznet/ecm/crypt"
)

// Password represents new Password button
type Password struct {
	window fyne.Window
	app    fyne.App
}

func newUIPassword(a fyne.App, w fyne.Window) *Password {
	return &Password{app: a, window: w}
}
func (r *Password) GeneratePassword() {
}
func (r *Password) CharactersChange(v string) {
}
func (r *Password) buildUI() *fyne.Container {
	// widgets
	genPassword := &widget.Entry{PlaceHolder: "generated password"}
	length := binding.NewInt()
	length.Set(16)
	strLength := binding.IntToString(length)
	size := widget.NewEntryWithData(strLength)
	names := []string{"letters", "letters+digits", "letters+digits+symbols"}
	characters := widget.NewSelect(names, r.CharactersChange)
	characters.SetSelected("letters+digits+symbols")

	// form widget
	passForm := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("New password", genPassword),
			widget.NewFormItem("Characters", characters),
			widget.NewFormItem("Size", size),
		},
		SubmitText: "Generate password",
		OnSubmit: func() {
			var hasNumbers, hasSymbols bool
			idx := characters.SelectedIndex()
			if idx > -1 && strings.Contains(names[idx], "digits") {
				hasNumbers = true
			} else if idx > -1 && strings.Contains(names[idx], "symbols") {
				hasSymbols = true
			}
			val, err := length.Get()
			if err != nil {
				log.Println("ERROR:", "TODO SOMETHING")
			}
			genPassword.Text = crypt.CreatePassword(val, hasNumbers, hasSymbols)
			genPassword.Refresh()
			r.window.Clipboard().SetContent(genPassword.Text)
		},
	}

	// final container
	return container.NewVBox(
		passForm,
	)
}
func (r *Password) tabItem() *container.TabItem {
	//     return &container.TabItem{Text: "Password", Icon: theme.VisibilityIcon(), Content: r.buildUI()}
	return &container.TabItem{Text: "Password", Icon: passImage.Resource, Content: r.buildUI()}
}

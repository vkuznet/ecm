package main

import (
	"image/color"

	fyne "fyne.io/fyne/v2"
	theme "fyne.io/fyne/v2/theme"
)

// see more:
// https://github.com/andydotxyz/beebui/blob/master/theme.go

type grayTheme struct{}

var _ fyne.Theme = (*grayTheme)(nil)

func (m grayTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		if variant == theme.VariantLight {
			//                 return color.White
			// if we coose light schema, we'll use gray color
			return &color.NRGBA{0xd6, 0xd6, 0xd6, 0xff}
		}
		return color.Black
	}
	return theme.DefaultTheme().Color(name, variant)
}
func (m grayTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	if name == theme.IconNameHome {
		//         fyne.NewStaticResource("myHome", homeBytes)
	}
	return theme.DefaultTheme().Icon(name)
}
func (m grayTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTextMonospaceFont()
	//     return theme.DefaultTheme().Font(style)
}

func (m grayTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

package gtkMod

import "github.com/gotk3/gotk3/gtk"

//IconButtonNew makes a button with icon from file
func IconButtonNew(filename string) (*gtk.Button, *gtk.Image) {
	button, _ := gtk.ButtonNew()
	icon, _ := gtk.ImageNewFromFile(filename)
	button.SetImage(icon)

	return button, icon
}

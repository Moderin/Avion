package gtkMod

import "github.com/gotk3/gotk3/gtk"

// EventImageNew makes an eventBox with image
func EventImageNew(path string) (*gtk.EventBox, *gtk.Image) {
	eventBox, _ := gtk.EventBoxNew()
	image, _ := gtk.ImageNewFromFile(path)
	eventBox.Add(image)

	return eventBox, image
}

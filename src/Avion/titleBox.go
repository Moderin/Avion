package main

import (
	"utilites"

	"github.com/gotk3/gotk3/gtk"
)

// TitleBox is in headerbar
type TitleBox struct {
	nameEntry      *gtk.Entry
	statusEntry    *gtk.Entry
	nameRevealer   *gtk.Revealer
	statusRevealer *gtk.Revealer
	picture        *utilites.PictureCircle
}

// UpdateData checks if user changed name/status and saves it.
func (t *TitleBox) UpdateData() {
	t.nameRevealer.SetRevealChild(false)
	t.statusRevealer.SetRevealChild(false)

	name, _ := t.nameEntry.GetText()
	status, _ := t.statusEntry.GetText()
	var save = false
	// if user changed name
	if *config["name"] != name {
		config["name"] = &name
		tox.SelfSetName(name)
		save = true
	}

	if config[status] == nil || *config["status"] != status {
		config["status"] = &status
		tox.SelfSetStatusMessage(status)
		save = true
	}

	if save {
		configFileSave()
	}
}

func titleBoxNew() *TitleBox {
	t := new(TitleBox)
	titleBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	titleEBox, _ := gtk.EventBoxNew()
	t.picture = utilites.AvatarDrNew(config["avatar"], 16, 0)

	t.nameRevealer, _ = gtk.RevealerNew()
	t.nameRevealer.SetTransitionType(gtk.REVEALER_TRANSITION_TYPE_SLIDE_LEFT)

	t.statusRevealer, _ = gtk.RevealerNew()
	t.statusRevealer.SetTransitionType(gtk.REVEALER_TRANSITION_TYPE_SLIDE_RIGHT)

	titleEBox.Connect("button_press_event", func() {
		if t.nameRevealer.GetRevealChild() {
			t.UpdateData()
		} else {
			t.nameRevealer.SetRevealChild(true)
			t.statusRevealer.SetRevealChild(true)
		}
	})

	t.nameEntry, _ = gtk.EntryNew()
	t.nameEntry.SetText(*config["name"])
	t.nameEntry.Connect("activate", t.UpdateData)
	style, _ := t.nameEntry.GetStyleContext()
	style.AddClass("label-entry")
	t.nameEntry.SetAlignment(1)

	t.statusEntry, _ = gtk.EntryNew()
	if config["status"] != nil {
		t.statusEntry.SetText(*config["status"])
	}
	t.statusEntry.Connect("activate", t.UpdateData)
	style, _ = t.statusEntry.GetStyleContext()
	style.AddClass("label-entry")
	style.AddClass("italic")

	t.nameRevealer.Add(t.nameEntry)
	t.statusRevealer.Add(t.statusEntry)
	titleEBox.Add(t.picture)

	titleBox.Add(t.nameRevealer)
	titleBox.Add(titleEBox)
	titleBox.Add(t.statusRevealer)

	topBar.SetCustomTitle(titleBox)

	return t
}

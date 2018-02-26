package main

import (
	"utilites"
	"encoding/json"

	"github.com/gotk3/gotk3/gtk"
)

const (
	Waiting = 0
	Sent    = 1
	User    = 50000
	Friend  = 3
)

type message interface {
	send(friend uint32) bool //	used for queued messages
	getBox() *gtk.Box
	SetState(state uint8)
	updateFace()
	SaveJSONData(*json.Encoder)
}

type messageBase struct {
	statusIcon *gtk.Image
	face       *utilites.PictureCircle
	box        *gtk.Box
	state      uint8
	author 	   uint32
}


func (m *messageBase) getBox() *gtk.Box {
	return m.box
}

func (m *messageBase) updateFace() {
	m.face.UpdateSource()
}

// SetState changes message status icon
func (m *messageBase) SetState(state uint8) {
	if m.state == Waiting && state == Sent {
		m.statusIcon.SetFromFile("img/message_status/sent.png")
		m.state = Sent
		return
	}

	if state == Waiting {
		m.statusIcon.SetFromFile("img/message_status/waiting.png")
		m.state = Waiting
		return
	}

}

func (m *messageBase) loadBase(author uint32) {
	m.author = author
	m.statusIcon, _ = gtk.ImageNew()
	m.statusIcon.SetMarginEnd(20)

	m.box, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	m.box.SetMarginStart(15)
	m.box.SetMarginEnd(15)
	m.box.SetMarginTop(10)
	m.box.SetMarginBottom(10)

	if author != User {
		m.face = utilites.AvatarDrNew(&contacts[author].avatarName, 16, 0)
	} else {
		m.face = utilites.AvatarDrNew(config["avatar"], 16, 0)
	}

	m.box.Add(m.statusIcon)
	m.box.Add(m.face)
	m.box.SetHExpand(true)
}

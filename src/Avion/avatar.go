package main

import (
	"encoding/hex"
	"math"
	"utilites"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

const (
	offline = 0
	online  = 1
	away    = 2
	busy    = 3
)

// Avatar is struct with profile picture and status background (this colourful circle)
type Avatar struct {
	overlay    *gtk.Overlay
	background *gtk.DrawingArea
	picture    *utilites.PictureCircle
	status     uint8
	newMessage bool
}

// SetNewMessage indicates, that contact got new message
func (a *Avatar) SetNewMessage(newMessage bool) {
	a.newMessage = newMessage
	a.background.QueueDraw()
}

// SetStatus changes background color of a avatar border
func (a *Avatar) SetStatus(status uint8) {
	a.status = status
	a.background.QueueDraw()
}

// AvatarNew makes new avatar
func AvatarNew(friend uint32) *Avatar {
	a := new(Avatar)

	publicKey, _ := tox.FriendGetPublickey(friend)
	contacts[friend].avatarName = hex.EncodeToString(publicKey) + ".png"
	a.picture = utilites.AvatarDrNew(&contacts[friend].avatarName, 32, 5)

	a.background, _ = gtk.DrawingAreaNew()
	a.background.SetSizeRequest(74, 74)
	a.background.Connect("draw", func(dr *gtk.DrawingArea, context *cairo.Context) {
		switch a.status {
		case offline:
			context.SetSourceRGB(0.84, 0.21, 0.26)
		case online:
			context.SetSourceRGB(0.56, 0.76, 0.25)
		case away:
			context.SetSourceRGB(0.94, 0.8, 0.35)
		case busy:
			context.SetSourceRGB(0.02, 0.63, 0.7)
		}

		context.SetLineWidth(5)
		if a.newMessage == true {
			context.Arc(37, 37, 34, 4, 0.9*2*math.Pi+4)
		} else {
			context.Arc(37, 37, 34, 0, 2*math.Pi)
		}

		context.Stroke()

		context.Clip()

		context.Paint()
	})

	a.overlay, _ = gtk.OverlayNew()
	a.overlay.AddOverlay(a.picture)
	a.overlay.Add(a.background)

	return a
}

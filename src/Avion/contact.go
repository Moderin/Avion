package main

import (
	"color"
	"encoding/hex"
	"fmt"
	"os"
	"unicode/utf8"
	"utilites"

	gotox "github.com/codedust/go-tox"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// Contact in list on the left
type Contact struct {
	friend                       uint32
	avatarName                   string
	Contantiner, nameContantiner *gtk.Box
	EventBox                     *gtk.EventBox
	nameLabel, statusLabel       *gtk.Label
	avatar                       *Avatar
	messagesBox                  *MessagesBox

	//	right-click menu
	contextMenu  *gtk.Menu
	removeButton *gtk.MenuItem
}

// UpdateName changes contact name
func (c *Contact) UpdateName(name string) {
	c.nameLabel.SetMarkup("<big>" + utilites.Shorten(name, 8) + "</big>")
	if utf8.RuneCountInString(name) > 8 {
		c.nameLabel.SetTooltipText(name)
	} else {
		c.nameLabel.SetTooltipText("")
	}
}

// UpdateStatus changes avatar background
func (c *Contact) UpdateStatus(status gotox.ToxUserStatus) {
	switch status {
	case gotox.TOX_USERSTATUS_NONE:
		c.avatar.SetStatus(online)
	case gotox.TOX_USERSTATUS_AWAY:
		c.avatar.SetStatus(away)
	case gotox.TOX_USERSTATUS_BUSY:
		c.avatar.SetStatus(busy)
	}
}

// UpdateStatusMsg changes status message
func (c *Contact) UpdateStatusMsg(status string) {
	c.statusLabel.SetMarkup("<small>" + utilites.Shorten(status, 15) + "</small>")
	if utf8.RuneCountInString(status) > 15 {
		c.statusLabel.SetTooltipText(status)
	} else {
		c.statusLabel.SetTooltipText("")
	}
}

// Init inits contact
func (c *Contact) Init(friend uint32, name, status string) *gtk.EventBox {
	c.friend = friend
	publicKey, _ := tox.FriendGetPublickey(friend)
	c.avatarName = hex.EncodeToString(publicKey) + ".png"
	c.avatar = AvatarNew(friend)

	c.contextMenu, _ = gtk.MenuNew()
	c.removeButton, _ = gtk.MenuItemNewWithLabel("Remove")
	c.removeButton.Connect("button_press_event", func() {
		publicKey, _ := tox.FriendGetPublickey(friend)
		avatarFilename := hex.EncodeToString(publicKey) + ".png"

		err := tox.FriendDelete(friend)
		if err != nil {
			fmt.Println(color.Red("Error removing friend: "), err)
			return
		}
		fmt.Println("rm!", friend)
		parent, _ := c.EventBox.GetParent()
		contactsBox.Remove(parent)
		if activeMsgBox == c.messagesBox {
			mainContantiner.Remove(activeMsgBox.scroll)
			activeMsgBox = nil
		}

		// remove profile picture if exists
		os.Remove(avatarFilename)
	})

	c.contextMenu.Append(c.removeButton)
	c.contextMenu.ShowAll()

	////////////////////////////////////////////////////////////////////////////////////////

	c.nameContantiner, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)

	// Add label with name
	c.nameLabel, _ = gtk.LabelNew("")
	c.UpdateName(name)
	c.nameLabel.SetHAlign(gtk.ALIGN_START)
	c.nameContantiner.Add(c.nameLabel)

	//Add label with status
	c.statusLabel, _ = gtk.LabelNew("")
	c.UpdateStatusMsg(status)
	c.statusLabel.SetHAlign(gtk.ALIGN_START)
	c.nameContantiner.Add(c.statusLabel)

	c.Contantiner, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 20)

	c.Contantiner.SetMarginStart(20)
	c.Contantiner.SetMarginEnd(20)
	c.Contantiner.SetMarginTop(10)
	c.Contantiner.SetMarginBottom(10)

	c.Contantiner.Add(c.avatar.overlay)
	c.Contantiner.Add(c.nameContantiner)

	c.EventBox, _ = gtk.EventBoxNew()
	c.EventBox.Add(c.Contantiner)
	c.EventBox.Connect("button_press_event", func(win *gtk.Window, event *gdk.Event) {
		pressEvent := gdk.EventButton{event}
		if pressEvent.Button() == 3 {
			c.contextMenu.PopupAtMouseCursor(nil, nil, 3, pressEvent.Time())
		} else {
			if settingsBox.active {
				settingsBox.close()
			}
			if activeMsgBox != nil {
				mainContantiner.Remove(activeMsgBox.scroll)
			}
			activeMsgBox = c.messagesBox
			c.avatar.SetNewMessage(false)
			mainContantiner.PackEnd(activeMsgBox.scroll, true, true, 0)
			mainContantiner.ShowAll()
		}
	})

	c.messagesBox = new(MessagesBox)
	c.messagesBox.Init(friend)

	return c.EventBox
}

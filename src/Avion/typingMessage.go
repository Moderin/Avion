package main

import (
	"color"
	"fmt"
	"gtkMod"
	"os"
	"strings"
	"utilites"

	gotox "github.com/codedust/go-tox"
	"github.com/gotk3/gotk3/gtk"
)

type typingMessage struct {
	messageBase
	entry                  *gtk.Entry
	fallbackPos            int
	emoiButton, fileButton *gtk.EventBox
	ePopover               *gtk.Popover
	isVisible              bool
}

func (t *typingMessage) SaveJSONData() {
	
}

func (t *typingMessage) restore() {
	activeMsgBox.box.Add(t.box)
	t.isVisible = true

	t.entry.GrabFocus()
	t.entry.SetPosition(t.fallbackPos)
	activeMsgBox.sizeAllocate()
}

func (t *typingMessage) remove() {
	t.fallbackPos = t.entry.GetPosition()
	activeMsgBox.box.Remove(t.box)
	t.isVisible = false
}



func (t *typingMessage) send(data string) {

	//	try to send message
	id, _ := tox.FriendSendMessage(activeMsgBox.friend, gotox.TOX_MESSAGE_TYPE_NORMAL, data)
	var msg message
	if emoticons[data] != "" {
		msg = emojiMessageNew(User, data)
		activeMsgBox.AddUserMessage(msg, false)
	} else {
		msg = textMessageNew(User, data)
		activeMsgBox.AddUserMessage(msg, false)
	}
	t.entry.SetText("")
	if id == 0 {
		activeMsgBox.queue = append(activeMsgBox.queue, activeMsgBox.message)
		activeMsgBox.message.SetState(Waiting)
	}

	// save to file
	activeMsgBox.save(msg)
}

// TypingMessageNew returns empty user message
func typingMessageNew() *typingMessage {
	msg := new(typingMessage)
	msg.statusIcon, _ = gtk.ImageNewFromFile("img/message_status/typing.png")
	msg.statusIcon.SetMarginEnd(20)

	msg.face = utilites.AvatarDrNew(config["avatar"], 16, 0)
	msg.face.SetSizeRequest(32, 32)

	msg.entry, _ = gtk.EntryNew()
	style, _ := msg.entry.GetStyleContext()
	style.AddClass("message-entry")

	msg.entry.Connect("activate", func() {
		entryText, _ := msg.entry.GetText()
		if strings.TrimSpace(entryText) != "" {
			msg.send(entryText)
		}
	})

	// -----------------------------------------
	// 						Emoticons
	// ----------------------------------------

	msg.emoiButton, _ = gtkMod.EventImageNew("img/emoji/smile.png")

	msg.ePopover, _ = gtk.PopoverNew(msg.emoiButton)
	eBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)

	for _, code := range userEmoticons {
		pCode := code
		imageName := emoticons[code]
		button, _ := gtkMod.EventImageNew(imageName)
		button.Connect("button_press_event", func() { msg.send(pCode) })
		button.SetBorderWidth(5)
		eBox.Add(button)
	}

	msg.ePopover.Add(eBox)

	msg.emoiButton.Connect("button_press_event", func() {

		if msg.ePopover.GetVisible() {
			msg.ePopover.Hide()
		} else {
			msg.ePopover.ShowAll()
		}
	})

	msg.fileButton, _ = gtkMod.EventImageNew("img/file.png")

	// ----------------------------------
	// 					file chooser
	// ---------------------------------

	msg.fileButton.Connect("button_press_event", func() {
		dialog, _ := gtk.FileChooserDialogNewWith2Buttons("Choose attachment", win, gtk.FILE_CHOOSER_ACTION_OPEN, "That's it!", gtk.RESPONSE_OK, "Cancel", gtk.RESPONSE_CANCEL)

		if dialog.Run() == int(gtk.RESPONSE_OK) {
			fmt.Println(color.Green("File choosen: "), dialog.GetFilename())

			file, err := os.Open(dialog.GetFilename())
			if err == nil {
				stat, _ := file.Stat()
				transfer := &FileTransfer{size: uint64(stat.Size()), handle: file, name: stat.Name()}
				activeMsgBox.box.Remove(msg.box)
				fileMsg := transferMessageNew(User, transfer)
				activeMsgBox.box.Add(fileMsg.box)
				fileMsg.box.ShowAll()
				activeMsgBox.box.Add(msg.box)

				//try to send a file
				done := fileMsg.send(activeMsgBox.friend)
				if !done {
					activeMsgBox.queue = append(activeMsgBox.queue, fileMsg)
				}

			} else {
				fmt.Println(color.Red("Cannot open selected file :( "), err)
			}
		}
		dialog.Destroy()

		// semd file and grab focus on entry
	})

	msg.box, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	msg.box.SetMarginStart(15)
	msg.box.SetMarginEnd(15)
	msg.box.SetMarginTop(10)
	msg.box.SetMarginBottom(10)

	// put everything together
	msg.box.Add(msg.statusIcon)
	msg.box.Add(msg.face)
	msg.box.PackStart(msg.entry, true, true, 0)
	msg.box.Add(msg.emoiButton)
	msg.box.Add(msg.fileButton)

	msg.box.SetHExpand(true)
	msg.box.ShowAll()

	return msg
}

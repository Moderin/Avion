package main

import (
	"gtkMod"
	"utilites"
	"encoding/json"

	gotox "github.com/codedust/go-tox"
	"github.com/gotk3/gotk3/gtk"
	"github.com/skratchdot/open-golang/open"
)

type transferMessage struct {
	messageBase
	transfer     *FileTransfer
	downloaded   bool
	messageLabel *gtk.Label
}

type transferMessageJSON struct {
	Author uint32
	Path	string
}

func (m *transferMessage) SaveJSONData(enc *json.Encoder){
	data := transferMessageJSON{
		Author: m.author,
		/*Path: m.transfer.handle.Name(),*/
	}
	
	enc.Encode(data)
}

func (m *transferMessage) send(friend uint32) bool {
	fileNumber, err := tox.FileSend(friend, gotox.TOX_FILE_KIND_DATA, m.transfer.size, nil, m.transfer.name)
	if err == nil {
		transfers[fileNumber] = m.transfer
		return true
	}

	return false
}

func transferMessageNew(author uint32, transfer *FileTransfer) *transferMessage {
	m := new(transferMessage)
	m.loadBase(author)

	m.transfer = transfer

	var animation *TransferAnimation
	transfer.message = m

	if author != User {
		// if transfer comes from friend, load download animation
		m.statusIcon.SetFromFile("img/message_status/download.png")
		animation = TransferAnimationNew(0.47, 0.85, 0.07, Download, transfer)
		transfer.animation = animation
		//when clicked
		animation.eventBox.Connect("button_press_event", func() {
			if m.downloaded {
				open.Start(transfer.handle.Name())
			} else {
				if transfer != nil {
					// if file isn't downloading right now
					if transfer.handle == nil {
						transfer.startDownload(author)
					}
				}
			}
		})
		// auto-start download, if user allowed
		if config["Auto-accept files"] != nil {
			if *config["Auto-accept files"] == "true" {
				transfer.animation.icon.Show()
				defer transfer.startDownload(author)
			}
		}
	} else {
		// if it's user file, load upload button
		m.statusIcon.SetFromFile("img/message_status/upload.png")
		animation = TransferAnimationNew(0.99, 0.47, 0.0, Upload, transfer)
		transfer.animation = animation
		filePath := transfer.handle.Name()
		animation.eventBox.Connect("button_press_event", func() {
			open.Start(filePath)
		})
	}

	// create label with name
	name := utilites.Shorten(transfer.name, 30)
	utilites.ReplaceHTMLSymbols(&name)
	m.messageLabel = gtkMod.MessageLabelNew("<b>" + utilites.FileSizeToString(transfer.size) + "</b>      " + name)

	// put everything together
	m.box.Add(animation.overlay)
	m.box.Add(m.messageLabel)
	m.box.ShowAll()

	return m
}

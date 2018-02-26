package main

import (
	"gtkMod"
	"utilites"
	"encoding/json"

	gotox "github.com/codedust/go-tox"
	"github.com/gotk3/gotk3/gtk"
)

type textMessage struct {
	messageBase
	messageLabel *gtk.Label
}

type textMessageJSON struct {
	Author uint32
	Text string	
}

func (m *textMessage) SaveJSONData(enc *json.Encoder){
	text, _ := m.messageLabel.GetText();
	data := textMessageJSON{
		Author: m.author,
		Text: text,
	}
	enc.Encode(data)
}

func (m *textMessage) send(friend uint32) bool {
	text, _ := m.messageLabel.GetText()
	id, _ := tox.FriendSendMessage(friend, gotox.TOX_MESSAGE_TYPE_NORMAL, text)
	if id != 0 {
		return true
	}

	return false
}

func textMessageNew(author uint32, text string) *textMessage {
	m := new(textMessage)
	m.loadBase(author)
	utilites.ReplaceHTMLSymbols(&text)

	m.statusIcon.SetFromFile("img/message_status/sent.png")
	m.messageLabel = gtkMod.MessageLabelNew(text)
	m.box.Add(m.messageLabel)
	m.box.ShowAll()

	return m
}

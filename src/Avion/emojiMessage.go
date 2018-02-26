package main

import (
	"strings"
	"encoding/json"

	gotox "github.com/codedust/go-tox"
	"github.com/gotk3/gotk3/gtk"
)

type emojiMessage struct {
	messageBase
	emojiCode string
}

type emojiMessageJSON struct {
	Author uint32
	Code string
}

func (m *emojiMessage) SaveJSONData(enc *json.Encoder){
	data := emojiMessageJSON{
		Author: m.author,
		Code: m.emojiCode,
	}
	
	enc.Encode(data)
}

func (m *emojiMessage) send(friend uint32) bool {
	id, _ := tox.FriendSendMessage(friend, gotox.TOX_MESSAGE_TYPE_NORMAL, m.emojiCode)
	if id != 0 {
		return true
	}
	return false
}

func emojiMessageNew(author uint32, code string) *emojiMessage {
	m := new(emojiMessage)
	m.loadBase(author)

	m.emojiCode = code
	m.statusIcon.SetFromFile("img/message_status/sent.png")
	emoticon, _ := gtk.ImageNewFromFile(emoticons[strings.TrimSpace(code)])
	m.box.Add(emoticon)
	m.box.ShowAll()

	return m
}

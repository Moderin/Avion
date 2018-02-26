package main

import (
	"encoding/hex"
	"gtkMod"
	"strings"
	"utilites"


	gotox "github.com/codedust/go-tox"
	"github.com/ctcpip/notifize"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// MessagesBox contains ScrolledWindow with messages and av contantiner
type MessagesBox struct {
	friend            uint32				// whith who we are talking with?

	/* Contantiners for messages */
	scroll            *gtk.ScrolledWindow
	box               *gtk.Grid

	/* Some messages */
	typingMsg         *typingMessage
	message           message 				// used temporarly when adding msgs
	queue             []message				// messages waiting for sending
	friendMessages    []message				// stored, so we can change avatars
	userMessages      []message				// same as here ^


	/* Used for history loading */
	fileName          string				// "[toxId]-messages"
	loadHistoryButton *gtk.Button			// :D
	historyLines	  uint32				// how much we hadn't loaded yet?
	firstMsg		  *gtk.Widget 			// handle for putting history msgs


	/* AV */
	avBox				*AvBox				// will be shown after pressing cam

}

func (mBox *MessagesBox) putMsg(msg message, fromHistory bool) {
	if fromHistory {
		if mBox.firstMsg != nil {
			mBox.box.InsertNextTo(mBox.firstMsg, gtk.POS_TOP)
			mBox.box.AttachNextTo(msg.getBox(), mBox.firstMsg, gtk.POS_TOP, 1, 1)
		} else {
			mBox.box.Add(msg.getBox())
		}
	} else {
		if mBox.typingMsg.isVisible {
			mBox.box.InsertNextTo(mBox.typingMsg.getBox(), gtk.POS_TOP)
			mBox.box.AttachNextTo(msg.getBox(), mBox.typingMsg.getBox(), gtk.POS_TOP, 1, 1)
		} else {
			mBox.box.Add(msg.getBox())
		}
	}
}

// AddUserMessage creates and shows empty message with entry
func (mBox *MessagesBox) AddUserMessage(msg message, fromHistory bool) {
	mBox.message = msg
	mBox.userMessages = append(mBox.userMessages, mBox.message)
	mBox.putMsg(msg, fromHistory)
	mBox.sizeAllocate()
}


// AddFriendMessage creates and shows a message from friend or file from us
func (mBox *MessagesBox) AddFriendMessage(msg string, friendNm uint32, transfer *FileTransfer, fromHistory bool) {
	// Replace < > &  with its HTML names, to pass Pango Markup
	utilites.ReplaceHTMLSymbols(&msg)

	// check if window is visible and display notification if isn't
	if win.IsActive() == false && activeStatus != gotox.TOX_USERSTATUS_BUSY {
		friendName, _ := tox.FriendGetName(friendNm)
		wd := utilites.GetWorkingDirectory()
		if transfer == nil {
			notifize.Display(friendName, utilites.Shorten(msg, 30), false, wd+"/img/mail-message-new.xpm")
		} else {
			notifize.Display(friendName, "File received", false, wd+"/img/file.xpm")
		}
	}

	// if this msgBox isn't visible, display indicator on friends avatar
	if activeMsgBox != mBox || settingsBox.active {
		contacts[friendNm].avatar.SetNewMessage(true)
	}


	// recognize type of data, create proper message and put it in msgBox
	if transfer != nil {
		mBox.message = transferMessageNew(mBox.friend, transfer)
	} else {


		// add it
		if emoticons[strings.TrimSpace(msg)] != "" {
			mBox.message = emojiMessageNew(mBox.friend, msg)
		} else {
			mBox.message = textMessageNew(mBox.friend, msg)
		}
	}
	mBox.friendMessages = append(mBox.friendMessages, mBox.message)

	mBox.putMsg(mBox.message, fromHistory)

	mBox.sizeAllocate()

	if(!fromHistory) {
		mBox.save(mBox.message)
	}
}

/*  basic messagesBox structure:

			messagesBox
			|		|
	  		|		ScrooledWindow
			|		|
			|		Grid -- Messages
			|
			Grid -- AV
*/


// Init makes the MessagesBox
func (mBox *MessagesBox) Init(friend uint32) *gtk.ScrolledWindow {

	mBox.friend = friend
	mBox.scroll, _ = gtk.ScrolledWindowNew(nil, nil)
	mBox.box, _ = gtk.GridNew()
	mBox.box.SetOrientation(gtk.ORIENTATION_VERTICAL)
	mBox.box.SetHExpand(true)

	// save the name of messages file
	publicKey, _ := tox.FriendGetPublickey(mBox.friend)
	mBox.fileName = hex.EncodeToString(publicKey) + "-messages"

	
	mBox.historyLines = utilites.CountLines(mBox.fileName)
	
	// if history isn't empty
	if mBox.historyLines > 0 {
		// make an button for chat history loading
		mBox.loadHistoryButton, _ = gtkMod.IconButtonNew("img/upload.png")
		style, _ := mBox.loadHistoryButton.GetStyleContext()
		style.AddClass("circular")
		mBox.loadHistoryButton.SetTooltipText("Load earlier conversation")
		mBox.loadHistoryButton.SetHAlign(gtk.ALIGN_CENTER)
		mBox.loadHistoryButton.SetBorderWidth(20)
		mBox.loadHistoryButton.Connect("clicked", mBox.loadMessages)
		mBox.loadHistoryButton.SetHExpand(true)

		mBox.box.Add(mBox.loadHistoryButton)
	}

		// add typing message
	mBox.typingMsg = typingMessageNew()

	mBox.scroll.Add(mBox.box)



	return mBox.scroll
}

func (mBox *MessagesBox) sizeAllocate() {
	mBox.scroll.CheckResize()
	vadj := mBox.scroll.GetVAdjustment()
	vadj.SetValue(vadj.GetUpper() - vadj.GetPageSize())
}


/*
	Hide or show typing message if needed
	and scroll messagesBox down
*/

func (mBox *MessagesBox) callback(widget *gtk.Window, event *gdk.Event) {
	keyEvent := gdk.EventKey{event}

	if gtkMod.Keys[keyEvent.KeyVal()] == "esc" {
		mBox.typingMsg.remove()
		return
	}

	if !mBox.typingMsg.isVisible {
		mBox.typingMsg.restore()
	}

	if !mBox.typingMsg.entry.HasFocus() {
		mBox.typingMsg.entry.GrabFocus()
	}

	mBox.sizeAllocate()
}

package main

import (
	"color"
	"encoding/hex"
	"flag"
	"fmt"
	"gtkMod"
	"io/ioutil"
	"log"
	"utilites"

	gotox "github.com/codedust/go-tox"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	// DefaultThemePath is a path for CSS theme file
	DefaultThemePath = "./themes/Arc-Darker/gtk-3.0/gtk.css"
)

var activeStatus gotox.ToxUserStatus

var activeMsgBox *MessagesBox
var mainContantiner *gtk.Box
var topBar *gtk.HeaderBar
var msgBoxes = make(map[uint32]*MessagesBox)
var contacts = make(map[uint32]*Contact)
var contactsBox *gtk.ListBox
var popover *gtk.Popover
var addFriendRBox *gtk.ListBoxRow
var settingsBox *SettingsBox
var titleBox *TitleBox
var themeProvider *gtk.CssProvider

var win *gtk.Window

func main() {

	minimized := flag.Bool("minimized", false, "start minimized")
	flag.Parse()

	
	/*
	 *	Say Welcome and print some debug info
	 */
	
	fmt.Println(color.Green("Welcome!"), "\t Avion 1.0\n")
	fmt.Println(color.Orange("Working directory:\t"), utilites.GetWorkingDirectory())
	fmt.Println(color.Orange("Downloads directory:\t"), utilites.GetDownloadsDirectory()+"\n")
	
	
	msgBoxes = make(map[uint32]*MessagesBox)
	contacts = make(map[uint32]*Contact)

	loadConfig()

	gtk.Init(nil)

	// make a window
	var err error
	win, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal(color.Red("Can't create a window: "), err)
	}

	// load CSS theme
	themeProvider, _ = gtk.CssProviderNew()
	err = themeProvider.LoadFromPath(DefaultThemePath)
	if err != nil {
		fmt.Println(color.Red("Cannot load theme :( "), err)
	}
	screen, _ := win.GetScreen()

	if config["Use system theme"] != nil && *config["Use system theme"] == "true" {
		gtk.AddProviderForScreen(screen, themeProvider, gtk.STYLE_PROVIDER_PRIORITY_THEME)
	} else {
		gtk.AddProviderForScreen(screen, themeProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	}

	loadTox()
	// header bar
	topBar, _ = gtk.HeaderBarNew()
	topBar.SetTitle("Avion")

	// settings button
	settingsBox = SettingsBoxNew()
	settingsButton, _ := gtkMod.IconButtonNew("img/gear.png")
	settingsButton.SetTooltipText("Settings")
	settingsButton.Connect("clicked", func() {
		if settingsBox.active {
			settingsBox.close()
		} else {
			settingsBox.open()
		}
	})
	topBar.Add(settingsButton)
	
	videoButton, _ := gtkMod.IconButtonNew("img/camera.png")
	style, _ := videoButton.GetStyleContext()
	style.AddClass("circular")

//  not yet. Maybe even never :'(
//	topBar.PackEnd(videoButton)
	// set up titleBox
	titleBox = titleBoxNew()

	topBar.SetShowCloseButton(true)

	//-------------------------------------------------------------
	//		Add button and menu for changing status
	//------------------------------------------------------------

	statusMenu, _ := gtk.RevealerNew()
	statusMenu.SetTransitionType(gtk.REVEALER_TRANSITION_TYPE_SLIDE_RIGHT)

	busyIcon, _ := gtk.ImageNewFromFile("img/user_status/busy.png")
	awayIcon, _ := gtk.ImageNewFromFile("img/user_status/away.png")

	statusButton, avaliableIcon := gtkMod.IconButtonNew("img/user_status/avaliable.png")
	statusButton.Connect("clicked", func() {
		if statusMenu.GetRevealChild() {
			statusMenu.SetRevealChild(false)
		} else {
			statusMenu.SetRevealChild(true)
		}
	})

	statusMenuContantiner, _ := gtk.GridNew()

	avaliableButton, _ := gtkMod.IconButtonNew("img/user_status/avaliable.png")
	avaliableButton.Connect("clicked", func() {
		activeStatus = gotox.TOX_USERSTATUS_NONE
		statusMenu.SetRevealChild(false)
		statusButton.SetImage(avaliableIcon)
		statusMenuContantiner.ShowAll()
		avaliableButton.Hide()
	})

	busyButton, _ := gtkMod.IconButtonNew("img/user_status/busy.png")
	busyButton.Connect("clicked", func() {
		activeStatus = gotox.TOX_USERSTATUS_BUSY
		statusMenu.SetRevealChild(false)
		statusButton.SetImage(busyIcon)
		statusMenuContantiner.ShowAll()
		busyButton.Hide()
	})

	awayButton, _ := gtkMod.IconButtonNew("img/user_status/away.png")
	awayButton.Connect("clicked", func() {
		activeStatus = gotox.TOX_USERSTATUS_AWAY
		statusMenu.SetRevealChild(false)
		statusButton.SetImage(awayIcon)
		statusMenuContantiner.ShowAll()
		awayButton.Hide()
	})

	statusMenuContantiner.Add(avaliableButton)
	statusMenuContantiner.Add(awayButton)
	statusMenuContantiner.Add(busyButton)

	statusMenu.Add(statusMenuContantiner)

	topBar.Add(statusButton)
	topBar.Add(statusMenu)

	// window body
	mainContantiner, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	contactsBox, _ = gtk.ListBoxNew()
	contactsScroll, _ := gtk.ScrolledWindowNew(nil, nil)
	contactsScroll.Add(contactsBox)

	//-------------------------------------------------------------
	//		Add button and menu to add new friends
	//------------------------------------------------------------

	addFriendButton, _ := gtkMod.IconButtonNew("img/add.png")
	addFriendButton.SetBorderWidth(10)
	addFriendRBox, _ = gtk.ListBoxRowNew()
	addFriendRBox.Add(addFriendButton)
	contactsBox.Add(addFriendRBox)

	popover, _ = gtk.PopoverNew(addFriendButton)
	entry, _ := gtk.EntryNew()
	entry.Connect("activate", func() {
		stringKey, _ := entry.GetText()
		friendKey, errr := hex.DecodeString(stringKey)
		friendN, errr := tox.FriendAdd(friendKey, "Hi, Avion here")
		if errr == nil {
			addFriend(friendN)
			mainContantiner.ShowAll()
		} else {
			fmt.Println(color.Red("Cannot add friend :("))
		}
	})

	label, _ := gtk.LabelNew("Friend Tox ID:")

	friendAddBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)

	friendAddBox.SetBorderWidth(10)

	friendAddBox.Add(label)
	friendAddBox.Add(entry)

	popover.Add(friendAddBox)

	addFriendButton.Connect("clicked", func() {
		if popover.GetVisible() {
			popover.Hide()
		} else {
			popover.ShowAll()
		}
	})

	//-------------------------------------------------------------
	//		Load Tox and list of friends
	//------------------------------------------------------------

	friends, _ := tox.SelfGetFriendlist()

	for friend := range friends {
		addFriend(uint32(friend))
	}

	contactsScroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	mainContantiner.Add(contactsScroll)

	win.Connect("destroy", func() {

		tox.Kill()
		gtk.MainQuit()

	})

	win.Connect("key-press-event", func(win *gtk.Window, event *gdk.Event) {
		if activeMsgBox != nil && settingsBox.active == false && popover.GetVisible() == false && titleBox.nameEntry.IsFocus() == false && titleBox.statusEntry.IsFocus() == false {
			activeMsgBox.callback(win, event)
		}
	})

	win.Add(mainContantiner)
	win.SetTitlebar(topBar)
	win.SetDefaultSize(860, 500)
	win.SetSizeRequest(860, 500)
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.SetIconFromFile("img/sTox.png")
	if *minimized == true {
		win.Iconify()
	}
	win.ShowAll()
	avaliableButton.Hide()

	glib.TimeoutAdd(100, func() bool {
		tox.Iterate()
		//save data
		data, _ := tox.GetSavedata()
		err = ioutil.WriteFile("toxConfig", data, 0644)
		if err != nil {
			log.Fatal(color.Red("Cannot save tox configuration: "), err)
		}
		tox.SelfSetStatus(activeStatus)

		return true
	})

	gtk.Main()
}

func addFriend(friend uint32) {
	contact := new(Contact)

	friendName, _ := tox.FriendGetName(friend)
	friendStatusMsg, _ := tox.FriendGetStatusMessage(friend)

	contactsBox.Remove(addFriendRBox)
	contacts[friend] = contact
	contactsBox.Add(contact.Init(friend, friendName, friendStatusMsg))
	contactsBox.Add(addFriendRBox)

	msgBoxes[friend] = contact.messagesBox

	// later - default conversation
	if friend == 0 {
		activeMsgBox = contact.messagesBox
		mainContantiner.PackEnd(activeMsgBox.scroll, true, true, 0)
	}
}

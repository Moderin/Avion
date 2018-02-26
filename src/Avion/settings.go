package main

import (
	"color"
	"encoding/hex"
	"fmt"
	"image"
	"os"
	"runtime"
	"utilites"

	gotox "github.com/codedust/go-tox"
	"github.com/gotk3/gotk3/gtk"
	"github.com/nfnt/resize"

	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"math/rand"
)

//SettingsBox is a grid with settings
type SettingsBox struct {
	active bool
	grid   *gtk.Grid
}

func (s *SettingsBox) open() {
	if activeMsgBox != nil {
		mainContantiner.Remove(activeMsgBox.scroll)
	}
	mainContantiner.PackEnd(s.grid, true, true, 0)
	s.grid.ShowAll()
	s.active = true
}

func (s *SettingsBox) close() {
	mainContantiner.Remove(s.grid)
	if activeMsgBox != nil {
		mainContantiner.PackEnd(activeMsgBox.scroll, true, true, 0)
	}
	s.active = false
}

var row = 0

func (s *SettingsBox) addHeader(text string) {
	label, _ := gtk.LabelNew("")
	label.SetMarkup("<big>" + text + "</big>")
	label.SetHAlign(gtk.ALIGN_START)
	s.grid.Attach(label, 1, row, 1, 1)
	row++
}

func (s *SettingsBox) addSwitch(text string, f func()) {
	label, _ := gtk.LabelNew(text)
	s.grid.Attach(label, 1, row, 1, 1)
	row++
	label.SetHAlign(gtk.ALIGN_END)

	switcher, _ := gtk.SwitchNew()
	if config[text] != nil {
		if *config[text] == "true" {
			switcher.SetActive(true)
		}
	}
	switcher.SetHAlign(gtk.ALIGN_START)
	s.grid.AttachNextTo(switcher, label, gtk.POS_RIGHT, 1, 1)

	switcher.Connect("state-set", func() {
		var value string
		if switcher.GetActive() == true {
			value = "true"
		} else {
			value = "false"
		}
		config[text] = &value
		configFileSave()
		f()
	})
}

/*SettingsBoxNew makes a grid with settings,
whitch is opened when user clicks settings button.
*/
func SettingsBoxNew() *SettingsBox {
	s := new(SettingsBox)
	s.grid, _ = gtk.GridNew()
	s.grid.SetBorderWidth(50)
	s.grid.SetColumnSpacing(20)
	s.grid.SetRowSpacing(15)

	/*
		Display User's Tox ID and provide button to change NoSpam
	*/
	s.addHeader("Your <span color='#e77d00'>Tox ID</span>:")

	noSpamButton, _ := gtk.ButtonNewWithLabel("New NoSpam")
	noSpamButton.SetHAlign(gtk.ALIGN_END)
	s.grid.Attach(noSpamButton, 2, row-1, 1, 1)

	label, _ := gtk.LabelNew("")

	// display Tox ID, coloring NoSpam
	updateLabelMarkup := func() {
		a, _ := tox.SelfGetAddress()
		sa := hex.EncodeToString(a)
		label.SetMarkup("<tt><small>" +
			sa[:64] + // public key
			"<span color='#56c719'>" + sa[64:72] + "</span>" + // NoSpam
			sa[72:] + "</small></tt>") // checksum
	}

	noSpamButton.Connect("clicked", func() {
		tox.SelfSetNospam(rand.Uint32())
		updateLabelMarkup()
	})

	updateLabelMarkup()
	label.SetSelectable(true)
	label.SetMarginBottom(30)
	s.grid.Attach(label, 1, row, 2, 1)
	row++

	/*
		Display profile picture, and provide option to change it
	*/
	s.addHeader("Profile picture")
	avatarEBox, _ := gtk.EventBoxNew()
	avatar := utilites.AvatarDrNew(config["avatar"], 32, 1)
	avatarEBox.Add(avatar)
	avatar.SetHAlign(gtk.ALIGN_CENTER)
	avatar.SetMarginBottom(10)

	avatarEBox.Connect("button_press_event", func() {
		//	set up file chooser and images filter
		dialog, _ := gtk.FileChooserDialogNewWith2Buttons("Choose new profile picture", win, gtk.FILE_CHOOSER_ACTION_OPEN, "That's it!", gtk.RESPONSE_OK, "Cancel", gtk.RESPONSE_CANCEL)
		filter, _ := gtk.FileFilterNew()
		filter.SetName("Images")
		filter.AddPattern("*.png")
		filter.AddPattern("*.jpg")
		filter.AddPattern("*.gif")
		dialog.AddFilter(filter)
		if dialog.Run() == int(gtk.RESPONSE_OK) {
			// try to open selected image
			name := dialog.GetFilename()
			newAvatar, err := os.Open(name)
			if err != nil {
				fmt.Println(color.Red("Cannot open selected file: "), err)
				return
			}
			// decode image's config, to check its size
			aconfig, _, _ := image.DecodeConfig(newAvatar)
			newAvatar.Close()

			// again open image
			newAvatar, _ = os.Open(name)
			img, _, errr := image.Decode(newAvatar)

			// if size > desired, scale it
			if errr == nil {
				if aconfig.Height > 150 || aconfig.Width > 150 {
					if aconfig.Height > aconfig.Width {
						img = resize.Resize(0, 150, img, resize.MitchellNetravali)
					} else {
						img = resize.Resize(150, 0, img, resize.MitchellNetravali)
					}
				}

				//	export to png
				outImg, _ := os.Create("img/userAvatar.png")
				defer outImg.Close()
				png.Encode(outImg, img)

				// change config and update all avatars we have opened
				*config["avatar"] = "img/userAvatar.png"
				avatar.UpdateSource()
				titleBox.picture.UpdateSource()

				for _, msgBox := range msgBoxes {
					for _, msg := range msgBox.userMessages {
						msg.updateFace()
					}
				}

				// check whitch friends are avaliable and send them new avatar.
				friends, _ := tox.SelfGetFriendlist()

				for friend := range friends {
					if status, _ := tox.FriendGetConnectionStatus(uint32(friend)); status != gotox.TOX_CONNECTION_NONE {
						sendAvatar(uint32(friend))
					}
				}
			} else {
				fmt.Println(color.Red("Cannot decode file: "), errr)
			}
			newAvatar.Close()
		}
		dialog.Destroy()
	})
	s.grid.Attach(avatarEBox, 1, row, 2, 1)
	row += 2

	/*
		Add switches and connect them
	*/
	s.addHeader("Settings")

	s.addSwitch("Auto-accept files", func() {
	})

	s.addSwitch("Use system theme", func() {
		screen, _ := win.GetScreen()
		gtk.RemoveProviderForScreen(screen, themeProvider)
		if config["Use system theme"] != nil && *config["Use system theme"] == "true" {
			gtk.AddProviderForScreen(screen, themeProvider, gtk.STYLE_PROVIDER_PRIORITY_THEME)
		} else {
			gtk.AddProviderForScreen(screen, themeProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
		}

	})

	if runtime.GOOS == "linux" { // so we can process on ~/.config/autostart
		s.addSwitch("Run on startup", func() {

			if config["Run on startup"] != nil && *config["Run on startup"] == "true" {
				// get current working directory
				wd := utilites.GetWorkingDirectory()

				// prepare data to write (I know, it looks strange, we cannot insert \t)
				data := []byte(`[Desktop Entry]
Type=Application
Name=Avion
Exec=bash -c "cd ` + wd + `;./main -minimized"`)

				// add .desktop file
				err := ioutil.WriteFile(os.Getenv("HOME")+"/.config/autostart/Avion.desktop", data, 0644)
				if err != nil {
					fmt.Println(color.Red("Cannot create desktop file"))
				}
			} else {
				// remove .desktop file
				os.Remove(os.Getenv("HOME") + "/.config/autostart/Avion.desktop")
			}

		})
	}
	return s
}

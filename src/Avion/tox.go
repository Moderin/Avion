package main

import (
	"color"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"utilites"

	gotox "github.com/codedust/go-tox"
	//gotox "github.com/kitech/go-toxcore"
)

// FileTransfer is struct for file transfers
type FileTransfer struct {
	number uint32
	size   uint64
	handle *os.File

	name            string
	transferredSize uint64
	message         *transferMessage
	animation       *TransferAnimation
}

func sendAvatar(friend uint32) {
	avatar, _ := os.Open(*config["avatar"])
	stat, _ := avatar.Stat()
	fileNumber, _ := tox.FileSend(friend, gotox.TOX_FILE_KIND_AVATAR, uint64(stat.Size()), nil, "avatar.png")
	transfers[fileNumber] = &FileTransfer{size: uint64(stat.Size()), handle: avatar}
}

func (transfer *FileTransfer) startDownload(friend uint32) {
	var err error
	friendName, _ := tox.FriendGetName(friend)
	path := utilites.GetDownloadsDirectory() + "/Avion/" + friendName

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		fmt.Println(utilites.LogStampErr(), "Error creating download directory: ", err)
	}

	transfer.handle, err = os.Create(path + "/" + transfer.name)
	if err != nil {
		fmt.Println(utilites.LogStampErr(), "Error creating download file: ", err)
	}
	tox.FileControl(friend, transfer.number, gotox.TOX_FILE_CONTROL_RESUME)
	transfer.animation.eventBox.Remove(transfer.animation.icon)
	transfer.animation.eventBox.Add(transfer.animation.label)
	transfer.animation.label.Show()
}

// Map of active file transfers
var transfers = make(map[uint32]*FileTransfer)

var tox *gotox.Tox

func loadTox() {
	savedData, err := ioutil.ReadFile("toxConfig")
	var options *gotox.Options
	if err == nil {
		options = &gotox.Options{
			IPv6Enabled:  true,
			UDPEnabled:   true,
			ProxyType:    gotox.TOX_PROXY_TYPE_NONE,
			ProxyHost:    "127.0.0.1",
			ProxyPort:    5555,
			StartPort:    0,
			EndPort:      0,
			TcpPort:      0,
			SaveDataType: gotox.TOX_SAVEDATA_TYPE_TOX_SAVE,
			SaveData:     savedData}
	} else {
		options = nil
	}

	/*
		Let's try to create Tox instance.
	*/

	tox, err = gotox.New(options)





	/*		OK, now there are 2 options.
			1. It has loaded without errors
			2. It hasn't, propably toxConfig file is broken
	*/



	if err != nil {
		// Option 2
		fmt.Println(color.Orange("Warning: "), " I can't load Your Tox instance from toxConfig file")

		/*	Right, now we need to check if backup file exists, and load it
			If it doesn't let's just exit.
		*/

		if _, err := os.Stat("toxConfig.backup"); err == nil {
			fmt.Println("I' ve found a backup file, I'll try to load from it \n")

			// Let's copy toxConfig.backup to toxConfig
			utilites.CopyFile("toxConfig.backup", "toxConfig")
		}


		/*	Now we have to re-init everything :D */
		defer loadTox()
		return

	} else {
		// Option 1, so everthing is right. Let's make a backup now
		utilites.CopyFile("toxConfig", "toxConfig.backup")
	}






	tox.SelfSetName(*config["name"])
	if config["status"] != nil {
		tox.SelfSetStatusMessage(*config["status"])
	}

	// get self Tox ID and print it
	a, err := tox.SelfGetAddress()
	if err != nil {
		log.Fatal(color.Red("Error getting adress: "), err)
	}
	fmt.Println(color.Green("Your Tox ID: "), "\n\t", hex.EncodeToString(a), "\n")

	activeStatus = gotox.TOX_USERSTATUS_NONE
	//tox.SelfSetStatus(gotox.TOX_USERSTATUS_NONE)

	// load tox node(s)
	pubkey, _ := hex.DecodeString("788236D34978D1D5BD822F0A5BEBD2C53C64CC31CD3149350EE27D4D9A2F9B6B")
	err = tox.Bootstrap("178.62.250.138", 3389, pubkey)

	pubkey, _ = hex.DecodeString("461FA3776EF0FA655F1A05477DF1B3B614F7D6B124F7DB1DD4FE3C08B03B640F")
	err = tox.Bootstrap("130.133.110.14", 33445, pubkey)

	pubkey, _ = hex.DecodeString("5823FB947FF24CF83DDFAC3F3BAA18F96EA2018B16CC08429CB97FA502F40C23")
	err = tox.Bootstrap("95.215.46.114", 33445, pubkey)

	if err != nil {
		panic(err)
	}

	//---------------------------------------
	// 				Tox Callbacks
	// --------------------------------------

	tox.CallbackFriendRequest(func(t *gotox.Tox, key []byte, message string) {
		fmt.Println(utilites.LogStamp()+color.Green("Friend request from:\n"), hex.EncodeToString(key))

		friend, _ := t.FriendAddNorequest(key)

		addFriend(friend)
		mainContantiner.ShowAll()
	})

	tox.CallbackSelfConnectionStatusChanges(func(t *gotox.Tox, connectionStatus gotox.ToxConnection) {
	})

	tox.CallbackFriendMessage(func(t *gotox.Tox, friend uint32, messageType gotox.ToxMessageType, message string) {
		if msgBoxes[friend] != nil {
			msgBoxes[friend].AddFriendMessage(message, friend, nil, false)
		}
	})

	tox.CallbackFriendNameChanges(func(t *gotox.Tox, friend uint32, name string) {
		if contacts[friend] != nil {
			contacts[friend].UpdateName(name)
		}
	})

	tox.CallbackFriendStatusChanges(func(t *gotox.Tox, friend uint32, status gotox.ToxUserStatus) {
		if contacts[friend] != nil {
			contacts[friend].UpdateStatus(status)
		}
	})

	tox.CallbackFriendConnectionStatusChanges(func(t *gotox.Tox, friend uint32, status gotox.ToxConnection) {
		if contacts[friend] != nil {
			if status == gotox.TOX_CONNECTION_NONE {
				contacts[friend].avatar.SetStatus(offline)

				// maybe not yet?
				// remove unfinished transfers
				/*	for _, msg := range msgBoxes[friend].friendMessages {
						if msg.transfer != nil && msg.downloaded == false {
							t.FileControl(friend, msg.transfer.number, gotox.TOX_FILE_CONTROL_CANCEL)
						}
					}

					for _, msg := range msgBoxes[friend].userMessages {
						if msg.transfer != nil && msg.downloaded == false {
							t.FileControl(friend, msg.transfer.number, gotox.TOX_FILE_CONTROL_CANCEL)
						}
					}*/
			} else {
				sendAvatar(friend)
				//	sens waiting messages if any
				for _, message := range msgBoxes[friend].queue {
					done := message.send(friend)
					if done {
						message.SetState(Sent)
						msgBoxes[friend].queue = msgBoxes[friend].queue[1:]
					}
				}
			}
		}

	})

	tox.CallbackFriendStatusMessageChanges(func(t *gotox.Tox, friend uint32, status string) {
		if contacts[friend] != nil {
			contacts[friend].UpdateStatusMsg(status)
		}
	})

	tox.CallbackFileRecv(func(t *gotox.Tox, friend, fileNumber uint32, kind gotox.ToxFileKind, filesize uint64, filename string) {
		if kind == gotox.TOX_FILE_KIND_AVATAR {
			if filesize > 65536 {
				t.FileControl(friend, fileNumber, gotox.TOX_FILE_CONTROL_CANCEL)
				return
			}

			publicKey, _ := t.FriendGetPublickey(friend)
			file, _ := os.Create(hex.EncodeToString(publicKey) + ".png")

			transfers[fileNumber] = &FileTransfer{size: filesize, handle: file}
			t.FileControl(friend, fileNumber, gotox.TOX_FILE_CONTROL_RESUME)
		} else {
			transfers[fileNumber] = &FileTransfer{size: filesize, name: filename, number: fileNumber}
			msgBoxes[friend].AddFriendMessage("", friend, transfers[fileNumber], false)
		}
	})

	tox.CallbackFileRecvControl(func(t *gotox.Tox, friend uint32, fileNumber uint32, fileControl gotox.ToxFileControl) {
		switch fileControl {
		case gotox.TOX_FILE_CONTROL_CANCEL:
			transfers[fileNumber].handle.Sync()
			transfers[fileNumber].handle.Close()
			if transfers[fileNumber].message != nil {
				msgBoxes[friend].box.Remove(transfers[fileNumber].message.box)
			}
			delete(transfers, fileNumber)
		}
	})
	/* for sending */
	tox.CallbackFileChunkRequest(func(t *gotox.Tox, friend uint32, fileNumber uint32, position uint64, length uint64) {
		if transfers[fileNumber] != nil {
			if length == 0 {
				transfers[fileNumber].handle.Close()
				if animation := transfers[fileNumber].animation; animation != nil {
					animation.Done()
				}
				friendName, _ := tox.FriendGetName(friend)
				fmt.Println(utilites.LogStamp(), transfers[fileNumber].handle.Name(), " sent to "+friendName)
				delete(transfers, fileNumber)

			} else {
				data := make([]byte, length)
				if animation := transfers[fileNumber].animation; animation != nil {
					transfers[fileNumber].transferredSize = position + length
					animation.circle.QueueDraw()
				}
				transfers[fileNumber].handle.ReadAt(data, int64(position))
				t.FileSendChunk(friend, fileNumber, position, data)
			}
		}
	})
	//*/

	tox.CallbackFileRecvChunk(func(t *gotox.Tox, friend uint32, fileNumber uint32, position uint64, data []byte) {
		if len(data) == 0 {

		} else {
			transfers[fileNumber].handle.WriteAt(data, int64(position))

			// if transfer completed
			if position+uint64(len(data)) >= transfers[fileNumber].size {
				transfers[fileNumber].handle.Sync()
				transfers[fileNumber].handle.Close()
				// if file isn't avatar
				if transfers[fileNumber].animation != nil {
					transfers[fileNumber].animation.Done()
					transfers[fileNumber].animation.eventBox.Remove(transfers[fileNumber].animation.label)
					transfers[fileNumber].animation.eventBox.Add(transfers[fileNumber].animation.icon)

					dA := transfers[fileNumber].animation.circle

					transfers[fileNumber].message.downloaded = true

					delete(transfers, fileNumber)
					dA.QueueDraw()
				} else {
					friendName, _ := tox.FriendGetName(friend)
					fmt.Println(utilites.LogStamp(), "Received avatar from"+friendName)
					contacts[friend].avatar.picture.UpdateSource()
					for _, msg := range msgBoxes[friend].friendMessages {
						msg.updateFace()
					}
					delete(transfers, fileNumber)
				}

			} else {
				transfers[fileNumber].transferredSize = position + uint64(len(data))
				if transfers[fileNumber].animation != nil {
					transfers[fileNumber].animation.circle.QueueDraw()
					if transfers[fileNumber].animation.label != nil {
						progress := float64(transfers[fileNumber].transferredSize) / float64(transfers[fileNumber].size)
						transfers[fileNumber].animation.label.SetLabel(strconv.Itoa(int(progress * 100)))
					}
				}

			}
		}
	})

	err = tox.Iterate()
	if err != nil {
		log.Fatal(color.Red("Error connecting tox: "), err)
	}

}

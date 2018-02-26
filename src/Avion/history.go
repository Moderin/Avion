package main

import (
	"encoding/json"
	"fmt"
	"os"
	"utilites"
)

const messagesToLoad = 15

/*
 *      Load last 15 messages from file
 */

func (mBox *MessagesBox) loadMessages() {
	// save first msg
	mBox.firstMsg, _ = mBox.box.GetChildAt(0, 1)


	// open messages file

	file, err := os.Open(mBox.fileName)
	if err != nil {
		fmt.Println(utilites.LogStampErr(), "Cannot open messages file: ", err)
		return
	}
	defer file.Close()



	dec := json.NewDecoder(file)

	var lineNumber uint32
	var decodedN uint32


	// go throw the lines of messages adding proper messages on screen

	for dec.More() {

		lineNumber++


		// check if actual line contains too old messages, and move on if true

		if lineNumber <= (mBox.historyLines-messagesToLoad) && mBox.historyLines > messagesToLoad {
			dec.Decode(nil)
			continue
		}


		// decode json message

		m := make(map[string]interface{})
		dec.Decode(&m)


		/*  this gives us ability to do not put
			same messages when msgLoadBtn pressed again */

		decodedN += 1


		/* 	Now, depending of the message type
			add proper msg to msgBox */

		if m["Text"] != nil {
			if m["Author"].(float64) == User {
				msg := textMessageNew(User, m["Text"].(string))
				mBox.AddUserMessage(msg, true)
			} else {
				mBox.AddFriendMessage(m["Text"].(string), uint32(m["Author"].(float64)), nil, true)
			}
		} else if m["Code"] != nil {
			if m["Author"].(float64) == User {
				msg := emojiMessageNew(User, m["Code"].(string))
				mBox.AddUserMessage(msg, true)
			} else {
				mBox.AddFriendMessage(m["Code"].(string), uint32(m["Author"].(float64)), nil, true)
			}
		} else if m["Path"] != nil {
			if m["Author"].(float64) == User {

			} else {
				// TODO: file transfer
			}
		}


		/* 	check if we loaded all messages that we had to load,
			if true end the loop */

		if lineNumber >= mBox.historyLines {
			mBox.historyLines -= decodedN
			decodedN = 0

			// if all loaded, remove button for history loading
			if mBox.historyLines == 0 {
				mBox.box.Remove(mBox.loadHistoryButton)
			}

			// end loop, messages loaded
			return
		}

	}

}

/*
 *  Open toxID-messages file and append to it JSON of msg
 */

func (mBox *MessagesBox) save(msg message) {
	// open file with APPEND flag
	file, err := os.OpenFile(mBox.fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(utilites.LogStampErr(), "Cannot open messages file: ", err)
		return
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	msg.SaveJSONData(enc)

}

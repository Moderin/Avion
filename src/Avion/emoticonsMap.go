package main

/*	emoticons, whitch user can send,
	simply comment to disable one. Cool, isn't it? */
var userEmoticons = [...]string{
	"😄",
	"😉",
	"😂",
	//"😎",
	"😏",
	"😟",
}

/*	Map of unicode emoticons, and PNG images, whitch user see  */
var emoticons = map[string]string{
	"😄": "img/emoji/smile.png",
	"😂": "img/emoji/joy.png",
	"😏": "img/emoji/smriking.png",
	"😟": "img/emoji/worried.png",
	"😎": "img/emoji/sunglasses.png",
	"😉": "img/emoji/winking.png",
}

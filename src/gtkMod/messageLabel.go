package gtkMod

import (
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/mvdan/xurls"
)

// MessageLabelNew returns label adjusted to use in message
func MessageLabelNew(markup string) *gtk.Label {
	label, _ := gtk.LabelNew("")
	label.SetLineWrap(true)
	label.SetLineWrapMode(pango.WRAP_CHAR)
	label.SetSelectable(true)

	/*	check if message contains link, and make it clickable if any */
	if url := xurls.Strict.FindString(markup); url != "" &&
		(strings.Contains(url, "http://") || strings.Contains(url, "https://")) {
		markup = strings.Replace(markup, url, "<a href='"+url+"'>"+url+"</a>", -1)
		label.SetTrackVisitedLinks(false)
	}
	label.SetMarkup("<big>" + markup + "</big>")

	return label
}

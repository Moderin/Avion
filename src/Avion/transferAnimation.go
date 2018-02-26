package main

import (
	"math"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

const (
	// Upload transfer mode
	Upload = true
	// Download transfer mode
	Download = false
)

// TransferAnimation is type with circle with progress animation
type TransferAnimation struct {
	overlay  *gtk.Overlay
	eventBox *gtk.EventBox

	circle *gtk.DrawingArea
	icon   *gtk.Image
	label  *gtk.Label

	done bool
}

// Done is called when fle transfer completed
func (t *TransferAnimation) Done() {
	t.icon.SetFromFile("img/ok.png")
	t.done = true
}

// TransferAnimationNew return ready to use animation
func TransferAnimationNew(r, g, b float64, mode bool, transfer *FileTransfer) *TransferAnimation {
	t := new(TransferAnimation)

	// draw a circle
	t.circle, _ = gtk.DrawingAreaNew()
	t.circle.SetSizeRequest(32, 32)
	t.circle.Connect("draw", func(dr *gtk.DrawingArea, context *cairo.Context) {
		var progress float64
		if t.done == false && transfer != nil {
			progress = float64(transfer.transferredSize) / float64(transfer.size)
		}

		// draw ring
		context.SetSourceRGB(r, g, b)
		context.SetLineWidth(1)
		context.Arc(16, 16, 13, 0, 2*math.Pi)
		context.Stroke()

		// fill it white
		context.Arc(16, 16, 12, 0, 2*math.Pi)
		context.SetSourceRGB(1, 1, 1)
		context.Fill()

		// Progress
		context.SetSourceRGBA(r, g, b, 0.3)
		context.Arc(16, 16, 12, 0, 2*math.Pi)
		context.Clip()
		context.Rectangle(4, 29-24*progress, 24, 24*progress)

		context.Clip()

		context.Paint()
	})

	t.eventBox, _ = gtk.EventBoxNew()

	if mode == Download {
		t.icon, _ = gtk.ImageNewFromFile("img/download.png")
		t.label, _ = gtk.LabelNew("0")
	} else {
		t.icon, _ = gtk.ImageNewFromFile("img/upload.png")
	}

	t.eventBox.Add(t.icon)

	t.overlay, _ = gtk.OverlayNew()
	t.overlay.Add(t.circle)
	t.overlay.AddOverlay(t.eventBox)

	return t
}

package utilites

import (
	"image"
	_ "image/png"
	"math"
	"os"
	"strconv"
	"unicode/utf8"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// FileSizeToString convert filesize to human readable string, eg. 1000 -> 1KB
func FileSizeToString(fileSize uint64) string {
	switch {
	case fileSize < 1000:
		return strconv.Itoa(int(fileSize)) + " B"
	case fileSize >= 1000 && fileSize < 1000*1000:
		return strconv.FormatFloat(float64(fileSize)/1000, 'f', 1, 64) + " KB"
	case fileSize >= 1000*1000 && fileSize < 1000*2024*1000:
		return strconv.FormatFloat(float64(fileSize)/1000/1000, 'f', 1, 64) + " MB"
	case fileSize >= 1000*1000*1000 && fileSize < 1000*1000*1000*1000:
		return strconv.FormatFloat(float64(fileSize)/1000/1000/1000, 'f', 1, 64) + " GB"
	case fileSize >= 1000*1000*1000*1000:
		return strconv.FormatFloat(float64(fileSize)/1000/1000/1000/1000, 'f', 1, 64) + " TB"
	}
	return "Err"
}

// Shorten cuts string to desired length and if cutted adds "..."
func Shorten(data string, length int) string {
	runesl := utf8.RuneCountInString(data)
	if runesl > length {
		runes := []rune(data)
		return string(runes[:length]) + "..."
	}

	return data
}

type PictureCircle struct {
	*gtk.DrawingArea
	upToDate bool
}

func (pc *PictureCircle) UpdateSource() {
	pc.upToDate = false
	pc.QueueDraw()
}

// AvatarDrNew gets image from file, crops it to circle and returns DrawingArea
func AvatarDrNew(path *string, radius float64, margin float64) *PictureCircle {
	dr, _ := gtk.DrawingAreaNew()
	pc := &PictureCircle{dr, false}
	var scale float64
	var surface *cairo.Surface
	loadSurface := func() {
		//	check if file exists
		var localPath string
		if _, err := gdk.PixbufNewFromFile(*path); err != nil {
			localPath = "img/avatar.png"
		} else {
			localPath = *path
		}

		file, _ := os.Open(localPath)
		imgCfg, _, _ := image.DecodeConfig(file)

		file.Close()

		if imgCfg.Height > imgCfg.Width {
			scale = (radius*2 + margin*2) / float64(imgCfg.Width)
		} else {
			scale = (radius*2 + margin*2) / float64(imgCfg.Height)
		}
		surface, _ = cairo.NewSurfaceFromPNG(localPath)
		surface = surface.CreateForRectangle(float64(imgCfg.Width/2)-(radius+margin/2)/scale, float64(imgCfg.Height/2)-(radius+margin/2)/scale, (2*radius+2*margin)/scale, (2*radius+2*margin)/scale)
		pc.upToDate = true
	}

	loadSurface()
	pc.Connect("draw", func(dr *gtk.DrawingArea, context *cairo.Context) {
		if pc.upToDate == false {
			loadSurface()
		}
		context.Scale(scale, scale)
		context.SetSourceSurface(surface, margin, margin)

		context.Arc((radius+margin)/scale, (radius+margin)/scale, radius/scale, 0, 2*math.Pi)
		context.Clip()

		context.Paint()
	})

	pc.SetSizeRequest(int((radius+margin)*2), int((radius+margin)*2))
	return pc
}

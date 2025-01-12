package ui

import (
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FileView struct {
	*tview.Box
	text     *tview.TextView
	image    *tview.Image
	file     string
	mimeType string
}

func NewFileView() *FileView {
	return &FileView{
		Box:   tview.NewBox(),
		text:  tview.NewTextView(),
		image: tview.NewImage(),
	}
}

func (f *FileView) Draw(screen tcell.Screen) {
	f.Box.DrawForSubclass(screen, f)

	if f.file == "" {
		f.image.SetImage(nil)
		f.text.Clear()
		return
	}

	reader, err := os.Open(f.file)
	if err != nil {
		return
	}
	defer reader.Close()

	x, y, width, height := f.GetInnerRect()
	switch f.mimeType {
	// TODO: images seem a bit slow
	case "image/png":
		f.text.SetRect(x, y, 0, 0)
		graphics, _ := png.Decode(reader)
		f.image.SetRect(x, y, width, height)
		f.image.SetImage(graphics)

	case "image/jpeg", "image/jpg":
		f.text.SetRect(x, y, 0, 0)
		graphics, _ := jpeg.Decode(reader)
		f.image.SetRect(x, y, width, height)
		f.image.SetImage(graphics)

	default:
		f.image.SetRect(x, y, 0, 0)
		f.image.SetImage(nil)
		f.text.SetRect(x, y, width, height)
		var sb strings.Builder
		io.Copy(&sb, reader) // TODO: size limit?
		f.text.SetText(sb.String())
	}

	f.text.Draw(screen)
	f.image.Draw(screen)
}

func (f *FileView) SetFile(file, mimeType string) {
	f.file = file
	f.mimeType = mimeType
}

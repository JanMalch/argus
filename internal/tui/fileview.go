package tui

import (
	"image/jpeg"
	"image/png"
	"mime"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FileView struct {
	*tview.Box
	*CodeView
	*tview.Image
}

func NewFileView() *FileView {
	return &FileView{
		Box: tview.NewBox(),
	}
}

func (v *FileView) onError(err error) {
	v.Image = nil
	if v.CodeView == nil {
		v.CodeView = NewCodeView()
	}
	v.CodeView.SetText(err.Error(), "")
}

func (v *FileView) SetFile(filename string) {
	mimeType := mime.TypeByExtension(filepath.Ext(filename))
	if mimeType == "image/jpeg" || mimeType == "image/jpg" || mimeType == "image/png" {
		file, err := os.Open(filename)
		if err != nil {
			v.onError(err)
			return
		}
		defer file.Close()
		v.CodeView = nil
		v.Image = tview.NewImage()
		if mimeType[6] == 'p' {
			graphics, err := png.Decode(file)
			if err != nil {
				v.onError(err)
				return
			}
			v.Image.SetImage(graphics)
		} else {
			photo, err := jpeg.Decode(file)
			if err != nil {
				v.onError(err)
				return
			}
			v.Image.SetImage(photo)
		}
	} else {
		v.Image = nil
		content, err := os.ReadFile(filename)
		if err != nil {
			v.onError(err)
			return
		}
		if v.CodeView == nil {
			v.CodeView = NewCodeView()
		}
		v.CodeView.SetText(string(content), mimeType)
	}
}

func (v *FileView) Focus(delegate func(p tview.Primitive)) {
	if v.CodeView != nil {
		delegate(v.CodeView)
	} else if v.Image != nil {
		delegate(v.Image)
	} else {
		delegate(v.Box)
	}
}

func (v *FileView) HasFocus() bool {
	if v.CodeView != nil {
		return v.CodeView.HasFocus()
	} else if v.Image != nil {
		return v.Image.HasFocus()
	} else {
		return v.Box.HasFocus()
	}
}

func (v *FileView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return v.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if v.CodeView != nil && v.CodeView.HasFocus() {
			if handler := v.CodeView.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
		if v.Image != nil && v.Image.HasFocus() {
			if handler := v.Image.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
	})
}

func (v *FileView) Draw(screen tcell.Screen) {
	v.Box.DrawForSubclass(screen, v)

	x, y, width, height := v.GetInnerRect()
	if v.CodeView != nil {
		v.CodeView.SetRect(x, y, width, height)
		v.CodeView.Draw(screen)
	} else if v.Image != nil {
		v.Image.SetRect(x, y, width, height)
		v.Image.Draw(screen)
	}
}

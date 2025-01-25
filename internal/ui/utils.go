package ui

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/janmalch/argus/internal/handler"
)

func (t *tui) currentExchange() *handler.Exchange {
	row, _ := t.table.GetSelection()
	return t.timeline.At(row)
}

func (t *tui) currentResponseBodyString() string {
	e := t.currentExchange()
	if e == nil {
		return ""
	}
	file := t.fileOf("res", e.Id, e.Response.Get("Content-Type"))
	reader, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer reader.Close()
	var sb strings.Builder
	_, err = io.Copy(&sb, reader)
	if err != nil {
		return ""
	}
	return sb.String()
}

func (t *tui) currentRequestBodyString() string {
	e := t.currentExchange()
	if e == nil {
		return ""
	}
	file := t.fileOf("req", e.Id, e.Request.Get("Content-Type"))
	reader, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer reader.Close()
	var sb strings.Builder
	_, err = io.Copy(&sb, reader)
	if err != nil {
		return ""
	}
	return sb.String()
}

func contentTypeOf(file string) string {
	contentType := mime.TypeByExtension(filepath.Ext(file))
	if contentType == "" {
		if f, err := os.Open(file); err == nil {
			lr := io.LimitReader(f, 512)
			if b, err := io.ReadAll(lr); err == nil {
				contentType = http.DetectContentType(b)
			}
			f.Close()
		}
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}

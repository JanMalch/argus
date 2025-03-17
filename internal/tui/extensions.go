package tui

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

var (
	extLUT = make(map[string]string)
)

func extensionByType(contentType, fallback string) string {
	switch contentType {
	case "text/plain":
		return ".txt"
	case "text/html":
		return ".html"
	}
	fb := fallback
	if fb != "" && fb[0] != '.' {
		fb = "." + fb
	}
	extension, ok := extLUT[contentType]
	if !ok {
		picked := ""
		extensions, _ := mime.ExtensionsByType(contentType)
		if len(extensions) > 0 {
			picked = extensions[0]
		}
		// store misses too
		extLUT[contentType] = picked
		extension = picked
	}
	if extension == "" {
		return fb
	}
	return extension
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

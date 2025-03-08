package main

import (
	"time"

	"github.com/janmalch/argus/internal/tui"
	"github.com/janmalch/argus/pkg/fmthttp"
	"github.com/rivo/tview"
)

func main() {
	now := time.Now()
	past := now.Add(time.Duration(-24*60+40) * time.Minute)
	future := now.Add(time.Duration(73) * time.Minute)
	headers := fmthttp.NewHeaders(
		"Accept", "application/json",
		"Content-Length", "7851213",
		"Cache-Control", "public, max-age=31536000",
		"Last-Modified", past.Format(time.RFC1123),
		"Expires", future.Format(time.RFC1123),
		"Authorization", "Basic YWRtaW46YWRtaW4=",
		"Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
	)
	v := tui.NewHeaderView()
	v.SetBorder(true).SetTitle(" Headers Example ")
	v.SetHeaders(headers)

	if err := tview.NewApplication().SetRoot(v, true).Run(); err != nil {
		panic(err)
	}
}

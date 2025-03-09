package main

import (
	"time"

	"github.com/janmalch/argus/internal/tui"
	"github.com/janmalch/argus/pkg/fmthttp"
	"github.com/rivo/tview"
)

func main() {
	params := fmthttp.NewParameters(
		"a", "a",
		"b", "Hello%2C%20World!",
		"bd", "Hello%252C%2520World!",
		"c", "c1",
		"c", "c2",

		"1", "x",
		"2", "x",
		"3", "x",
		"4", "x",
		"5", "x",
		"6", "x",
		"7", "x",
		"8", "x",
		"9", "x",
	)
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
	v := tui.NewExchangeView()
	v.SetRequest(params, headers, ".assets/example.min.json").
		SetResponse(headers, ".assets/example.min.json")
	v.SetBorder(true).
		SetTitle(" Exchange2 Example ")

	if err := tview.NewApplication().SetRoot(v, true).Run(); err != nil {
		panic(err)
	}
}

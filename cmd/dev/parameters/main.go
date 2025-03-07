package main

import (
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
	v := tui.NewParameterView()
	v.SetBorder(true).SetTitle(" Parameters Example ")
	v.SetParameters(params)

	if err := tview.NewApplication().SetRoot(v, true).Run(); err != nil {
		panic(err)
	}
}

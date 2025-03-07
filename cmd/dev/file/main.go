package main

import (
	"github.com/janmalch/argus/internal/tui"
	"github.com/rivo/tview"
)

func main() {
	v := tui.NewFileView()
	v.SetBorder(true).SetTitle(" File Example ")
	// v.SetFile(".assets/maurice-dt--NdiLqADQcU-unsplash.jpg")
	v.SetFile(".assets/example.min.json")

	if err := tview.NewApplication().SetRoot(v, true).Run(); err != nil {
		panic(err)
	}
}

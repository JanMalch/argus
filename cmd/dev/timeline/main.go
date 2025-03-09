package main

import (
	"net/url"
	"time"

	"github.com/janmalch/argus/internal/tui"
	"github.com/rivo/tview"
)

func main() {
	defaultColumns := []string{
		// Casing aswell as " " or "_" doesn't matter
		"ID",
		"start",
		"method",
		"host",
		"request target",
		"end",
		"duration",
		"status_code",
		"status_Text",
	}
	app := tview.NewApplication()
	v := tui.NewTimelineView(defaultColumns)
	u, _ := url.Parse("https://gobyexample.com/url-parsing")

	go func() {
		time.Sleep(1 * time.Second)
		app.QueueUpdateDraw(func() {
			v.AddRequest(1, time.Now(), "GET", u)
		})
	}()

	go func() {
		time.Sleep(2 * time.Second)
		app.QueueUpdateDraw(func() {
			v.AddRequest(2, time.Now(), "POST", u)
		})
	}()

	go func() {
		time.Sleep(3 * time.Second)
		app.QueueUpdateDraw(func() {
			v.AddResponse(1, time.Now(), 304, "")
		})
	}()

	go func() {
		time.Sleep(6 * time.Second)
		app.QueueUpdateDraw(func() {
			v.SetColumns([]string{"ID", "method", "request target", "status code", "duration"})
		})
	}()

	go func() {
		time.Sleep(8 * time.Second)
		app.QueueUpdateDraw(func() {
			v.SetColumns(defaultColumns)
		})
	}()

	v.SetBorder(true).
		SetTitle(" Timeline Example ")

	if err := app.SetRoot(v, true).Run(); err != nil {
		panic(err)
	}
}

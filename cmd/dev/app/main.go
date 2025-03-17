package main

import (
	"net/http"
	"os"
	"time"

	"github.com/janmalch/argus/internal/config"
	"github.com/janmalch/argus/internal/tui"
	"github.com/janmalch/argus/pkg/fmthttp"
)

func main() {
	app := tui.NewApp(".argus", config.UI{
		TimelineColumns: []string{"ID",
			"start",
			"method",
			"host",
			"request target",
			"end",
			"duration",
			"status_code",
			"status_Text",
		},
		Horizontal:   false,
		GrowTimeline: 2,
		GrowExchange: 3,
	})
	req0, _ := http.NewRequest("GET", "https://example.com/argus/0?demo=works", http.NoBody)
	f, _ := os.Open(".assets/example.json")
	defer f.Close()
	req1, _ := http.NewRequest("POST", "https://example.com/headers-and-body", f)
	req1.Header.Add("Authorization", "Basic YWRtaW46YWRtaW4=")
	req1.Header.Add("Content-Type", "application/json")
	go func() {
		time.Sleep(1000 * time.Millisecond)
		app.AddRequest(0, req0, time.Now())
	}()
	go func() {
		time.Sleep(1100 * time.Millisecond)
		app.AddRequest(1, req1, time.Now())
	}()
	go func() {
		time.Sleep(1400 * time.Millisecond)
		app.AddResponse(0, &fmthttp.Response{
			ResponseHead: fmthttp.NewResponseHead("HTTP/1.0", 200, "200 OK", http.Header{}),
		}, time.Now())
	}()

	if err := app.Run(); err != nil {
		panic(err)
	}
}

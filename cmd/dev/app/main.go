package main

import (
	"net/http"
	"os"
	"time"

	"github.com/janmalch/argus/internal/tui"
)

func main() {
	app := tui.NewApp([]string{"ID",
		"start",
		"method",
		"host",
		"request target",
		"end",
		"duration",
		"status_code",
		"status_Text",
	}, true, 2, 3)
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
		app.AddResponse(0, &http.Response{
			StatusCode: 200,
			Status:     "200 OK",
		}, time.Now())
	}()

	if err := app.Run(); err != nil {
		panic(err)
	}
}

package tui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/janmalch/argus/pkg/fmthttp"
	"github.com/rivo/tview"
)

type TuiApp struct {
	tmpDir    string
	app       *tview.Application
	container *tview.Flex
	timeline  *TimelineView
	exchange  *ExchangeView
	lut       map[uint64]*storedEntry
}

type storedEntry struct {
	parameters  fmthttp.Parameters
	reqHeaders  fmthttp.Headers
	reqBodyFile string
	resHeaders  fmthttp.Headers
	resBodyFile string
}

func NewApp(columns []string, layoutVertical bool, layoutTimeline, layoutExchange int) *TuiApp {
	app := tview.NewApplication()
	lut := make(map[uint64]*storedEntry)
	timeline := NewTimelineView(columns)
	exchange := NewExchangeView()
	timeline.SetSelectedEntryChangedFunc(func(entry *timelineEntry) {
		go app.QueueUpdateDraw(func() {
			if entry == nil {
				// FIXME: what do?
				return
			}
			e, ok := lut[entry.id]
			if ok {
				exchange.SetRequest(e.parameters, e.reqHeaders, e.reqBodyFile)
				exchange.SetResponse(e.resHeaders, e.resBodyFile)
			}
		})
	})
	timeline.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab && timeline.HasFocus() {
			app.SetFocus(exchange)
			return nil
		}
		return event
	})
	exchange.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc && exchange.HasFocus() {
			app.SetFocus(timeline)
			return nil
		}
		return event
	})

	container := tview.NewFlex()
	if layoutVertical {
		container.SetDirection(tview.FlexRow)
	} else {
		container.SetDirection(tview.FlexColumn)
	}
	container.AddItem(timeline, 0, layoutTimeline, true)
	container.AddItem(exchange, 0, layoutExchange, false)
	app.SetRoot(container, true)

	tmpDir := "temp"
	// FIXME: tmpDir, err := os.MkdirTemp("", "argus-")
	err := os.Mkdir(tmpDir, 0644)
	if err != nil {
		panic(err)
	}

	return &TuiApp{
		tmpDir:    tmpDir,
		app:       app,
		container: container,
		timeline:  timeline,
		exchange:  exchange,
		lut:       lut,
	}
}

func (a *TuiApp) Run() error {
	defer os.RemoveAll(a.tmpDir)
	return a.app.Run()
}

func (v *TuiApp) AddRequest(id uint64, req *http.Request, timestamp time.Time) string {
	reqBodyFile := ""
	if req.Body != nil && req.Body != http.NoBody {
		defer req.Body.Close()
		ext := extensionByType(req.Header.Get("Content-Type"), ".dat")
		reqBodyFile = filepath.Join(v.tmpDir, fmt.Sprintf("req_%v%s", id, ext))
		f, err := os.Create(reqBodyFile)
		if err != nil {
			reqBodyFile = ""
		} else {
			defer f.Close()
			_, err = io.Copy(f, req.Body)
			if err != nil {
				reqBodyFile = ""
			}
		}
	}
	v.lut[id] = &storedEntry{
		parameters:  fmthttp.CopyToParameters(req.URL),
		reqHeaders:  fmthttp.CopyToHeaders(req.Header),
		reqBodyFile: reqBodyFile,
	}
	go v.app.QueueUpdateDraw(func() {
		v.timeline.AddRequest(id, timestamp, req.Method, req.URL)
	})
	return reqBodyFile
}

func (v *TuiApp) AddResponse(id uint64, res *http.Response, timestamp time.Time) string {
	e, ok := v.lut[id]
	if !ok {
		// FIXME: panic?
		return ""
	}
	resBodyFile := ""
	if res.Body != nil && res.Body != http.NoBody {
		defer res.Body.Close()
		ext := extensionByType(res.Header.Get("Content-Type"), ".dat")
		resBodyFile = filepath.Join(v.tmpDir, fmt.Sprintf("res_%v%s", id, ext))
		f, err := os.Create(resBodyFile)
		if err != nil {
			resBodyFile = ""
		} else {
			defer f.Close()
			_, err = io.Copy(f, res.Body)
			if err != nil {
				resBodyFile = ""
			}
		}
	}

	e.resHeaders = fmthttp.CopyToHeaders(res.Header)
	e.resBodyFile = resBodyFile
	var statusText string
	if len(res.Status) > 4 {
		statusText = res.Status[4:]
	}
	go v.app.QueueUpdateDraw(func() {
		v.timeline.AddResponse(id, timestamp, res.StatusCode, statusText)
	})
	return resBodyFile
}

func (v *TuiApp) SetColumns(columns []string) {
	v.app.QueueUpdateDraw(func() {
		v.timeline.SetColumns(columns)
	})
}

func (v *TuiApp) SetLayout(vertical bool, timeline, exchange int) {
	v.app.QueueUpdateDraw(func() {
		if vertical {
			v.container.SetDirection(tview.FlexColumn)
		} else {
			v.container.SetDirection(tview.FlexRow)
		}
		v.container.ResizeItem(v.timeline, 0, timeline)
		v.container.ResizeItem(v.exchange, 0, exchange)
	})
}

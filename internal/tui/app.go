package tui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/janmalch/argus/internal/config"
	"github.com/janmalch/argus/pkg/fmthttp"
	"github.com/rivo/tview"
)

type TuiApp struct {
	directory string
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

func NewApp(directory string, ui config.UI) *TuiApp {
	app := tview.NewApplication()
	lut := make(map[uint64]*storedEntry)
	timeline := NewTimelineView(ui.TimelineColumns)
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
	if ui.Horizontal {
		container.SetDirection(tview.FlexColumn)
	} else {
		container.SetDirection(tview.FlexRow)
	}
	container.AddItem(timeline, 0, ui.GrowTimeline, true)
	container.AddItem(exchange, 0, ui.GrowExchange, false)
	app.SetRoot(container, true)

	tmpDir, err := os.MkdirTemp("", "argus-")
	if err != nil {
		panic(err)
	}

	return &TuiApp{
		directory: directory,
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

func (v *TuiApp) ReadFile(file string) (io.ReadCloser, string, error) {
	contentType := contentTypeOf(file)
	r, err := os.Open(file)
	return r, contentType, err
}

func (v *TuiApp) Log(id uint64, msg string) {
	now := time.Now().Format(time.RFC3339Nano)
	f, _ := os.OpenFile(filepath.Join(v.directory, "log.txt"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(fmt.Sprintf("[%s | %d] %s\n", now, id, msg))
}

func (v *TuiApp) AddRequest(id uint64, req *http.Request, timestamp time.Time) (string, error) {
	reqBodyFile := ""
	var rErr error
	if req.Body != nil && req.Body != http.NoBody {
		defer req.Body.Close()
		ext := extensionByType(req.Header.Get("Content-Type"), ".dat")
		reqBodyFile = filepath.Join(v.tmpDir, fmt.Sprintf("req_%v%s", id, ext))
		f, err := os.Create(reqBodyFile)
		if err != nil {
			reqBodyFile = ""
			rErr = err
		} else {
			defer f.Close()
			_, err = io.Copy(f, req.Body)
			if err != nil {
				reqBodyFile = ""
				rErr = err
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
	return reqBodyFile, rErr
}

func (v *TuiApp) AddResponse(id uint64, res *fmthttp.Response, timestamp time.Time) (string, error) {
	e, ok := v.lut[id]
	if !ok {
		panic(fmt.Sprintf("failed to find exchange with ID %v", id))
	}
	var rErr error
	resBodyFile := ""
	if res.Body != nil && res.Body != http.NoBody {
		defer res.Body.Close()
		ext := extensionByType(res.Headers.Get("Content-Type"), ".dat")
		resBodyFile = filepath.Join(v.tmpDir, fmt.Sprintf("res_%v%s", id, ext))
		f, err := os.Create(resBodyFile)
		if err != nil {
			resBodyFile = ""
			rErr = err
		} else {
			defer f.Close()
			_, err = io.Copy(f, res.Body)
			if err != nil {
				resBodyFile = ""
				rErr = err
			}
		}
	}

	e.resHeaders = res.Headers
	e.resBodyFile = resBodyFile
	go v.app.QueueUpdateDraw(func() {
		v.timeline.AddResponse(id, timestamp, res.StatusCode, res.StatusText, res.Headers.Get("Content-Type"))
	})
	return resBodyFile, rErr
}

func (v *TuiApp) SetUI(ui config.UI) {
	v.app.QueueUpdateDraw(func() {
		v.timeline.SetColumns(ui.TimelineColumns)
		if ui.Horizontal {
			v.container.SetDirection(tview.FlexColumn)
		} else {
			v.container.SetDirection(tview.FlexRow)
		}
		v.container.ResizeItem(v.timeline, 0, ui.GrowTimeline)
		v.container.ResizeItem(v.exchange, 0, ui.GrowExchange)
	})
}

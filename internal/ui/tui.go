package ui

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/janmalch/argus/internal/handler"
	"github.com/rivo/tview"
)

type UI interface {
	handler.Hooks
	Run() error
	Stop()
}

type tui struct {
	// Directory to put files in during the session.
	// Usages must assume that this is a temporary directory, which is removed at the end of the program.
	sessionDir      string
	directory       string
	logFile         string
	timeline        *timeline
	app             *tview.Application
	header          *tview.TextView
	table           *tview.Table
	requestHeaders  *tview.Table
	requestBody     *FileView
	responseHeaders *tview.Table
	responseBody    *FileView
	footer          *tview.TextView
}

func (t *tui) Run() error {
	return t.app.Run()
}

func (t *tui) Stop() {
	t.app.Stop()
}

func NewTerminalUI(directory, sessionDir, logFile string) UI {
	tui := tui{
		directory:  directory,
		sessionDir: sessionDir,
		logFile:    logFile,
		timeline:   newTimeline(),
	}

	app := tview.NewApplication()
	tui.app = app

	header := tview.NewTextView()
	tui.header = header
	tui.setHeader()

	requestHeaders := tview.NewTable()
	requestHeaders.SetTitle(" (I)ncoming headers ")
	requestHeaders.SetBorder(true)
	tui.requestHeaders = requestHeaders

	requestBody := NewFileView()
	requestBody.SetTitle(" (i)ncoming body ")
	requestBody.SetBorder(true)
	tui.requestBody = requestBody

	responseHeaders := tview.NewTable()
	responseHeaders.SetTitle(" (O)utgoing headers ")
	responseHeaders.SetBorder(true)
	tui.responseHeaders = responseHeaders

	responseBody := NewFileView()
	responseBody.SetTitle(" (o)utgoing body ")
	responseBody.SetBorder(true)
	tui.responseBody = responseBody

	table := tview.NewTable().
		SetSelectable(true, false).
		SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGray)).
		SetSelectionChangedFunc(func(row, column int) {
			tui.setExchange()
		})
	table.SetTitle(" timeline ")
	table.Box.SetBorder(true)
	tui.table = table

	footer := tview.NewTextView().SetLabel("Commands: ").SetText("(y)ank request URL, (q)uit")
	tui.footer = footer

	tui.setupInputCaptures()

	grid := tview.NewGrid().
		SetRows(1, -2, -1, -3, 1).
		SetColumns(0, 0).
		AddItem(header, 0, 0, 1, 2, 0, 0, false).
		AddItem(table, 1, 0, 1, 2, 0, 0, true).
		AddItem(requestHeaders, 2, 0, 1, 1, 0, 0, false).
		AddItem(responseHeaders, 2, 1, 1, 1, 0, 0, false).
		AddItem(requestBody, 3, 0, 1, 1, 0, 0, false).
		AddItem(responseBody, 3, 1, 1, 1, 0, 0, false).
		AddItem(footer, 4, 0, 1, 2, 0, 0, false)
	app.SetRoot(grid, true).SetFocus(table)

	return &tui
}

func (t *tui) setExchange() {
	e := t.currentExchange()
	if e == nil {
		return
	}

	row := 0
	longestReqHeaderNameLen := len(e.Request.Headers.LongestName()) + 4
	t.requestHeaders.Clear()
	for _, h := range e.Request.Headers {
		for _, v := range h.Values {
			t.requestHeaders.SetCellSimple(row, 0, fmt.Sprintf("%-*s", longestReqHeaderNameLen, h.Key))
			t.requestHeaders.SetCellSimple(row, 1, v)
			row++
		}
	}
	t.requestHeaders.ScrollToBeginning()

	reqContent := e.Request.Get("Content-Type")
	t.requestBody.SetFile(t.fileOf("req", e.Id, reqContent), reqContent)

	if e.Response != nil {
		row = 0
		longestResHeaderNameLen := len(e.Response.Headers.LongestName()) + 4
		t.responseHeaders.Clear()
		for _, h := range e.Response.Headers {
			for _, v := range h.Values {
				t.responseHeaders.SetCellSimple(row, 0, fmt.Sprintf("%-*s", longestResHeaderNameLen, h.Key))
				t.responseHeaders.SetCellSimple(row, 1, v)
				row++
			}
		}
		t.responseHeaders.ScrollToBeginning()

		resContent := e.Response.Get("Content-Type")
		t.responseBody.SetFile(t.fileOf("res", e.Id, resContent), resContent)
	} else {
		t.responseBody.SetFile("", "")
	}
}

const ftime = "15:04:05.000000"

func (t *tui) setTable() {
	row := 0
	for _, id := range t.timeline.order {
		e := t.timeline.data[id]
		t.table.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf(" %d [%s]", id, e.Request.Timestamp.Format(ftime))).
			SetAlign(tview.AlignRight).
			SetTextColor(tcell.ColorGray),
		)

		t.table.SetCell(row, 1, tview.NewTableCell(e.Request.Method).
			SetAlign(tview.AlignRight).
			SetStyle(tcell.StyleDefault.Bold(true).Foreground(methodColor(e.Request.Method))),
		)
		t.table.SetCellSimple(row, 2, e.Request.Url)

		t.table.SetCell(row, 3, tview.NewTableCell(">>>").SetTextColor(tcell.ColorGray))

		if e.Response != nil {
			t.table.SetCell(row, 4, tview.NewTableCell(strconv.Itoa(e.Response.StatusCode)).
				SetAlign(tview.AlignCenter).
				SetStyle(tcell.StyleDefault.Bold(true).Foreground(statusColor(e.Response.StatusCode))),
			)
			t.table.SetCellSimple(row, 5, e.Response.StatusText)
			t.table.SetCell(row, 6, tview.NewTableCell(fmt.Sprintf("(%s)", e.Response.Timestamp.Sub(e.Request.Timestamp))).
				SetTextColor(tcell.ColorGray),
			)
			t.table.SetCellSimple(row, 7, e.Response.Get("Content-Type"))
		} else {
			t.table.SetCellSimple(row, 4, "")
			t.table.SetCellSimple(row, 5, "")
			t.table.SetCellSimple(row, 6, "")
			t.table.SetCellSimple(row, 7, "")
		}
		row++
	}
}

// etype is "req" or "res"
func (t *tui) fileOf(etype string, id uint64, contentType string) string {
	if contentType == "text/plain" {
		return ".txt"
	}
	exts, _ := mime.ExtensionsByType(contentType)
	var ext string
	if len(exts) > 0 {
		ext = exts[0]
	} else {
		ext = ".dat"
	}
	return filepath.Join(t.sessionDir, fmt.Sprintf("%d.%s%s", id, etype, ext))
}

func (t *tui) ReadFile(file string) (io.ReadCloser, string, error) {
	f := filepath.Join(t.directory, file)
	contentType := contentTypeOf(f)
	r, err := os.Open(f)
	return r, contentType, err
}

func (t *tui) onRequest(e *handler.Exchange, r *http.Request) (string, error) {
	defer r.Body.Close()
	hasBody := r.Body != http.NoBody
	file := t.fileOf("req", e.Id, r.Header.Get("Content-Type"))

	if hasBody {
		wf, err := os.Create(file)
		if err != nil {
			return "", err
		}
		_, err = io.Copy(wf, r.Body)
		if err != nil {
			wf.Close()
			return "", err
		}
		fs, err := wf.Stat()
		if err != nil {
			t.timeline.setReqBodySize(e.Id, -1)
		} else {
			t.timeline.setReqBodySize(e.Id, fs.Size())
		}
		wf.Close()
	}

	// Only add to timeline when body is written, so it can be displayed instantly
	go func() {
		t.timeline.add(e)
		t.app.QueueUpdateDraw(func() {
			t.setTable()
			t.setExchange()
			t.setHeader()
		})
	}()

	if !hasBody {
		return "", nil
	}
	return file, nil
}

func (t *tui) OnRequestWithoutFurtherBodyUsage(e *handler.Exchange, r *http.Request) error {
	_, err := t.onRequest(e, r)
	return err
}

func (t *tui) OnRequest(e *handler.Exchange, r *http.Request) (io.ReadCloser, error) {
	file, err := t.onRequest(e, r)
	if err != nil || file == "" {
		return http.NoBody, err
	}
	return os.Open(file)
}

func (t *tui) OnResponse(e *handler.Exchange, body io.ReadCloser, contentType string) (string, error) {
	defer body.Close()
	hasBody := body != http.NoBody
	file := t.fileOf("res", e.Id, contentType)

	if hasBody {
		wf, err := os.Create(file)
		if err != nil {
			return "", err
		}
		_, err = io.Copy(wf, body)
		if err != nil {
			wf.Close()
			return "", err
		}
		fs, err := wf.Stat()
		if err != nil {
			t.timeline.setResBodySize(e.Id, -1)
		} else {
			t.timeline.setResBodySize(e.Id, fs.Size())
		}
		wf.Close()
	}

	go func() {
		t.timeline.add(e)
		t.app.QueueUpdateDraw(func() {
			t.setTable()
			t.setExchange()
			t.setHeader()
		})
	}()

	if !hasBody {
		return "", nil
	}
	return file, nil
}

func (t *tui) Log(id uint64, msg string) {
	now := time.Now().Format(time.RFC3339Nano)
	f, _ := os.OpenFile(t.logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(fmt.Sprintf("[%s | %d] %s\n", now, id, msg))
}

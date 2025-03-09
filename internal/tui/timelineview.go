package tui

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const ftime = "15:04:05.000000"

var (
	knownColumns = []string{
		"id",
		"start",
		"method",
		"host",
		"request_target",
		"status_code",
		"status_text",
		"end",
		"duration",
	}

	methodDefaultStyle = tcell.StyleDefault.Bold(true)
	methodGetStyle     = methodDefaultStyle.Foreground(tcell.ColorGreen)
	methodPostStyle    = methodDefaultStyle.Foreground(tcell.ColorBlue)
	methodDeleteStyle  = methodDefaultStyle.Foreground(tcell.ColorRed)
	methodPutStyle     = methodDefaultStyle.Foreground(tcell.ColorYellow)
	methodPatchStyle   = methodDefaultStyle.Foreground(tcell.ColorTeal)
	methodHeadStyle    = methodDefaultStyle.Foreground(tcell.ColorPurple)

	statusDefaultStyle     = tcell.StyleDefault.Bold(true)
	statusSuccessStyle     = statusDefaultStyle.Foreground(tcell.ColorGreen)
	statusRedirectStyle    = statusDefaultStyle.Foreground(tcell.ColorLightCyan)
	statusClientErrorStyle = statusDefaultStyle.Foreground(tcell.ColorYellow)
	statusServerErrorStyle = statusDefaultStyle.Foreground(tcell.ColorRed)
)

type timelineData struct {
	tview.TableContentReadOnly
	columns []string
	entries []uint64
	lut     map[uint64]*timelineEntry
}

type timelineEntry struct {
	id            uint64
	start         time.Time
	fmtStart      string
	method        string
	host          string
	requestTarget string
	statusCode    int
	statusText    string
	fmtEnd        string
	fmtDuration   string
}

func pad(s string) string {
	return " " + s + " "
}

func unpad(s string) string {
	l := len(s)
	return s[1 : l-1]
}

func (d *timelineData) GetCell(row, column int) *tview.TableCell {
	if row == 0 {
		text := ""
		switch d.columns[column] {
		case "id":
			text = "ID"
		case "start":
			text = "Start"
		case "method":
			text = "Method"
		case "host":
			text = "Host"
		case "request_target":
			text = "Request Target"
		case "status_code":
			text = "Status Code"
		case "status_text":
			text = "Status Text"
		case "end":
			text = "End"
		case "duration":
			text = "Duration"
		}
		return tview.NewTableCell(pad(text)).SetStyle(headStyle).SetSelectable(false)
	}
	if row == 1 {
		return tview.NewTableCell("").SetSelectable(false)
	}
	row -= 2
	if row >= len(d.entries) {
		return nil
	}
	id := d.entries[row]
	e, ok := d.lut[id]
	if !ok {
		return nil
	}

	var style tcell.Style
	text := ""
	switch d.columns[column] {
	case "id":
		text = fmt.Sprintf("#%v", e.id)
	case "start":
		text = e.fmtStart
	case "method":
		text = e.method
		switch e.method {
		case "GET":
			style = methodGetStyle
		case "POST":
			style = methodPostStyle
		case "DELETE":
			style = methodDeleteStyle
		case "PUT":
			style = methodPutStyle
		case "PATCH":
			style = methodPatchStyle
		case "HEAD":
			style = methodHeadStyle
		default:
			style = methodDefaultStyle
		}
	case "host":
		text = e.host
	case "request_target":
		text = e.requestTarget
	case "status_code":
		if e.statusCode != 0 {
			text = strconv.Itoa(e.statusCode)
		}
		if e.statusCode >= 600 {
			style = statusDefaultStyle
		} else if e.statusCode >= 500 {
			style = statusServerErrorStyle
		} else if e.statusCode >= 400 {
			style = statusClientErrorStyle
		} else if e.statusCode >= 300 {
			style = statusRedirectStyle
		} else if e.statusCode >= 200 {
			style = statusSuccessStyle
		} else {
			style = statusDefaultStyle
		}
	case "status_text":
		text = e.statusText
	case "end":
		text = e.fmtEnd
	case "duration":
		text = e.fmtDuration
	}
	cell := tview.NewTableCell(pad(text)).SetSelectable(text != "")
	if style != tcell.StyleDefault {
		cell.SetStyle(style)
	}
	return cell
}

func (d *timelineData) GetRowCount() int {
	return len(d.entries) + 2
}

func (d *timelineData) GetColumnCount() int {
	return len(d.columns)
}

type TimelineView struct {
	*tview.Table
	data *timelineData
}

func toKnownColumns(columns []string) []string {
	res := make([]string, 0)
	for _, c := range columns {
		nc := strings.ToLower(strings.ReplaceAll(c, " ", "_"))
		if slices.Contains(knownColumns, nc) {
			res = append(res, nc)
		}
	}
	return res
}

// Known column names:
// - id
// - start
// - method
// - host
// - requestTarget
// - statusCode
// - statusText
// - end
// - duration
func NewTimelineView(columns []string) *TimelineView {
	table := tview.NewTable().SetFixed(2, 1).SetSelectable(true, true)
	table.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		tx, ty, twidth, theight := table.GetInnerRect()
		ysep := ty + 1
		for cx := tx; cx < tx+twidth; cx++ {
			screen.SetContent(cx, ysep, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault)
		}
		return tx, ty, twidth, theight
	})
	data := &timelineData{
		columns: toKnownColumns(columns),
		entries: make([]uint64, 0),
		lut:     make(map[uint64]*timelineEntry),
	}
	table.SetContent(data)
	return &TimelineView{
		Table: table,
		data:  data,
	}
}

// Known column names:
// - id
// - start
// - method
// - host
// - requestTarget
// - statusCode
// - statusText
// - end
// - duration
func (v *TimelineView) SetColumns(columns []string) {
	v.data.columns = toKnownColumns(columns)
}

func (v *TimelineView) AddRequest(
	id uint64,
	start time.Time,
	method string,
	url *url.URL,
) {
	requestTarget := url.Path
	if url.RawQuery != "" {
		requestTarget = requestTarget + "?" + url.RawQuery
	}
	entry := timelineEntry{
		id:            id,
		start:         start,
		fmtStart:      start.Format(ftime),
		method:        strings.ToUpper(method),
		host:          url.Host,
		requestTarget: requestTarget,
	}
	v.data.lut[id] = &entry
	v.data.entries = append(v.data.entries, id)
}

func (v *TimelineView) AddResponse(
	id uint64,
	end time.Time,
	statusCode int,
	statusText string,
) {
	e, ok := v.data.lut[id]
	if !ok {
		return
	}
	e.fmtEnd = end.Format(ftime)
	e.fmtDuration = end.Sub(e.start).String()
	e.statusCode = statusCode
	if statusText == "" {
		e.statusText = http.StatusText(statusCode)
	} else {
		e.statusText = statusText
	}
}

func (v *TimelineView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return v.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if v.Table.HasFocus() {
			switch event.Rune() {
			case 'y':
				sx, sy := v.Table.GetSelection()
				sc := v.Table.GetCell(sx, sy)
				if sc != nil {
					clipboard.WriteAll(unpad(sc.Text))
				}
			default:
				if handler := v.Table.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			}
		}
	})
}

func (v *TimelineView) Draw(screen tcell.Screen) {
	v.Table.Draw(screen)
}

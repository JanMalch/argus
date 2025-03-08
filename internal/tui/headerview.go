package tui

import (
	"encoding/base64"
	"slices"
	"strconv"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/gdamore/tcell/v2"
	"github.com/janmalch/argus/pkg/fmthttp"
	"github.com/pquerna/cachecontrol/cacheobject"
	"github.com/rivo/tview"
)

var httpDateHeaders = []string{
	"Last-Modified",
	"Expires",
	"Date",
	"If-Modified-Since",
	"If-Unmodified-Since",
	"Retry-After",
}

type Clock func() time.Time

type HeaderView struct {
	*tview.Table
	now Clock
}

func NewHeaderView() *HeaderView {
	return &HeaderView{
		Table: tview.NewTable().SetFixed(1, 3).SetSelectable(true, true).SetEvaluateAllRows(true),
		now: func() time.Time {
			return time.Now()
		},
	}
}

func analyze(key, value string, now time.Time) string {
	if key == "Cache-Control" {
		cc, err := cacheobject.ParseResponseCacheControl(value)
		if err == nil {
			maxAge := time.Duration(cc.MaxAge) * time.Second
			return strings.TrimSpace(humanize.RelTime(now.Add(maxAge.Abs()), now, "", ""))
		}
		return ""
	}

	if key == "Content-Length" {
		octets, err := strconv.ParseUint(value, 10, 64)
		if err == nil {
			return humanize.Bytes(octets)
		}
		return ""
	}

	if key == "Authorization" {
		if strings.HasPrefix(value, "Basic ") {
			data, err := base64.StdEncoding.DecodeString(value[6:])
			if err == nil {
				return string(data)
			}
		} else if strings.HasPrefix(value, "Bearer ey") {
			return stringifyJwt(value[7:])
		}
		return ""
	}

	if slices.Contains(httpDateHeaders, key) {
		parsedTime, err := time.Parse(time.RFC1123, value)
		if err == nil {
			return humanize.RelTime(parsedTime, now, "ago", "from now")
		}
		return ""
	}

	return ""
}

func stringifyJwt(tokenString string) string {
	parts := strings.Split(tokenString, ".")
	var sb strings.Builder

	header, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return ""
	}
	sb.WriteString(string(header))

	if len(parts) > 1 {
		payload, err := base64.RawStdEncoding.DecodeString(parts[1])
		if err != nil {
			return ""
		}
		sb.WriteRune('.')
		sb.WriteString(string(payload))
	}

	return sb.String()
}

func (v *HeaderView) SetHeaders(h fmthttp.Headers) {
	v.Table.Clear()
	v.Table.SetCell(0, 0, tview.NewTableCell("Key").SetSelectable(false).SetStyle(headStyle))
	v.Table.SetCell(0, 1, tview.NewTableCell("Value").SetSelectable(false).SetStyle(headStyle))
	v.Table.SetCell(0, 2, tview.NewTableCell("Analyzed").SetSelectable(false).SetStyle(headStyle))

	row := 1
	for _, header := range h {
		for _, hv := range header.Values {
			v.Table.SetCellSimple(row, 0, header.Key)
			// FIXME: find a better way than hardcoding max width ... maybe hijack Draw() for measuring?
			v.Table.SetCell(row, 1, tview.NewTableCell(hv).SetMaxWidth(40))
			v.Table.SetCellSimple(row, 2, analyze(header.Key, hv, v.now()))
			row++
		}
	}
}

func (v *HeaderView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return v.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if v.Table.HasFocus() {
			switch event.Rune() {
			case 'y':
				yankSelectedCellText(v.Table)
			default:
				if handler := v.Table.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			}
		}
	})
}

func (v *HeaderView) Draw(screen tcell.Screen) {
	v.Table.Draw(screen)
}

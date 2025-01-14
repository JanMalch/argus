package ui

import (
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (t *tui) consumeGlobalEvents(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'q':
		t.app.Stop()
		return nil
	case 'I':
		t.app.SetFocus(t.requestHeaders)
		t.footer.SetText("(y)ank headers, (q)uit, (Esc)ape to timeline")
		return nil
	case 'i':
		t.app.SetFocus(t.requestBody)
		t.footer.SetText("(y)ank body, (q)uit, (Esc)ape to timeline")
		return nil
	case 'O':
		t.app.SetFocus(t.responseHeaders)
		t.footer.SetText("(y)ank headers, (q)uit, (Esc)ape to timeline")
		return nil
	case 'o':
		t.app.SetFocus(t.responseBody)
		t.footer.SetText("(y)ank body, (q)uit, (Esc)ape to timeline")
		return nil
	}
	return event
}

func (t *tui) setupTableInputCapture(yankUrl func()) {
	t.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'y':
			yankUrl()
			return nil
		}
		return t.consumeGlobalEvents(event)
	})
}

func (t *tui) setupHeadBodyInputCaptures(
	headers *tview.Table,
	body *FileView,
	yankHeaders func(),
	yankBody func(),
) {
	headers.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			t.app.SetFocus(t.table)
			return nil
		}
		switch event.Rune() {
		case 'y':
			yankHeaders()
			t.app.SetFocus(t.table)
			return nil
		}
		return t.consumeGlobalEvents(event)
	})
	body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			t.app.SetFocus(t.table)
			return nil
		}
		switch event.Rune() {
		case 'y':
			yankBody()
			t.app.SetFocus(t.table)
			return nil
		}
		return t.consumeGlobalEvents(event)
	})
}

func (t *tui) setupInputCaptures() {
	t.setupTableInputCapture(func() {
		if e := t.currentExchange(); e != nil {
			clipboard.WriteAll(e.Request.Url)
		}
	})
	t.setupHeadBodyInputCaptures(t.requestHeaders, t.requestBody, func() {
		if e := t.currentExchange(); e != nil {
			clipboard.WriteAll(e.Request.Headers.String())
		}
	}, func() {
		if e := t.currentExchange(); e != nil {
			txt := t.currentRequestBodyString()
			if txt != "" {
				clipboard.WriteAll(txt)
			}
		}
	})
	t.setupHeadBodyInputCaptures(t.responseHeaders, t.responseBody, func() {
		if e := t.currentExchange(); e != nil {
			clipboard.WriteAll(e.Response.Headers.String())
		}
	}, func() {
		if e := t.currentExchange(); e != nil {
			txt := t.currentResponseBodyString()
			if txt != "" {
				clipboard.WriteAll(txt)
			}
		}
	})

}

package ui

import (
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type state struct {
	isInYankMode bool
}

func (t *tui) consumeGlobalEvents(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'q':
		t.app.Stop()
		return nil
	case 'I':
		t.app.SetFocus(t.requestHeaders)
		return nil
	case 'i':
		t.app.SetFocus(t.requestBody)
		return nil
	case 'O':
		t.app.SetFocus(t.responseHeaders)
		return nil
	case 'o':
		t.app.SetFocus(t.responseBody)
		return nil
	}
	return event
}

func (t *tui) setupTableInputCapture(s *state, yankUrl func()) {
	t.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			s.isInYankMode = false
			return nil
		}
		switch event.Rune() {
		case 'y':
			if !s.isInYankMode {
				s.isInYankMode = true
				return nil
			}
			fallthrough
		case 'Y':
			yankUrl()
			s.isInYankMode = false
			return nil
		}
		return t.consumeGlobalEvents(event)
	})
}

func (t *tui) setupHeadBodyInputCaptures(
	s *state,
	headers *tview.Table,
	body *FileView,
	yankHeaders func(),
	yankBody func(),
) {
	headers.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			s.isInYankMode = false
			t.app.SetFocus(t.table)
			return nil
		}
		switch event.Rune() {
		case 'y':
			if !s.isInYankMode {
				s.isInYankMode = true
				return nil
			}
			fallthrough
		case 'Y':
			yankHeaders()
			s.isInYankMode = false
			t.app.SetFocus(t.table)
			return nil
		}
		return t.consumeGlobalEvents(event)
	})
	body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			s.isInYankMode = false
			t.app.SetFocus(t.table)
			return nil
		}
		switch event.Rune() {
		case 'y':
			if !s.isInYankMode {
				s.isInYankMode = true
				return nil
			}
			fallthrough
		case 'Y':
			yankBody()
			s.isInYankMode = false
			t.app.SetFocus(t.table)
			return nil
		}
		return t.consumeGlobalEvents(event)
	})
}

func (t *tui) setupInputCaptures() {
	s := state{isInYankMode: false}
	t.setupTableInputCapture(&s, func() {
		if e := t.currentExchange(); e != nil {
			clipboard.WriteAll(e.Request.Url)
		}
	})
	t.setupHeadBodyInputCaptures(&s, t.requestHeaders, t.requestBody, func() {
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
	t.setupHeadBodyInputCaptures(&s, t.responseHeaders, t.responseBody, func() {
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

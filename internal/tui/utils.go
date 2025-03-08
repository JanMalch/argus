package tui

import (
	"github.com/atotto/clipboard"
	"github.com/rivo/tview"
)

func yankSelectedCellText(t *tview.Table) {
	sx, sy := t.GetSelection()
	sc := t.GetCell(sx, sy)
	if sc != nil {
		clipboard.WriteAll(sc.Text)
	}
}

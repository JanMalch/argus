package tui

import (
	"net/url"

	"github.com/gdamore/tcell/v2"
	"github.com/janmalch/argus/pkg/fmthttp"
	"github.com/rivo/tview"
)

var (
	headStyle = tcell.StyleDefault.Bold(true)
)

type ParameterView struct {
	*tview.Table
	parameters fmthttp.Parameters
	decode     int
}

func NewParameterView() *ParameterView {
	return &ParameterView{
		Table: tview.NewTable().SetFixed(1, 2).SetSelectable(true, true).SetEvaluateAllRows(true),
	}
}

func (v *ParameterView) update() {
	v.Table.Clear()
	v.Table.SetCell(0, 0, tview.NewTableCell("Key").SetSelectable(false).SetStyle(headStyle))
	v.Table.SetCell(0, 1, tview.NewTableCell("Value").SetSelectable(false).SetStyle(headStyle))

	row := 1
	for _, param := range v.parameters {
		for _, pv := range param.Values {
			v.Table.SetCellSimple(row, 0, param.Key)
			if v.decode == 0 {
				v.Table.SetCellSimple(row, 1, pv)
			} else {
				rem := v.decode
				dpv := pv
				var err error
				for rem > 0 {
					dpv, err = url.QueryUnescape(dpv)
					if err != nil {
						dpv = pv
						break
					}
					rem--
				}
				v.Table.SetCellSimple(row, 1, dpv)
			}

			row++
		}
	}
}

func (v *ParameterView) SetParameters(p fmthttp.Parameters) {
	v.parameters = p
	v.update()
}

func (v *ParameterView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return v.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if v.Table.HasFocus() {
			switch event.Rune() {
			case 'D':
				v.decode = max(0, v.decode-1)
				v.update()
			case 'd':
				v.decode++
				v.update()
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

func (v *ParameterView) Draw(screen tcell.Screen) {
	v.Table.Draw(screen)
}

package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/janmalch/argus/pkg/fmthttp"
	"github.com/rivo/tview"
)

var (
	tabDeselected = tcell.StyleDefault.Bold(false).Foreground(tcell.ColorGray)
	tabSelected   = tcell.StyleDefault.Bold(true).Foreground(tcell.ColorDefault)
)

type ExchangeView struct {
	*tview.Grid
	tab int

	parameterTitle *tview.TextView
	reqHeaderTitle *tview.TextView
	reqBodyTitle   *tview.TextView
	resHeaderTitle *tview.TextView
	resBodyTitle   *tview.TextView

	pages         *tview.Pages
	parameterView *ParameterView
	reqHeaderView *HeaderView
	reqBodyView   *FileView
	resHeaderView *HeaderView
	resBodyView   *FileView
}

func NewExchangeView() *ExchangeView {
	newTab := func(text string) *tview.TextView {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text).
			SetTextStyle(tabDeselected)
	}

	parameterTitle := newTab("Parameters")
	reqHeaderTitle := newTab("Req. Headers")
	reqBodyTitle := newTab("Req. Body")
	resHeaderTitle := newTab("Res. Headers")
	resBodyTitle := newTab("Res. Body")
	resBodyTitle.SetTextStyle(tabSelected)

	pages := tview.NewPages()
	parameterView := NewParameterView()
	parameterView.SetBorderPadding(0, 0, 1, 1)
	reqHeaderView := NewHeaderView()
	reqHeaderView.SetBorderPadding(0, 0, 1, 1)
	reqBodyView := NewFileView()
	reqBodyView.SetBorderPadding(0, 0, 1, 1)
	resHeaderView := NewHeaderView()
	resHeaderView.SetBorderPadding(0, 0, 1, 1)
	resBodyView := NewFileView()
	resBodyView.SetBorderPadding(0, 0, 1, 1)

	pages.AddPage("p0", parameterView, true, false)
	pages.AddPage("p1", reqHeaderView, true, false)
	pages.AddPage("p2", reqBodyView, true, false)
	pages.AddPage("p3", resHeaderView, true, false)
	pages.AddPage("p4", resBodyView, true, true)
	grid := tview.NewGrid().
		SetRows(1, 0).
		SetColumns(0, 0, 0, 0, 0).
		SetBorders(true).
		AddItem(parameterTitle, 0, 0, 1, 1, 0, 0, false).
		AddItem(reqHeaderTitle, 0, 1, 1, 1, 0, 0, false).
		AddItem(reqBodyTitle, 0, 2, 1, 1, 0, 0, false).
		AddItem(resHeaderTitle, 0, 3, 1, 1, 0, 0, false).
		AddItem(resBodyTitle, 0, 4, 1, 1, 0, 0, false).
		AddItem(pages, 1, 0, 1, 5, 0, 0, true)
	return &ExchangeView{
		tab:            4,
		Grid:           grid,
		parameterTitle: parameterTitle,
		reqHeaderTitle: reqHeaderTitle,
		reqBodyTitle:   reqBodyTitle,
		resHeaderTitle: resHeaderTitle,
		resBodyTitle:   resBodyTitle,
		pages:          pages,
		parameterView:  parameterView,
		reqHeaderView:  reqHeaderView,
		reqBodyView:    reqBodyView,
		resHeaderView:  resHeaderView,
		resBodyView:    resBodyView,
	}
}

func (v *ExchangeView) getActiveView() tview.Primitive {
	_, content := v.pages.GetFrontPage()
	if content == nil {
		panic("failed to determine content view for exchange")
	}
	return content
}

func (v *ExchangeView) onTabChange() tview.Primitive {
	v.parameterTitle.SetTextStyle(tabDeselected)
	v.reqHeaderTitle.SetTextStyle(tabDeselected)
	v.reqBodyTitle.SetTextStyle(tabDeselected)
	v.resHeaderTitle.SetTextStyle(tabDeselected)
	v.resBodyTitle.SetTextStyle(tabDeselected)
	switch v.tab {
	case 0:
		v.parameterTitle.SetTextStyle(tabSelected)
	case 1:
		v.reqHeaderTitle.SetTextStyle(tabSelected)
	case 2:
		v.reqBodyTitle.SetTextStyle(tabSelected)
	case 3:
		v.resHeaderTitle.SetTextStyle(tabSelected)
	case 4:
		v.resBodyTitle.SetTextStyle(tabSelected)
	}
	v.pages.SwitchToPage(fmt.Sprintf("p%d", v.tab))
	return v.getActiveView()
}

func (v *ExchangeView) SetRequest(parameters fmthttp.Parameters, headers fmthttp.Headers, bodyFile string) *ExchangeView {
	v.parameterView.SetParameters(parameters)
	v.reqHeaderView.SetHeaders(headers)
	v.reqBodyView.SetFile(bodyFile)
	return v
}

func (v *ExchangeView) SetResponse(headers fmthttp.Headers, bodyFile string) *ExchangeView {
	v.resHeaderView.SetHeaders(headers)
	v.resBodyView.SetFile(bodyFile)
	return v
}

func (v *ExchangeView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return v.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyTab:
			v.tab = (v.tab + 1) % 5
			setFocus(v.onTabChange())
			return
		case tcell.KeyBacktab:
			if v.tab == 0 {
				v.tab = 4
			} else {
				v.tab--
			}
			setFocus(v.onTabChange())
			return
		}
		content := v.getActiveView()
		if handler := content.InputHandler(); handler != nil {
			handler(event, setFocus)
			return
		}
	})
}

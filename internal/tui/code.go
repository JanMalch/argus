package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/lipgloss"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	lineNumberStyle = lipgloss.NewStyle().Inline(true).Foreground(lipgloss.AdaptiveColor{Light: "249", Dark: "8"})
)

type CodeView struct {
	*tview.TextView
}

func NewCodeView() *CodeView {
	return &CodeView{
		TextView: tview.NewTextView().SetDynamicColors(true),
	}
}
func (v *CodeView) SetText(text, contentType string) {
	v.TextView.ScrollToBeginning()
	content := text
	extension := ""
	if contentType[0] == '.' {
		extension = contentType
	} else {
		extension = extensionByType(contentType, "."+contentType)
	}

	content = prettier(content, extension)

	var sb strings.Builder
	err := quick.Highlight(&sb, content, extension, "terminal256", "dracula")
	if err != nil {
		v.TextView.SetText(content)
		return
	}
	writer := tview.ANSIWriter(v.TextView)
	writer.Write([]byte(applyLineNumbers(sb.String())))
}

func (v *CodeView) Draw(screen tcell.Screen) {
	v.TextView.Draw(screen)
}

// Prettifies the given content, based on the given extension.
// Returns the original content in case of internal errors.
func prettier(content, extension string) string {
	switch extension {
	case ".json":
		var prettyJson bytes.Buffer
		err := json.Indent(&prettyJson, []byte(content), "", "    ")
		if err != nil {
			return content
		}
		return prettyJson.String()
	}
	return content
}

// Returns the digits in i.
//
//	digits(100) == 3
func digits(i int) int {
	if i == 0 {
		return 1
	}
	count := 0
	for i != 0 {
		i /= 10
		count++
	}
	return count
}

// Returns a func, which will add linenumbers to the received string,
// if enabled.
func applyLineNumbers(s string) string {
	var sb strings.Builder
	lines := strings.Split(s, "\n")
	maxLineNumberWidth := digits(len(lines))
	for i, line := range lines {
		sb.WriteString(lineNumberStyle.Render(fmt.Sprintf("%*d ", maxLineNumberWidth, i+1)))
		sb.WriteString(line)
		if i < len(lines) {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

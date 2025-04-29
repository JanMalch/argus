package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	keyStyle = lipgloss.NewStyle().Inline(true).Foreground(lipgloss.AdaptiveColor{Light: "249", Dark: "8"})
)

func printKv(indent int, key, value string) {
	line := strings.Repeat(" ", indent) + keyStyle.Render(key+": ") + value
	fmt.Println(line)
}

func DebugPrint(c Config, path string) {
	printKv(0, "Config location", path)
	fmt.Println(keyStyle.Render(strings.Repeat("-", len(path)+17)))
	printKv(0, "Directory", c.Directory)
	fmt.Println()

	printKv(0, fmt.Sprintf("%d Servers", len(c.Servers)), "")
	for i, s := range c.Servers {
		fmt.Printf("  Server #%d:\n", i+1)
		printKv(4, "Upstream", s.Upstream.String())
		printKv(4, "Port", strconv.Itoa(s.Port))

		printKv(4, "Requests header overwrites", strconv.Itoa(len(s.Request.Headers)))
		for k, v := range s.Request.Headers {
			fmt.Printf("      %s: %s\n", k, v)
		}

		printKv(4, "Requests query parameter overwrites", strconv.Itoa(len(s.Request.Parameters)))
		for k, v := range s.Request.Parameters {
			fmt.Printf("      %s: %s\n", k, v)
		}

		printKv(4, "Response header overwrites", strconv.Itoa(len(s.Response.Headers)))
		for k, v := range s.Response.Headers {
			fmt.Printf("      %s: %s\n", k, v)
		}
		printKv(4, "Response overwrites", strconv.Itoa(len(s.Response.Overwrites)))
		for _, o := range s.Response.Overwrites {
			if o.Regex != nil {
				printKv(6, "Matches any method on paths with regex", o.Regex.String())
			} else if o.Method != "" {
				printKv(6, fmt.Sprintf("Matches Matches %s method on exactly", o.Method), o.Exact)
			} else {
				printKv(6, "Matches any method on exactly", o.Exact)
			}

			if o.File != "" {
				_, err := os.Stat(o.File)
				if err == nil {
					printKv(8, "Responds with file", o.File+" "+keyStyle.Render("(exists)"))
				} else if errors.Is(err, os.ErrNotExist) {
					printKv(8, "Responds with file", o.File+" "+keyStyle.Render("(does not exist)"))
				} else {
					log.Panic(err)
				}
			}
			if o.Status > 0 {
				printKv(8, "Responds with HTTP status", strconv.Itoa(o.Status))
			}
		}
		fmt.Println()
	}

	printKv(0, "UI", "")
	layout := "vertical"
	if c.UI.Horizontal {
		layout = "horizontal"
	}
	printKv(2, "Layout", layout)
	printKv(2, "Timeline grow factor", strconv.Itoa(c.UI.GrowTimeline))
	printKv(2, "Exchange grow factor", strconv.Itoa(c.UI.GrowExchange))
	printKv(2, fmt.Sprintf("%d timline columns", len(c.UI.TimelineColumns)), fmt.Sprintf("%s", c.UI.TimelineColumns))
}

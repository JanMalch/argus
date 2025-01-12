package ui

import "github.com/gdamore/tcell/v2"

func methodColor(method string) tcell.Color {
	switch method {
	case "GET":
		return tcell.ColorGreen
	case "POST":
		return tcell.ColorBlue
	case "DELETE":
		return tcell.ColorRed
	case "PUT":
		return tcell.ColorYellow
	case "PATCH":
		return tcell.ColorTeal
	case "HEAD":
		return tcell.ColorPurple
	default:
		return tcell.ColorDefault
	}
}

func statusColor(status int) tcell.Color {
	if status >= 600 {
		return tcell.ColorDefault
	} else if status >= 500 {
		return tcell.ColorRed
	} else if status >= 400 {
		return tcell.ColorYellow
	} else if status >= 300 {
		return tcell.ColorLightCyan
	} else if status >= 200 {
		return tcell.ColorGreen
	} else {
		return tcell.ColorDefault
	}
}

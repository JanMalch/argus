package ui

import (
	"fmt"
	"runtime"
)

func (t *tui) setHeader() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	mbUsage := float64(memStats.Alloc/1024) / 1024

	t.header.SetText(fmt.Sprintf("ARGUS │ exchanges: %d │ memory: %.2f MB", t.timeline.Len(), mbUsage))
}

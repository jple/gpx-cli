package tui

import (
	"strings"

	sym "github.com/jple/text-symbol"
)

type StringWriter struct {
	cursor  int // TUI cursor position
	k       int // writing (line) count
	builder *strings.Builder
}

func (w *StringWriter) Write(p []byte) (int, error) {
	cursorAtWritingLine := w.cursor == w.k
	w.k++ // update current writing line
	// NOTE: it supposes Write is call at each line

	if cursorAtWritingLine {
		q := sym.BgBrightGreen(string(p))
		return w.builder.WriteString(q)
	}
	return w.builder.Write(p)
}

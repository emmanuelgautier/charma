// Package output provides writers for the four supported output formats.
package output

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// IsTerminal reports whether f is connected to a terminal.
func IsTerminal(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// WriteTerminal writes styled lines to w.
// If noColor is true, or w is not a terminal, ANSI escape codes are stripped first.
func WriteTerminal(w io.Writer, lines []string, noColor bool) error {
	// Detect if w is a real terminal when w is *os.File.
	if f, ok := w.(*os.File); ok && !IsTerminal(f) {
		noColor = true
	}
	for _, line := range lines {
		text := line
		if noColor {
			text = stripANSI(text)
		}
		if _, err := fmt.Fprintln(w, text); err != nil {
			return err
		}
	}
	return nil
}

func stripANSI(s string) string {
	var b strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			i += 2
			for i < len(s) && (s[i] < 0x40 || s[i] > 0x7e) {
				i++
			}
			if i < len(s) {
				i++
			}
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}

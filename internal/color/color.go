// Package color applies ANSI colors and gradients to text lines using lipgloss.
package color

import (
	"fmt"
	stdcolor "image/color"
	"math"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
)

// AvailableColors lists all accepted color names.
var AvailableColors = []string{
	"default", "red", "green", "blue", "cyan", "magenta", "yellow", "white",
}

// IsValidColor reports whether name is a supported color name or a valid hex color.
func IsValidColor(name string) bool {
	for _, c := range AvailableColors {
		if c == name {
			return true
		}
	}
	return IsHexColor(name)
}

// IsHexColor reports whether s is a valid CSS hex color (#RGB or #RRGGBB).
func IsHexColor(s string) bool {
	if len(s) == 0 || s[0] != '#' {
		return false
	}
	hex := s[1:]
	if len(hex) != 3 && len(hex) != 6 {
		return false
	}
	for _, c := range hex {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false
		}
	}
	return true
}

// hexColors maps color names to hex strings for gradient interpolation.
var hexColors = map[string]string{
	"red":     "#FF0000",
	"green":   "#00CC00",
	"blue":    "#0066FF",
	"cyan":    "#00CCFF",
	"magenta": "#FF00FF",
	"yellow":  "#FFFF00",
	"white":   "#FFFFFF",
}

// ansiColors maps color names to ANSI bright color codes.
var ansiColors = map[string]string{
	"red":     "9",
	"green":   "10",
	"blue":    "12",
	"cyan":    "14",
	"magenta": "13",
	"yellow":  "11",
	"white":   "15",
}

// ApplyColor applies a uniform color to every line.
// If colorName is "default" or empty, lines are returned unchanged.
// colorName may be a named color or a hex value like "#FF0000" or "#F00".
func ApplyColor(lines []string, colorName string) []string {
	if colorName == "" || colorName == "default" {
		return lines
	}
	var lipglossColor stdcolor.Color
	if IsHexColor(colorName) {
		lipglossColor = lipgloss.Color(expandHex(colorName))
	} else {
		code, ok := ansiColors[colorName]
		if !ok {
			return lines
		}
		lipglossColor = lipgloss.Color(code)
	}
	style := lipgloss.NewStyle().Foreground(lipglossColor)
	result := make([]string, len(lines))
	for i, l := range lines {
		result[i] = style.Render(l)
	}
	return result
}

// resolveHex returns the hex string for a color name or hex value.
// For named colors it looks up hexColors; for hex values it expands if needed.
func resolveHex(name string) (string, bool) {
	if IsHexColor(name) {
		return expandHex(name), true
	}
	h, ok := hexColors[name]
	return h, ok
}

// ApplyGradient applies a two-color gradient across lines.
// gradient must be in "color1:color2" format (e.g. "red:blue" or "#FF0000:#0000FF").
func ApplyGradient(lines []string, gradient string) ([]string, error) {
	parts := strings.SplitN(gradient, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("gradient must be in 'color1:color2' format, got %q", gradient)
	}
	c1Name, c2Name := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	hex1, ok1 := resolveHex(c1Name)
	hex2, ok2 := resolveHex(c2Name)
	if !ok1 {
		return nil, fmt.Errorf("unsupported gradient color %q", c1Name)
	}
	if !ok2 {
		return nil, fmt.Errorf("unsupported gradient color %q", c2Name)
	}

	r1, g1, b1 := parseHex(hex1)
	r2, g2, b2 := parseHex(hex2)

	result := make([]string, len(lines))
	n := len(lines)
	for i, line := range lines {
		var t float64
		if n > 1 {
			t = float64(i) / float64(n-1)
		}
		r := lerp(r1, r2, t)
		g := lerp(g1, g2, t)
		b := lerp(b1, b2, t)
		hex := fmt.Sprintf("#%02X%02X%02X", r, g, b)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(hex))
		result[i] = style.Render(line)
	}
	return result, nil
}

// StripANSI removes ANSI CSI escape sequences from s.
func StripANSI(s string) string {
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

// StripANSILines calls StripANSI on each line.
func StripANSILines(lines []string) []string {
	result := make([]string, len(lines))
	for i, l := range lines {
		result[i] = StripANSI(l)
	}
	return result
}

// expandHex normalises a hex color to full #RRGGBB form.
// #RGB is expanded to #RRGGBB; #RRGGBB is returned unchanged.
func expandHex(s string) string {
	h := strings.TrimPrefix(s, "#")
	if len(h) == 3 {
		h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
	}
	return "#" + strings.ToUpper(h)
}

func parseHex(hex string) (uint8, uint8, uint8) {
	hex = strings.TrimPrefix(hex, "#")
	v, _ := strconv.ParseUint(hex, 16, 32)
	return uint8(v >> 16), uint8((v >> 8) & 0xFF), uint8(v & 0xFF)
}

func lerp(a, b uint8, t float64) uint8 {
	return uint8(math.Round(float64(a) + (float64(b)-float64(a))*t))
}

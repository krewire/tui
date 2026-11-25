package term

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// RuneWidth returns the terminal visual column width of a single rune:
// 0 for control characters, combining marks, format/zero-width characters;
// 2 for East Asian Wide / Fullwidth characters and emoji;
// 1 for ordinary halfwidth characters.
func RuneWidth(r rune) int {
	if r == 0 || r < 32 || (r >= 0x7F && r <= 0x9F) {
		return 0
	}
	// Zero-width characters & combining marks
	if unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || unicode.Is(unicode.Cf, r) {
		return 0
	}
	// Specific zero-width spaces and joiners
	if r == 0x200B || r == 0x200C || r == 0x200D || r == 0xFEFF {
		return 0
	}
	// Wide characters: East Asian Fullwidth/Wide, CJK Ideographs, Hangul, Emojis
	if isWide(r) {
		return 2
	}
	return 1
}

func isWide(r rune) bool {
	switch {
	case r >= 0x1100 && r <= 0x115F: // Hangul Jamo
		return true
	case r >= 0x2329 && r <= 0x232A:
		return true
	case r >= 0x2E80 && r <= 0xA4CF && r != 0x303F: // CJK Radicals, Kangxi, Ideographs, Yi
		return true
	case r >= 0xAC00 && r <= 0xD7A3: // Hangul Syllables
		return true
	case r >= 0xF900 && r <= 0xFAFF: // CJK Compatibility Ideographs
		return true
	case r >= 0xFE10 && r <= 0xFE19: // Vertical forms
		return true
	case r >= 0xFE30 && r <= 0xFE6F: // CJK Compatibility Forms
		return true
	case r >= 0xFF00 && r <= 0xFF60: // Fullwidth Forms
		return true
	case r >= 0xFFE0 && r <= 0xFFE6: // Fullwidth Symbols
		return true
	case r >= 0x1F300 && r <= 0x1FAFF: // Emojis (Pictographs, Transport, Alchemical, Supplemental)
		return true
	case r >= 0x20000 && r <= 0x2FA1F: // CJK Extension B, C, D, E, F, I
		return true
	}
	return false
}

// StringWidth returns the visual column width of s, excluding ANSI escape sequences.
func StringWidth(s string) int {
	clean := Sanitize(s)
	width := 0
	for _, r := range clean {
		if r == '\t' {
			width += 4 - (width % 4)
			continue
		}
		width += RuneWidth(r)
	}
	return width
}

// Wrap wraps text to a maximum visual width in terminal columns.
// Explicit newlines in text are preserved. Words longer than width are broken
// across multiple lines. If width <= 0, text is returned unmodified.
func Wrap(text string, width int) string {
	if width <= 0 || len(text) == 0 {
		return text
	}

	lines := strings.Split(text, "\n")
	var out []string

	for _, line := range lines {
		if StringWidth(line) <= width {
			out = append(out, line)
			continue
		}

		words := strings.Fields(line)
		if len(words) == 0 {
			out = append(out, "")
			continue
		}

		var currentLine strings.Builder
		currentWidth := 0

		for _, word := range words {
			wWidth := StringWidth(word)

			// If a single word exceeds width, break it down
			if wWidth > width {
				if currentWidth > 0 {
					out = append(out, currentLine.String())
					currentLine.Reset()
					currentWidth = 0
				}
				// Break the word rune by rune
				for len(word) > 0 {
					r, size := utf8.DecodeRuneInString(word)
					rw := RuneWidth(r)
					if currentWidth+rw > width && currentWidth > 0 {
						out = append(out, currentLine.String())
						currentLine.Reset()
						currentWidth = 0
					}
					currentLine.WriteRune(r)
					currentWidth += rw
					word = word[size:]
				}
				continue
			}

			if currentWidth == 0 {
				currentLine.WriteString(word)
				currentWidth = wWidth
			} else if currentWidth+1+wWidth <= width {
				currentLine.WriteByte(' ')
				currentLine.WriteString(word)
				currentWidth += 1 + wWidth
			} else {
				out = append(out, currentLine.String())
				currentLine.Reset()
				currentLine.WriteString(word)
				currentWidth = wWidth
			}
		}

		if currentWidth > 0 {
			out = append(out, currentLine.String())
		}
	}

	return strings.Join(out, "\n")
}

// Truncate truncates text to fit within width columns, appending tail if truncated.
// If width <= 0, it returns an empty string. If StringWidth(text) <= width,
// text is returned unmodified.
func Truncate(text string, width int, tail string) string {
	if width <= 0 {
		return ""
	}
	if StringWidth(text) <= width {
		return text
	}

	tailWidth := StringWidth(tail)
	if tailWidth >= width {
		// If tail itself doesn't fit, truncate tail
		var sb strings.Builder
		w := 0
		for _, r := range tail {
			rw := RuneWidth(r)
			if w+rw > width {
				break
			}
			sb.WriteRune(r)
			w += rw
		}
		return sb.String()
	}

	target := width - tailWidth
	var sb strings.Builder
	current := 0

	for _, r := range text {
		rw := RuneWidth(r)
		if current+rw > target {
			break
		}
		sb.WriteRune(r)
		current += rw
	}
	sb.WriteString(tail)
	return sb.String()
}

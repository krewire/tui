// Package term provides terminal I/O, output formatting, and color
// conventions for the Krewire ecosystem.
package term

import (
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Color is an ANSI 16-color palette value.
type Color int

const (
	// ColorBlack is ANSI black.
	ColorBlack Color = 30
	// ColorRed is ANSI red.
	ColorRed Color = 31
	// ColorGreen is ANSI green.
	ColorGreen Color = 32
	// ColorYellow is ANSI yellow.
	ColorYellow Color = 33
	// ColorBlue is ANSI blue.
	ColorBlue Color = 34
	// ColorMagenta is ANSI magenta.
	ColorMagenta Color = 35
	// ColorCyan is ANSI cyan.
	ColorCyan Color = 36
	// ColorWhite is ANSI white.
	ColorWhite Color = 37
	// ColorDefault is the terminal default foreground.
	ColorDefault Color = 39
)

// Style is a text rendering style.
type Style int

const (
	// StyleBold is bold weight.
	StyleBold Style = 1
	// StyleDim is dim intensity.
	StyleDim Style = 2
	// StyleUnderline is underlined text.
	StyleUnderline Style = 4
)

// Paint wraps text in ANSI escape sequences for the given color and styles.
// text is sanitized by default to strip control and escape sequences.
// Use RawPaint for trusted strings that legitimately contain ANSI sequences.
// Use Terminal.Paint instead when the output may not support colors.
func Paint(text string, color Color, styles []Style) string {
	return RawPaint(Sanitize(text), color, styles)
}

// RawPaint wraps text in ANSI escape sequences for color and styles without
// sanitizing the text. It is an explicit, trusted-only API intended for
// caller-controlled ANSI art or pre-formatted terminal sequences.
func RawPaint(text string, color Color, styles []Style) string {
	codes := make([]string, 0, len(styles)+1)
	for _, s := range styles {
		codes = append(codes, strconv.Itoa(int(s)))
	}
	codes = append(codes, strconv.Itoa(int(color)))
	return "\x1b[" + strings.Join(codes, ";") + "m" + text + "\x1b[0m"
}

// Sanitize removes ANSI/VT100 escape sequences (CSI, OSC, DCS, APC, PM) and
// dangerous C0/C1 control characters from text, preventing terminal injection
// attacks while preserving intentional newlines (\n), tabs (\t), and valid Unicode.
//
// Standalone carriage returns (\r) are stripped to prevent cursor-rewind attacks;
// CRLF pairs (\r\n) are normalized to newlines (\n).
func Sanitize(text string) string {
	if text == "" {
		return ""
	}

	// Fast path: if there are no escape characters, C0 controls (except \n, \t),
	// or C1 control bytes, return text without allocation.
	needsSanitize := false
	for i := 0; i < len(text); i++ {
		b := text[i]
		if b < 0x20 {
			if b != '\n' && b != '\t' {
				needsSanitize = true
				break
			}
		} else if b == 0x7F || b == 0x1B || b == 0xC2 { // 0xC2 starts UTF-8 C1 controls \u0080..\u009F
			needsSanitize = true
			break
		}
	}
	if !needsSanitize {
		return text
	}

	var sb strings.Builder
	sb.Grow(len(text))

	i := 0
	n := len(text)
	for i < n {
		b := text[i]

		// 1. ESC (0x1B) sequence
		if b == 0x1B {
			i++
			if i >= n {
				break
			}
			next := text[i]
			switch next {
			case '[': // CSI (Control Sequence Introducer)
				i++
				// CSI parameter bytes (0x30..0x3F) and intermediate bytes (0x20..0x2F)
				for i < n && (text[i] >= 0x20 && text[i] <= 0x3F) {
					i++
				}
				// Final byte (0x40..0x7E)
				if i < n && text[i] >= 0x40 && text[i] <= 0x7E {
					i++
				}
			case ']': // OSC (Operating System Command)
				i++
				// Read until BEL (\x07) or ST (\x1b\\ or \u009c) or end of string
				for i < n {
					if text[i] == 0x07 { // BEL
						i++
						break
					}
					if text[i] == 0x1B && i+1 < n && text[i+1] == '\\' { // ESC \
						i += 2
						break
					}
					if text[i] == 0xC2 && i+1 < n && text[i+1] == 0x9C { // UTF-8 \u009c
						i += 2
						break
					}
					i++
				}
			case 'P', '_', '^', 'X': // DCS, APC, PM, SOS
				i++
				// Read until ST (ESC \ or \u009c) or BEL or end of string
				for i < n {
					if text[i] == 0x07 {
						i++
						break
					}
					if text[i] == 0x1B && i+1 < n && text[i+1] == '\\' {
						i += 2
						break
					}
					if text[i] == 0xC2 && i+1 < n && text[i+1] == 0x9C {
						i += 2
						break
					}
					i++
				}
			default:
				// 2-byte escape sequence: ESC followed by any printable character (e.g. \x1b=, \x1bc)
				if next >= 0x20 && next <= 0x7E {
					i++
				}
			}
			continue
		}

		// 2. Carriage Return (\r)
		if b == '\r' {
			if i+1 < n && text[i+1] == '\n' {
				sb.WriteByte('\n')
				i += 2
			} else {
				// Strip isolated \r to prevent terminal cursor overwrite attacks
				i++
			}
			continue
		}

		// 3. Permitted C0 controls (\n, \t)
		if b == '\n' || b == '\t' {
			sb.WriteByte(b)
			i++
			continue
		}

		// 4. Other C0 controls and DEL (0x7F)
		if b < 0x20 || b == 0x7F {
			i++
			continue
		}

		// 5. UTF-8 decoding & C1 control check
		if b >= 0x80 {
			r, size := utf8.DecodeRuneInString(text[i:])
			if r == utf8.RuneError && size == 1 {
				// Invalid UTF-8 byte
				i++
				continue
			}
			// C1 control characters (U+0080 through U+009F)
			if r >= 0x80 && r <= 0x9F {
				// Check for C1 CSI (\u009B)
				if r == 0x9B {
					i += size
					for i < n && (text[i] >= 0x20 && text[i] <= 0x3F) {
						i++
					}
					if i < n && text[i] >= 0x40 && text[i] <= 0x7E {
						i++
					}
					continue
				}
				// Check for C1 OSC (\u009D), DCS (\u0090), APC (\u009F), PM (\u009E)
				if r == 0x9D || r == 0x90 || r == 0x9F || r == 0x9E {
					i += size
					for i < n {
						if text[i] == 0x07 {
							i++
							break
						}
						if text[i] == 0x1B && i+1 < n && text[i+1] == '\\' {
							i += 2
							break
						}
						if text[i] == 0xC2 && i+1 < n && text[i+1] == 0x9C {
							i += 2
							break
						}
						i++
					}
					continue
				}
				// Other C1 controls: drop
				i += size
				continue
			}

			// Valid Unicode rune above U+009F
			sb.WriteString(text[i : i+size])
			i += size
			continue
		}

		// 6. Regular printable ASCII
		sb.WriteByte(b)
		i++
	}

	return sb.String()
}

// Terminal guards against emitting color codes when the destination does not
// support them, and distinguishes data output (stdout) from diagnostics (stderr).
type Terminal struct {
	// ColorSupported reports whether ANSI color codes should be emitted.
	ColorSupported bool
	// Out is the destination for data output (defaults to os.Stdout).
	Out io.Writer
	// Err is the destination for diagnostics and warnings (defaults to os.Stderr).
	Err io.Writer
}

// IsTerminal reports whether f is an interactive terminal character device.
func IsTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// NewTerminal detects color support from environment variables (NO_COLOR, CLICOLOR_FORCE, TERM)
// and verifies whether standard output is an interactive terminal.
func NewTerminal() *Terminal {
	return NewTerminalFor(os.Stdout)
}

// NewTerminalFor detects color support for the given file handle.
//
// The NO_COLOR convention (https://no-color.org) applies when the variable is
// present *and not empty*, so an exported-but-blank NO_COLOR does not silently
// disable colour for a whole process tree.
func NewTerminalFor(f *os.File) *Terminal {
	var colorSupported bool
	if noColor, ok := os.LookupEnv("NO_COLOR"); ok && noColor != "" {
		colorSupported = false
	} else if force, ok := os.LookupEnv("CLICOLOR_FORCE"); ok && force != "" && force != "0" {
		colorSupported = true
	} else if force, ok := os.LookupEnv("FORCE_COLOR"); ok && force != "" && force != "0" {
		colorSupported = true
	} else if ter, _ := os.LookupEnv("TERM"); ter == "dumb" {
		colorSupported = false
	} else {
		colorSupported = IsTerminal(f)
	}

	var out io.Writer = f
	if f == nil {
		out = os.Stdout
	}
	return &Terminal{
		ColorSupported: colorSupported,
		Out:            out,
		Err:            os.Stderr,
	}
}

// DataWriter returns the writer for data output (Out or os.Stdout).
func (t *Terminal) DataWriter() io.Writer {
	if t != nil && t.Out != nil {
		return t.Out
	}
	return os.Stdout
}

// DiagWriter returns the writer for diagnostics and warnings (Err or os.Stderr).
func (t *Terminal) DiagWriter() io.Writer {
	if t != nil && t.Err != nil {
		return t.Err
	}
	return os.Stderr
}

// Paint colors text, degrading to plain sanitized text when color is unsupported.
// text is sanitized by default to strip control and escape sequences.
func (t *Terminal) Paint(text string, color Color, styles []Style) string {
	if t == nil || !t.ColorSupported {
		return Sanitize(text)
	}
	return Paint(text, color, styles)
}

// RawPaint colors text without sanitization, degrading to raw text when color is unsupported.
// This is an explicit, trusted-only API and must not receive untrusted inputs.
func (t *Terminal) RawPaint(text string, color Color, styles []Style) string {
	if t == nil || !t.ColorSupported {
		return text
	}
	return RawPaint(text, color, styles)
}

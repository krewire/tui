// Package term provides terminal I/O, output formatting, and color
// conventions for the Krewire ecosystem.
package term

import (
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
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

// Sanitize removes terminal control sequences and control characters from
// untrusted text. It preserves ordinary Unicode, tabs, and newlines. Complete
// and unterminated ANSI CSI, OSC, DCS, SOS, PM, and APC sequences are removed.
func Sanitize(text string) string {
	runes := []rune(text)
	var clean strings.Builder
	clean.Grow(len(text))

	for i := 0; i < len(runes); {
		r := runes[i]
		switch r {
		case '\n', '\t':
			clean.WriteRune(r)
			i++
		case '\x1b':
			i = skipEscapeSequence(runes, i)
		case '\u009b':
			i = skipCSI(runes, i+1)
		case '\u009d':
			i = skipOSC(runes, i+1)
		case '\u0090', '\u0098', '\u009e', '\u009f':
			i = skipStringSequence(runes, i+1)
		default:
			if !unicode.IsControl(r) {
				clean.WriteRune(r)
			}
			i++
		}
	}
	return clean.String()
}

func skipEscapeSequence(runes []rune, start int) int {
	if start+1 >= len(runes) {
		return len(runes)
	}

	switch runes[start+1] {
	case '[':
		return skipCSI(runes, start+2)
	case ']':
		return skipOSC(runes, start+2)
	case 'P', 'X', '^', '_':
		return skipStringSequence(runes, start+2)
	default:
		// Three-byte ESC sequences include an intermediate byte (0x20-0x2F)
		// followed by a final byte (0x30-0x7E).
		if runes[start+1] >= 0x20 && runes[start+1] <= 0x2f &&
			start+2 < len(runes) && runes[start+2] >= 0x30 && runes[start+2] <= 0x7e {
			return start + 3
		}
		// Two-byte ESC sequences are discarded.
		return start + 2
	}
}

func skipCSI(runes []rune, start int) int {
	for i := start; i < len(runes); i++ {
		// ECMA-48 CSI final bytes occupy 0x40 through 0x7E.
		if runes[i] >= 0x40 && runes[i] <= 0x7e {
			return i + 1
		}
	}
	return len(runes)
}

func skipOSC(runes []rune, start int) int {
	for i := start; i < len(runes); i++ {
		switch runes[i] {
		case '\a', '\u009c', '\x18', '\x1a':
			return i + 1
		case '\x1b':
			if i+1 < len(runes) && runes[i+1] == '\\' {
				return i + 2
			}
		}
	}
	return len(runes)
}

func skipStringSequence(runes []rune, start int) int {
	for i := start; i < len(runes); i++ {
		switch runes[i] {
		case '\u009c', '\x18', '\x1a':
			return i + 1
		case '\x1b':
			if i+1 < len(runes) && runes[i+1] == '\\' {
				return i + 2
			}
		}
	}
	return len(runes)
}

// Paint sanitizes untrusted text and wraps it in ANSI escape sequences for the
// given color and styles. Use PaintTrustedANSI only for application-controlled
// payloads that intentionally contain terminal control sequences.
func Paint(text string, color Color, styles []Style) string {
	return PaintTrustedANSI(Sanitize(text), color, styles)
}

// PaintTrustedANSI wraps trusted text in ANSI escape sequences for the given
// color and styles without sanitizing its payload. Never pass repository,
// configuration, user, or other untrusted text to this function.
func PaintTrustedANSI(text string, color Color, styles []Style) string {
	codes := make([]string, 0, len(styles)+1)
	for _, s := range styles {
		codes = append(codes, strconv.Itoa(int(s)))
	}
	codes = append(codes, strconv.Itoa(int(color)))
	return "\x1b[" + strings.Join(codes, ";") + "m" + text + "\x1b[0m"
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

// Paint colors sanitized text, degrading to sanitized plain text when color is unsupported.
func (t *Terminal) Paint(text string, color Color, styles []Style) string {
	text = Sanitize(text)
	if !t.ColorSupported {
		return text
	}
	return PaintTrustedANSI(text, color, styles)
}

// PaintTrustedANSI emits the supplied trusted ANSI payload only when color is
// supported. When color is unsupported, embedded control sequences are removed.
func (t *Terminal) PaintTrustedANSI(text string, color Color, styles []Style) string {
	if !t.ColorSupported {
		return Sanitize(text)
	}
	return PaintTrustedANSI(text, color, styles)
}

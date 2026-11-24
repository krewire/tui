package term

import (
	"strings"
	"testing"
	"unicode"
)

func TestSanitize_RemovesTerminalControlSequences(t *testing.T) {
	input := "CSI:\x1b[2Jvisible " +
		"OSC8:\x1b]8;;https://evil.example\x07linked\x1b]8;;\x07 " +
		"OSC52:\x1b]52;c;Ym9keQ==\x07 " +
		"BEL:\x07 CR:\r C1-CSI:\u009b31mcolored " +
		"C1-OSC:\u009d52;c;Ym9keQ==\u009c"
	got := Sanitize(input)
	for _, want := range []string{"visible", "linked", "colored"} {
		if !strings.Contains(got, want) {
			t.Errorf("Sanitize() = %q, want to preserve %q", got, want)
		}
	}
	for _, forbidden := range []string{"evil.example", "Ym9keQ==", "\x1b", "\x07", "\r"} {
		if strings.Contains(got, forbidden) {
			t.Errorf("Sanitize() = %q, must not contain %q", got, forbidden)
		}
	}
	for _, r := range got {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			t.Errorf("Sanitize() retained control rune U+%04X", r)
		}
	}
}

func TestSanitize_PreservesUnicodeWhitespaceAndPlainText(t *testing.T) {
	input := "Halo 世界\nbaris kedua\tkolom"
	if got := Sanitize(input); got != input {
		t.Errorf("Sanitize() = %q, want %q", got, input)
	}
}

func TestSanitize_DropsUnterminatedSequences(t *testing.T) {
	for _, input := range []string{
		"before\x1b[31",
		"before\x1b]52;c;payload",
		"before\u009b31",
		"before\u009d52;c;payload",
	} {
		if got := Sanitize(input); got != "before" {
			t.Errorf("Sanitize(%q) = %q, want %q", input, got, "before")
		}
	}
}

func TestPaint_SanitizesPayloadAndKeepsGeneratedStyle(t *testing.T) {
	payload := "before\x1b[2Jhidden\x1b]52;c;Ym9keQ==\x07after"
	got := Paint(payload, ColorGreen, []Style{StyleBold})
	if strings.Contains(got, "Ym9keQ==") || strings.Contains(got, "\x1b]") || strings.Count(got, "\x1b") != 2 {
		t.Errorf("Paint() retained attacker-controlled terminal sequences: %q", got)
	}
	if !strings.HasPrefix(got, "\x1b[1;32m") || !strings.HasSuffix(got, "\x1b[0m") {
		t.Errorf("Paint() lost generated style wrapper: %q", got)
	}
	if !strings.Contains(got, "beforehiddenafter") {
		t.Errorf("Paint() lost safe text: %q", got)
	}
}

func TestPaintTrustedANSI_IsExplicitEscapeHatch(t *testing.T) {
	payload := "\x1b[31mtrusted\x1b[0m"
	if got := PaintTrustedANSI(payload, ColorGreen, nil); !strings.Contains(got, payload) {
		t.Errorf("PaintTrustedANSI() = %q, want trusted payload intact", got)
	}
	if got := (&Terminal{ColorSupported: false}).PaintTrustedANSI(payload, ColorGreen, nil); strings.Contains(got, "\x1b") {
		t.Errorf("PaintTrustedANSI() without color support = %q, want controls stripped", got)
	}
}

func TestSanitize_RemovesThreeByteEscapeSequences(t *testing.T) {
	if got := Sanitize("before\x1b(Bafter"); got != "beforeafter" {
		t.Errorf("Sanitize() = %q, want %q", got, "beforeafter")
	}
}

func TestTerminalPaint_SanitizesBeforeAddingStyle(t *testing.T) {
	got := (&Terminal{ColorSupported: true}).Paint("before\x1b[2Jhidden", ColorGreen, []Style{StyleBold})
	if strings.Count(got, "\x1b") != 2 || !strings.HasPrefix(got, "\x1b[1;32m") || !strings.HasSuffix(got, "\x1b[0m") {
		t.Errorf("Terminal.Paint() = %q, want only generated ANSI wrappers", got)
	}
	if !strings.Contains(got, "beforehidden") {
		t.Errorf("Terminal.Paint() lost safe text: %q", got)
	}
}

package term_test

import (
	"testing"

	"github.com/krewire/tui/term"
)

func TestRuneWidth(t *testing.T) {
	cases := []struct {
		name string
		r    rune
		want int
	}{
		{"ASCII letter", 'a', 1},
		{"ASCII digit", '5', 1},
		{"ASCII space", ' ', 1},
		{"Null byte", 0, 0},
		{"Control char", '\x07', 0},
		{"Hangul Syllable", '한', 2},
		{"CJK Ideograph", '字', 2},
		{"Fullwidth Latin", 'Ａ', 2},
		{"Emoji Grinning", '😀', 2},
		{"Emoji Rocket", '🚀', 2},
		{"Zero-width space", 0x200B, 0},
		{"Zero-width joiner", 0x200D, 0},
		{"Combining acute accent", 0x0301, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := term.RuneWidth(tc.r)
			if got != tc.want {
				t.Errorf("RuneWidth(%q [%U]) = %d, want %d", tc.r, tc.r, got, tc.want)
			}
		})
	}
}

func TestStringWidth(t *testing.T) {
	cases := []struct {
		name string
		s    string
		want int
	}{
		{"empty", "", 0},
		{"plain ascii", "hello world", 11},
		{"cjk text", "你好世界", 8},
		{"mixed ascii and cjk", "Krewire 框架", 12}, // 7 + 1 + 4 = 12
		{"ansi colored text", "\x1b[31mhello\x1b[0m", 5},
		{"combining characters", "e\u0301", 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := term.StringWidth(tc.s)
			if got != tc.want {
				t.Errorf("StringWidth(%q) = %d, want %d", tc.s, got, tc.want)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	t.Run("zero or negative width returns unmodified", func(t *testing.T) {
		text := "quick brown fox"
		if got := term.Wrap(text, 0); got != text {
			t.Errorf("Wrap(..., 0) = %q, want %q", got, text)
		}
		if got := term.Wrap(text, -5); got != text {
			t.Errorf("Wrap(..., -5) = %q, want %q", got, text)
		}
	})

	t.Run("wraps at word boundaries", func(t *testing.T) {
		text := "The quick brown fox jumps over the lazy dog"
		wrapped := term.Wrap(text, 15)
		expected := "The quick brown\nfox jumps over\nthe lazy dog"
		if wrapped != expected {
			t.Errorf("Wrap() = %q, want %q", wrapped, expected)
		}
	})

	t.Run("preserves explicit newlines", func(t *testing.T) {
		text := "Line 1\nLine 2 is somewhat longer\nLine 3"
		wrapped := term.Wrap(text, 15)
		expected := "Line 1\nLine 2 is\nsomewhat longer\nLine 3"
		if wrapped != expected {
			t.Errorf("Wrap() = %q, want %q", wrapped, expected)
		}
	})

	t.Run("breaks words exceeding limit", func(t *testing.T) {
		text := "supercalifragilisticexpialidocious"
		wrapped := term.Wrap(text, 10)
		expected := "supercalif\nragilistic\nexpialidoc\nious"
		if wrapped != expected {
			t.Errorf("Wrap() = %q, want %q", wrapped, expected)
		}
	})

	t.Run("handles CJK wide characters wrapping", func(t *testing.T) {
		text := "你好 世界 欢迎 来到 这里"
		wrapped := term.Wrap(text, 10)
		// "你好 世界" is 4 + 1 + 4 = 9 width <= 10
		// "欢迎 来到" is 4 + 1 + 4 = 9 width <= 10
		// "这里" is 4 width <= 10
		expected := "你好 世界\n欢迎 来到\n这里"
		if wrapped != expected {
			t.Errorf("Wrap() = %q, want %q", wrapped, expected)
		}
	})
}

func TestTruncate(t *testing.T) {
	cases := []struct {
		name  string
		text  string
		width int
		tail  string
		want  string
	}{
		{"text fits", "hello", 10, "...", "hello"},
		{"exact fit", "hello", 5, "...", "hello"},
		{"truncate with ellipsis", "hello world", 8, "...", "hello..."},
		{"zero width", "hello", 0, "...", ""},
		{"negative width", "hello", -2, "...", ""},
		{"tail wider than width", "hello", 2, "...", ".."},
		{"truncate wide characters", "你好世界欢迎你", 9, "...", "你好世..."}, // 2*3 + 3 = 9 <= 9
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := term.Truncate(tc.text, tc.width, tc.tail)
			if got != tc.want {
				t.Errorf("Truncate(%q, %d, %q) = %q, want %q", tc.text, tc.width, tc.tail, got, tc.want)
			}
		})
	}
}

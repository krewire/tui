# `term` Module Reference

- **Module Path**: `github.com/krewire/tui/term` (`term/go.mod`)
- **Package Name**: `package term`

## Purpose & Domain Responsibility

`term` is a zero-dependency leaf package providing ANSI 16-color and text-style rendering, terminal capability detection, stdout/stderr stream separation, and defense-in-depth terminal injection sanitization.

## Architectural Position

- **Layer**: Layer 2 (Leaf Terminal Primitive)
- **Pure-Go Invariant**: Compiles cleanly with `CGO_ENABLED=0`
- **Dependencies**:
  - Third-party: None
  - Ecosystem: None (strictly Go standard library: `io`, `os`, `strconv`, `strings`, `unicode/utf8`)
  - **Leaf Invariant**: Enforced by `tui/core/architecture_test.go`, which fails if `term/` imports `github.com/krewire/tui` or any non-stdlib package.

## Core Concepts & APIs

### Colors & Styles

- **`Color`**: ANSI foreground palette constants:
  - `ColorBlack` (30), `ColorRed` (31), `ColorGreen` (32), `ColorYellow` (33), `ColorBlue` (34), `ColorMagenta` (35), `ColorCyan` (36), `ColorWhite` (37), `ColorDefault` (39).
- **`Style`**: Text attribute constants:
  - `StyleBold` (1), `StyleDim` (2), `StyleUnderline` (4).

### Safe Painting & Terminal Injection Protection

- **`Sanitize(text string) string`**: Strips ANSI/VT100 escape sequences (CSI, OSC, DCS, APC, PM, SOS, and 2-byte ESC sequences), C0 control characters (except `\n` and `\t`), DEL (`0x7F`), and C1 control characters (`U+0080`–`U+009F`). Isolated `\r` bytes are removed to prevent cursor-rewind spoofing, while `\r\n` is normalized to `\n` and valid Unicode above `U+009F` is preserved. Uses a zero-allocation fast path when no control bytes are present.
- **`Paint(text string, color Color, styles []Style) string`**: Sanitizes `text` via `Sanitize` and wraps it with ANSI SGR codes via `RawPaint`.
- **`RawPaint(text string, color Color, styles []Style) string`**: Trusted-only API that wraps `text` in `\x1b[<codes>m...\x1b[0m` without sanitization.

### Terminal Capability Detection & Streams

- **`Terminal`**: `{ ColorSupported bool, Out io.Writer, Err io.Writer }`.
- **`IsTerminal(f *os.File) bool`**: Reports whether `f` is a character device (`os.ModeCharDevice`).
- **`NewTerminal() *Terminal`** and **`NewTerminalFor(f *os.File) *Terminal`**: Detect color support in priority order:
  1. `NO_COLOR` present and non-empty → `false` (per no-color.org)
  2. `CLICOLOR_FORCE` or `FORCE_COLOR` present, non-empty, and `!= "0"` → `true`
  3. `TERM=dumb` → `false`
  4. Otherwise `IsTerminal(f)`
- **Stream Accessors**: **`(t *Terminal) DataWriter() io.Writer`** (defaults to `os.Stdout`) and **`(t *Terminal) DiagWriter() io.Writer`** (defaults to `os.Stderr`).
- **Guarded Painting**: **`(t *Terminal) Paint(text, color, styles)`** emits ANSI codes only when `t != nil && t.ColorSupported`, otherwise returning `Sanitize(text)`. **`(t *Terminal) RawPaint(text, color, styles)`** degrades to raw `text` when color is unsupported.

## Usage Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/krewire/tui/term"
)

func main() {
	tm := term.NewTerminalFor(os.Stderr)

	// Untrusted input is sanitized automatically by Paint
	badge := tm.Paint("READY", term.ColorGreen, []term.Style{term.StyleBold})
	msg := tm.Paint("user\x1b[31m-input", term.ColorCyan, nil)

	fmt.Fprintf(tm.DiagWriter(), "[%s] %s\n", badge, msg)
}
```

## Verification

```bash
cd /home/krewire/workspace/krewire/tui/term
gofmt -l .
go vet ./...
CGO_ENABLED=0 go test -v ./...
CGO_ENABLED=0 go build ./...
```

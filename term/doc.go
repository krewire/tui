// Package term provides low-level terminal I/O primitives, ANSI styling,
// and terminal capability detection for the Krewire ecosystem.
//
// # Boundary & Invariants
//
// The term package operates as a standalone terminal abstraction layer:
//   - Zero dependencies on the parent github.com/krewire/tui package.
//   - Zero external third-party dependencies (strictly Go standard library).
//   - Respects NO_COLOR, TERM=dumb, and non-TTY file descriptor redirection.
//   - Pure-Go implementation compatible with CGO_ENABLED=0 across all platforms.
package term

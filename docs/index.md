# Krewire TUI — Developer & Contributor Documentation

Welcome to the internal engineering documentation for **`krewire/tui`** (`github.com/krewire/tui`), the lightweight, pure-Go terminal user interface toolkit and ANSI styling engine for the Krewire digital SDLC ecosystem.

## Documentation Index

- [Architecture & Design](architecture.md): Two-package design (`core` / `package tui` and `term`), leaf invariant, downward dependency rules, and terminal primitives.
- [Getting Started](getting-started.md): Installation, quickstart examples, command dispatching, flag registration, and formatting.
- [Module Catalog](modules/index.md): Reference documentation for [`core`](modules/core.md) (`github.com/krewire/tui/core`) and [`term`](modules/term.md) (`github.com/krewire/tui/term`).
- [Contributing Guide](../CONTRIBUTING.md): Quality gates, pure-Go invariant, testing standards, and pull request workflow.

## Overview & Core Principles

- **Zero External Dependencies**: Implemented strictly using the Go standard library (`os`, `fmt`, `strings`, `strconv`, `io`, `unicode/utf8`, `flag`, `encoding/json`, `log/slog`, `sort`) plus `github.com/krewire/krewire/errs` in `core`.
- **Pure Go (`CGO_ENABLED=0`)**: Completely CGO-free, enabling seamless cross-compilation across Linux, macOS, and Windows.
- **Two-Module Structure**:
  - [`github.com/krewire/tui/core`](modules/core.md) (`package tui`): CLI command routing, flag binding, progressive-disclosure help, `<APP>_` environment config, JSON output, and typo suggestions.
  - [`github.com/krewire/tui/term`](modules/term.md) (`package term`): Leaf ANSI styling, terminal injection sanitization (`Sanitize`), capability detection (`NO_COLOR`, `TERM=dumb`, TTY), and stream separation.
- **Leaf Invariant**: `term` must never import `core` (`tui`) or any package outside the Go standard library, enforced by `core/architecture_test.go`.
